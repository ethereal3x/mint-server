package service

import (
	"context"

	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AgentServer AgentService gRPC 处理器
type AgentServer struct {
	agentpb.UnimplementedAgentServiceServer
	chat    *logic.Chat
	mapping *logic.Mapping
}

// NewAgentServer 创建 Agent gRPC 处理器
func NewAgentServer(chat *logic.Chat, mapping *logic.Mapping) *AgentServer {
	return &AgentServer{chat: chat, mapping: mapping}
}

type streamResult struct {
	usage *schema.TokenUsage
	err   error
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
	}

	contentChan := make(chan string)
	resultChan := make(chan streamResult, 1)

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	go func() {
		usage, err := s.chat.StreamChat(ctx, chatReq, contentChan)
		select {
		case resultChan <- streamResult{usage: usage, err: err}:
		case <-ctx.Done():
		}
	}()

	var fullAnswer string
	for {
		select {
		case content, ok := <-contentChan:
			if !ok {
				result := <-resultChan
				if result.err != nil {
					return status.Errorf(codes.Internal, "chat stream: %v", result.err)
				}
				if saveErr := s.chat.SaveRecord(ctx, chatReq, fullAnswer, result.usage); saveErr != nil {
					return status.Errorf(codes.Internal, "save record: %v", saveErr)
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
	if req.Question == "" || req.Model == "" {
		return nil, status.Errorf(codes.InvalidArgument, "chat: missing question or model")
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
	}

	answer, usage, err := s.chat.GenerateChat(ctx, chatReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "chat: %v", err)
	}

	if err := s.chat.SaveRecord(ctx, chatReq, answer, usage); err != nil {
		return nil, status.Errorf(codes.Internal, "save record: %v", err)
	}

	return &agentpb.GenerateChatResponse{
		Content:    answer,
		DialogueId: dialogueID,
		RecordId:   recordID,
	}, nil
}

// ListModels 获取可用模型列表
func (s *AgentServer) ListModels(ctx context.Context, _ *agentpb.ListModelsRequest) (*agentpb.ListModelsResponse, error) {
	mappings, err := s.mapping.ListAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list models: %v", err)
	}
	models := make([]*agentpb.ModelInfo, 0, len(mappings))
	for _, m := range mappings {
		models = append(models, &agentpb.ModelInfo{
			Model:        m.ModelType,
			Manufacturer: m.Manufacturer,
			Description:  m.Description,
		})
	}
	return &agentpb.ListModelsResponse{Models: models}, nil
}

// GetHistory 获取对话历史
func (s *AgentServer) GetHistory(ctx context.Context, req *agentpb.GetHistoryRequest) (*agentpb.GetHistoryResponse, error) {
	if req.DialogueId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "get history: missing dialogue_id")
	}
	records, err := s.chat.GetHistory(ctx, req.DialogueId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get history: %v", err)
	}
	pbRecords := make([]*agentpb.HistoryRecord, 0, len(records))
	for _, r := range records {
		pbRecords = append(pbRecords, &agentpb.HistoryRecord{
			RecordId:     r.RecordID,
			UserContent:  r.UserContent,
			AgentContent: r.AgentContent,
			Model:        r.Model,
			TotalTokens:  r.TotalToken,
		})
	}
	return &agentpb.GetHistoryResponse{Records: pbRecords}, nil
}

// ListDialogues 获取对话摘要列表
func (s *AgentServer) ListDialogues(ctx context.Context, req *agentpb.ListDialoguesRequest) (*agentpb.ListDialoguesResponse, error) {
	summaries, err := s.chat.ListDialogues(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list dialogues: %v", err)
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
	return &agentpb.ListDialoguesResponse{Dialogues: pbDialogues}, nil
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
