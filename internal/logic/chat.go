package logic

import (
	"context"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/model"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

// ModelConfigRepo 聊天场景需要的模型配置查询接口
type ModelConfigRepo interface {
	FindByModelType(ctx context.Context, modelType string) (*model.ChatModelConfig, error)
}

// ChatRecordRepo 聊天场景需要的对话记录接口
type ChatRecordRepo interface {
	ListByDialogueID(ctx context.Context, dialogueID string) ([]*model.DialogueRecord, error)
	Create(ctx context.Context, record *model.DialogueRecord) error
	ListDialogues(ctx context.Context, userID string) ([]*model.DialogueSummary, error)
	AggregateByModel(ctx context.Context) ([]*model.ModelStat, error)
}

// Chat 聊天业务逻辑，编排模型配置查询、上下文构建、模型调用、记录落库
type Chat struct {
	configRepo ModelConfigRepo
	recordRepo ChatRecordRepo
	adapter    LLMAdapter
}

// ChatResult 聊天调用结果
type ChatResult struct {
	Config  *model.ChatModelConfig
	Content string
	Usage   *schema.TokenUsage
}

// NewChat 创建聊天业务逻辑
func NewChat(configRepo ModelConfigRepo, recordRepo ChatRecordRepo, adapter LLMAdapter) *Chat {
	return &Chat{configRepo: configRepo, recordRepo: recordRepo, adapter: adapter}
}

// ChatRequest 聊天请求参数
type ChatRequest struct {
	UserID     string
	Question   string
	Model      string
	DialogueID string
	RecordID   string
	FileData   []byte
	FileName   string
	ImageURLs  []string
}

// prepareChat 查询模型配置→加载历史→构建上下文
func (s *Chat) prepareChat(ctx context.Context, req *ChatRequest) (*model.ChatModelConfig, []*schema.Message, error) {
	logger.ContextInfo(ctx, "Chat.prepareChat", zap.String("model", req.Model), zap.String("dialogue_id", req.DialogueID))

	config, err := s.configRepo.FindByModelType(ctx, req.Model)
	if err != nil {
		logger.ContextError(ctx, "Chat.prepareChat", zap.String("model", req.Model), zap.Error(err))
		return nil, nil, mint_err.ErrDBQuery
	}
	if config == nil {
		logger.ContextError(ctx, "Chat.prepareChat", zap.String("model", req.Model))
		return nil, nil, mint_err.ErrModelNotFound
	}

	history, err := s.recordRepo.ListByDialogueID(ctx, req.DialogueID)
	if err != nil {
		logger.ContextError(ctx, "Chat.prepareChat", zap.String("dialogue_id", req.DialogueID), zap.Error(err))
		return nil, nil, mint_err.ErrGetHistory
	}

	messages := NewContextBuilder().
		AddHistory(history).
		AddFile(req.FileData, req.FileName).
		AddUserQuestion(req.Question).
		AddImages(req.ImageURLs).
		Build()

	return config, messages, nil
}

// StreamChat 执行流式聊天，增量内容通过 contentChan 返回
func (s *Chat) StreamChat(ctx context.Context, req *ChatRequest, contentChan chan<- string) (*ChatResult, error) {
	ctx, span := tracing.Start(ctx, "logic.Chat.StreamChat")
	defer span.End()

	defer close(contentChan)

	config, messages, err := s.prepareChat(ctx, req)
	if err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	usage, err := s.adapter.Stream(ctx, config, messages, contentChan)
	if err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}
	return &ChatResult{Config: config, Usage: usage}, nil
}

// GenerateChat 非流式生成聊天内容
func (s *Chat) GenerateChat(ctx context.Context, req *ChatRequest) (*ChatResult, error) {
	ctx, span := tracing.Start(ctx, "logic.Chat.GenerateChat")
	defer span.End()

	config, messages, err := s.prepareChat(ctx, req)
	if err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	content, usage, err := s.adapter.Generate(ctx, config, messages)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.GenerateChat", zap.Error(err))
		return nil, mint_err.ErrChatGenerate
	}
	return &ChatResult{Config: config, Content: content, Usage: usage}, nil
}

// SaveRecord 持久化对话记录，根据 config 中的单价计算费用
func (s *Chat) SaveRecord(ctx context.Context, req *ChatRequest, answer string, config *model.ChatModelConfig, usage *schema.TokenUsage) error {
	ctx, span := tracing.Start(ctx, "logic.Chat.SaveRecord")
	defer span.End()

	record := &model.DialogueRecord{
		DialogueID:   req.DialogueID,
		RecordID:     req.RecordID,
		UserID:       req.UserID,
		Model:        req.Model,
		UserContent:  req.Question,
		AgentContent: answer,
	}
	if usage != nil {
		record.UserTokens = int64(usage.PromptTokens)
		record.AgentTokens = int64(usage.CompletionTokens)
		record.TotalTokens = int64(usage.TotalTokens)

		if config != nil {
			record.InputCost = float64(record.UserTokens) * config.InputPrice / 1000000
			record.OutputCost = float64(record.AgentTokens) * config.OutputPrice / 1000000
		}
	}
	if err := s.recordRepo.Create(ctx, record); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.SaveRecord", zap.String("record_id", record.RecordID), zap.Error(err))
		return mint_err.ErrSaveRecord
	}
	return nil
}
