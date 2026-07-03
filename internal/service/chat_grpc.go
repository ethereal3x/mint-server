package service

import (
	"context"

	"github.com/ethereal3x/apc/errs"
	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	"github.com/ethereal3x/mint-server/internal/dto"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/model"
	"google.golang.org/grpc"
)

// AgentServer AgentService gRPC 处理器
type AgentServer struct {
	agentpb.UnimplementedAgentServiceServer
	chat   AgentChatLogic
	config AgentConfigLogic
}

// NewAgentServer 创建 Agent gRPC 处理器
func NewAgentServer(chat AgentChatLogic, config AgentConfigLogic) *AgentServer {
	return &AgentServer{chat: chat, config: config}
}

// chatPrepareParams 聊天请求准备参数
type chatPrepareParams struct {
	question   string
	model      string
	dialogueID string
	recordID   string
	fileData   []byte
	fileName   string
	imageURLs  []string
}

// buildChatRequest 鉴权、生成 IDs 并构建 ChatRequest
func (server *AgentServer) buildChatRequest(ctx context.Context, params *chatPrepareParams) (string, string, *dto.ChatRequest, error) {
	userID, err := requireUserID(ctx)
	if err != nil {
		return "", "", nil, err
	}
	dialogueID, recordID := logic.ResolveChatIDs(params.dialogueID, params.recordID)
	chatReq := &dto.ChatRequest{
		UserID:     userID,
		Question:   params.question,
		Model:      params.model,
		DialogueID: dialogueID,
		RecordID:   recordID,
		FileData:   params.fileData,
		FileName:   params.fileName,
		ImageURLs:  params.imageURLs,
	}
	return dialogueID, recordID, chatReq, nil
}

// StreamChat 流式聊天
func (server *AgentServer) StreamChat(req *agentpb.StreamChatRequest, stream grpc.ServerStreamingServer[agentpb.StreamChatResponse]) error {
	if req.Question == "" || req.Model == "" {
		return sendStreamError(stream, mint_err.ErrParam)
	}
	dialogueID, recordID, chatReq, err := server.buildChatRequest(stream.Context(), &chatPrepareParams{
		question: req.Question, model: req.Model,
		dialogueID: req.DialogueId, recordID: req.RecordId,
		fileData: req.FileData, fileName: req.FileName, imageURLs: req.ImageUrls,
	})
	if err != nil {
		return sendStreamError(stream, err)
	}

	contentChan := make(chan string)
	resultChan := make(chan struct {
		result *dto.ChatResult
		err    error
	}, 1)

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	go func() {
		result, streamErr := server.chat.StreamChat(ctx, chatReq, contentChan)
		select {
		case resultChan <- struct {
			result *dto.ChatResult
			err    error
		}{result, streamErr}:
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
				if saveErr := server.chat.SaveRecord(ctx, &dto.SaveRecordRequest{ChatRequest: chatReq, Answer: fullAnswer, Config: res.result.Config, Usage: res.result.Usage}); saveErr != nil {
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
func (server *AgentServer) GenerateChat(ctx context.Context, req *agentpb.GenerateChatRequest) (*agentpb.GenerateChatResponse, error) {
	return errs.Handle(&agentpb.GenerateChatResponse{}, func(rsp *agentpb.GenerateChatResponse) error {
		if req.Question == "" || req.Model == "" {
			return mint_err.ErrParam
		}
		dialogueID, recordID, chatReq, err := server.buildChatRequest(ctx, &chatPrepareParams{
			question: req.Question, model: req.Model,
			dialogueID: req.DialogueId, recordID: req.RecordId,
			fileData: req.FileData, fileName: req.FileName, imageURLs: req.ImageUrls,
		})
		if err != nil {
			return err
		}
		result, err := server.chat.GenerateChat(ctx, chatReq)
		if err != nil {
			return err
		}
		if saveErr := server.chat.SaveRecord(ctx, &dto.SaveRecordRequest{ChatRequest: chatReq, Answer: result.Content, Config: result.Config, Usage: result.Usage}); saveErr != nil {
			return saveErr
		}
		rsp.Content = result.Content
		rsp.DialogueId = dialogueID
		rsp.RecordId = recordID
		return nil
	})
}

// ListModels 获取可用模型列表，按厂商分组
func (server *AgentServer) ListModels(ctx context.Context, _ *agentpb.ListModelsRequest) (*agentpb.ListModelsResponse, error) {
	return errs.Handle(&agentpb.ListModelsResponse{}, func(rsp *agentpb.ListModelsResponse) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		configs, err := server.config.ListAll(ctx, userID)
		if err != nil {
			return err
		}
		rsp.Manufacturers = logic.GroupModelsByManufacturer(configs)
		return nil
	})
}

// GetHistory 获取对话历史
func (server *AgentServer) GetHistory(ctx context.Context, req *agentpb.GetHistoryRequest) (*agentpb.GetHistoryResponse, error) {
	return errs.Handle(&agentpb.GetHistoryResponse{}, func(rsp *agentpb.GetHistoryResponse) error {
		if req.DialogueId == "" {
			return mint_err.ErrParam
		}
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		records, err := server.chat.GetHistory(ctx, &model.DialogueQuery{DialogueID: req.DialogueId, UserID: userID})
		if err != nil {
			return err
		}
		pbRecords := make([]*agentpb.HistoryRecord, 0, len(records))
		for _, record := range records {
			pbRecords = append(pbRecords, &agentpb.HistoryRecord{
				RecordId:     record.RecordID,
				UserContent:  record.UserContent,
				AgentContent: record.AgentContent,
				Model:        record.Model,
				TotalTokens:  record.TotalTokens,
				InputCost:    record.InputCost,
				OutputCost:   record.OutputCost,
			})
		}
		rsp.Records = pbRecords
		return nil
	})
}

// ListDialogues 获取对话摘要列表
func (server *AgentServer) ListDialogues(ctx context.Context, _ *agentpb.ListDialoguesRequest) (*agentpb.ListDialoguesResponse, error) {
	return errs.Handle(&agentpb.ListDialoguesResponse{}, func(rsp *agentpb.ListDialoguesResponse) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		summaries, err := server.chat.ListDialogues(ctx, userID)
		if err != nil {
			return err
		}
		pbDialogues := make([]*agentpb.DialogueSummary, 0, len(summaries))
		for _, summary := range summaries {
			pbDialogues = append(pbDialogues, &agentpb.DialogueSummary{
				DialogueId:   summary.DialogueID,
				Title:        summary.Title,
				MessageCount: summary.MessageCount,
				UpdatedTime:  summary.UpdatedTime,
			})
		}
		rsp.Dialogues = pbDialogues
		return nil
	})
}

// sendStreamError 将 error 写入流式响应并发送，标记流结束
func sendStreamError(stream grpc.ServerStreamingServer[agentpb.StreamChatResponse], err error) error {
	rsp := &agentpb.StreamChatResponse{Done: true}
	if setErr := errs.SetErrMsg(rsp, err); setErr != nil {
		rsp.Code = int32(mint_err.ERR_CODE_INTERNAL)
		rsp.Message = err.Error()
	}
	return stream.Send(rsp)
}
