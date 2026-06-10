package service

import (
	"context"

	"github.com/ethereal3x/apc/errs"
	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AgentServer AgentService gRPC 处理器
type AgentServer struct {
	agentpb.UnimplementedAgentServiceServer
	chat   *logic.Chat
	config *logic.Config
}

// NewAgentServer 创建 Agent gRPC 处理器
func NewAgentServer(chat *logic.Chat, config *logic.Config) *AgentServer {
	return &AgentServer{chat: chat, config: config}
}

// StreamChat 流式聊天
func (s *AgentServer) StreamChat(req *agentpb.StreamChatRequest, stream grpc.ServerStreamingServer[agentpb.StreamChatResponse]) error {
	if req.Question == "" || req.Model == "" {
		return status.Errorf(codes.InvalidArgument, "chat stream: missing question or model")
	}
	dialogueID, recordID := resolveIDs(req.DialogueId, req.RecordId)
	chatReq := &logic.ChatRequest{
		UserID:     req.UserId,
		Question:   req.Question,
		Model:      req.Model,
		DialogueID: dialogueID,
		RecordID:   recordID,
		FileData:   req.FileData,
		FileName:   req.FileName,
		ImageURLs:  req.ImageUrls,
	}

	contentChan := make(chan string)
	resultChan := make(chan struct {
		result *logic.ChatResult
		err    error
	}, 1)

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	go func() {
		result, err := s.chat.StreamChat(ctx, chatReq, contentChan)
		select {
		case resultChan <- struct {
			result *logic.ChatResult
			err    error
		}{result, err}:
		case <-ctx.Done():
		}
	}()

	var fullAnswer string
	for {
		select {
		case content, ok := <-contentChan:
			if !ok {
				res := <-resultChan
				if res.err != nil {
					return sendStreamError(stream, res.err)
				}
				if saveErr := s.chat.SaveRecord(ctx, chatReq, fullAnswer, res.result.Config, res.result.Usage); saveErr != nil {
					return sendStreamError(stream, saveErr)
				}
				return stream.Send(&agentpb.StreamChatResponse{Done: true, DialogueId: dialogueID, RecordId: recordID})
			}
			fullAnswer += content
			if sendErr := stream.Send(&agentpb.StreamChatResponse{Content: content, DialogueId: dialogueID, RecordId: recordID}); sendErr != nil {
				return sendErr
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// GenerateChat 非流式聊天
func (s *AgentServer) GenerateChat(ctx context.Context, req *agentpb.GenerateChatRequest) (*agentpb.GenerateChatResponse, error) {
	rsp := &agentpb.GenerateChatResponse{}
	if req.Question == "" || req.Model == "" {
		return errs.GenProtoReply(rsp, mint_err.ErrParam)
	}
	dialogueID, recordID := resolveIDs(req.DialogueId, req.RecordId)
	chatReq := &logic.ChatRequest{
		UserID:     req.UserId,
		Question:   req.Question,
		Model:      req.Model,
		DialogueID: dialogueID,
		RecordID:   recordID,
		FileData:   req.FileData,
		FileName:   req.FileName,
		ImageURLs:  req.ImageUrls,
	}

	result, err := s.chat.GenerateChat(ctx, chatReq)
	if err != nil {
		return errs.GenProtoReply(rsp, err)
	}

	if saveErr := s.chat.SaveRecord(ctx, chatReq, result.Content, result.Config, result.Usage); saveErr != nil {
		return errs.GenProtoReply(rsp, saveErr)
	}

	rsp.Content = result.Content
	rsp.DialogueId = dialogueID
	rsp.RecordId = recordID
	return rsp, nil
}

// ListModels 获取可用模型列表，按厂商分组
func (s *AgentServer) ListModels(ctx context.Context, _ *agentpb.ListModelsRequest) (*agentpb.ListModelsResponse, error) {
	rsp := &agentpb.ListModelsResponse{}
	configs, err := s.config.ListAll(ctx)
	if err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	groupMap := make(map[string][]*agentpb.ModelInfo)
	for _, c := range configs {
		groupMap[c.Manufacturer] = append(groupMap[c.Manufacturer], &agentpb.ModelInfo{
			Model:       c.ModelType,
			Description: c.Description,
		})
	}
	manufacturers := make([]*agentpb.ManufacturerGroup, 0, len(groupMap))
	for name, models := range groupMap {
		manufacturers = append(manufacturers, &agentpb.ManufacturerGroup{
			Manufacturer: name,
			Models:       models,
		})
	}
	rsp.Manufacturers = manufacturers
	return rsp, nil
}

// GetHistory 获取对话历史
func (s *AgentServer) GetHistory(ctx context.Context, req *agentpb.GetHistoryRequest) (*agentpb.GetHistoryResponse, error) {
	rsp := &agentpb.GetHistoryResponse{}
	if req.DialogueId == "" {
		return errs.GenProtoReply(rsp, mint_err.ErrParam)
	}
	records, err := s.chat.GetHistory(ctx, req.DialogueId)
	if err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	pbRecords := make([]*agentpb.HistoryRecord, 0, len(records))
	for _, r := range records {
		pbRecords = append(pbRecords, &agentpb.HistoryRecord{
			RecordId:     r.RecordID,
			UserContent:  r.UserContent,
			AgentContent: r.AgentContent,
			Model:        r.Model,
			TotalTokens:  r.TotalTokens,
			InputCost:    r.InputCost,
			OutputCost:   r.OutputCost,
		})
	}
	rsp.Records = pbRecords
	return rsp, nil
}

// ListDialogues 获取对话摘要列表
func (s *AgentServer) ListDialogues(ctx context.Context, req *agentpb.ListDialoguesRequest) (*agentpb.ListDialoguesResponse, error) {
	rsp := &agentpb.ListDialoguesResponse{}
	summaries, err := s.chat.ListDialogues(ctx, req.UserId)
	if err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	pbDialogues := make([]*agentpb.DialogueSummary, 0, len(summaries))
	for _, d := range summaries {
		pbDialogues = append(pbDialogues, &agentpb.DialogueSummary{
			DialogueId:   d.DialogueID,
			Title:        d.Title,
			MessageCount: d.MessageCount,
			UpdatedTime:  d.UpdatedTime,
		})
	}
	rsp.Dialogues = pbDialogues
	return rsp, nil
}

func resolveIDs(dialogueID, recordID string) (string, string) {
	if dialogueID == "" {
		dialogueID = uuid.NewString()
	}
	if recordID == "" {
		recordID = uuid.NewString()
	}
	return dialogueID, recordID
}

// sendStreamError 将 error 写入流式响应并发送
func sendStreamError(stream grpc.ServerStreamingServer[agentpb.StreamChatResponse], err error) error {
	rsp := &agentpb.StreamChatResponse{}
	if setErr := errs.SetErrMsg(rsp, err); setErr != nil {
		rsp.Code = int32(mint_err.ERR_CODE_INTERNAL)
		rsp.Message = err.Error()
	}
	return stream.Send(rsp)
}
