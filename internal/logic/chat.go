package logic

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"github.com/ethereal3x/mint-server/internal/model"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"go.uber.org/zap"
)

const documentPromptPrefix = "以下是用户传输文件内容，根据用户提问结合文件内容回答："

// ChatMappingRepo 聊天场景需要的模型映射查询接口
type ChatMappingRepo interface {
	FindByModelType(ctx context.Context, modelType string) (*model.ModelMapping, error)
}

// ChatRecordRepo 聊天场景需要的对话记录接口
type ChatRecordRepo interface {
	ListByDialogueID(ctx context.Context, dialogueID string) ([]*model.DialogueRecord, error)
	Create(ctx context.Context, record *model.DialogueRecord) error
	ListDialogues(ctx context.Context, userID string) ([]*model.DialogueSummary, error)
}

// Chat 聊天业务逻辑，编排策略解析、模型调用、记录落库
type Chat struct {
	mappingRepo ChatMappingRepo
	recordRepo  ChatRecordRepo
	strategy    *Strategy
}

// NewChat 创建聊天业务逻辑
func NewChat(mappingRepo ChatMappingRepo, recordRepo ChatRecordRepo, strategy *Strategy) *Chat {
	return &Chat{
		mappingRepo: mappingRepo,
		recordRepo:  recordRepo,
		strategy:    strategy,
	}
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
}

// newChatModel 从策略创建 Eino ChatModel，rule 的 API Key 必须已由 GetByManufacturer 解密
func newChatModel(ctx context.Context, rule *model.StrategyRule) (*openai.ChatModel, error) {
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  rule.APIKey,
		Model:   rule.AgentModel,
		BaseURL: rule.URL,
	})
}

// prepareChat 执行模型映射查询→策略获取→历史加载→消息组装，是 StreamChat/GenerateChat 的公共前置步骤
func (s *Chat) prepareChat(ctx context.Context, req *ChatRequest) (*model.StrategyRule, []*schema.Message, error) {
	logger.ContextInfo(ctx, "Chat.prepareChat", zap.String("model", req.Model), zap.String("dialogue_id", req.DialogueID))

	mapping, err := s.mappingRepo.FindByModelType(ctx, req.Model)
	if err != nil {
		logger.ContextError(ctx, "Chat.prepareChat", zap.String("model", req.Model), zap.Error(err))
		return nil, nil, fmt.Errorf("get model mapping: %w", err)
	}
	if mapping == nil {
		logger.ContextError(ctx, "Chat.prepareChat", zap.String("model", req.Model))
		return nil, nil, fmt.Errorf("model mapping not found: %s", req.Model)
	}

	strategy, err := s.strategy.GetByManufacturer(ctx, mapping.Manufacturer)
	if err != nil {
		logger.ContextError(ctx, "Chat.prepareChat", zap.String("manufacturer", mapping.Manufacturer), zap.Error(err))
		return nil, nil, fmt.Errorf("get strategy: %w", err)
	}
	if strategy == nil {
		logger.ContextError(ctx, "Chat.prepareChat", zap.String("manufacturer", mapping.Manufacturer))
		return nil, nil, fmt.Errorf("strategy not found: %s", mapping.Manufacturer)
	}

	history, err := s.recordRepo.ListByDialogueID(ctx, req.DialogueID)
	if err != nil {
		logger.ContextError(ctx, "Chat.prepareChat", zap.String("dialogue_id", req.DialogueID), zap.Error(err))
		return nil, nil, fmt.Errorf("list history: %w", err)
	}

	messages := buildMessages(req.Question, history, req.FileData, req.FileName)
	return strategy, messages, nil
}

// StreamChat 执行流式聊天，增量内容通过 contentChan 返回，返回 token 用量
func (s *Chat) StreamChat(ctx context.Context, req *ChatRequest, contentChan chan<- string) (*schema.TokenUsage, error) {
	ctx, span := tracing.Start(ctx, "logic.Chat.StreamChat")
	defer span.End()

	defer close(contentChan)

	strategy, messages, err := s.prepareChat(ctx, req)
	if err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	usage, err := callStream(ctx, strategy, messages, contentChan)
	if err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}
	return usage, nil
}

// GenerateChat 非流式生成聊天内容，返回内容和 token 用量
func (s *Chat) GenerateChat(ctx context.Context, req *ChatRequest) (string, *schema.TokenUsage, error) {
	ctx, span := tracing.Start(ctx, "logic.Chat.GenerateChat")
	defer span.End()

	strategy, messages, err := s.prepareChat(ctx, req)
	if err != nil {
		tracing.RecordError(ctx, err)
		return "", nil, err
	}

	cm, err := newChatModel(ctx, strategy)
	if err != nil {
		tracing.RecordError(ctx, err)
		return "", nil, err
	}

	msg, err := cm.Generate(ctx, messages)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.GenerateChat", zap.Error(err))
		return "", nil, fmt.Errorf("chat model generate: %w", err)
	}
	return msg.Content, msg.ResponseMeta.Usage, nil
}

// GetHistory 获取对话历史记录
func (s *Chat) GetHistory(ctx context.Context, dialogueID string) ([]*model.DialogueRecord, error) {
	records, err := s.recordRepo.ListByDialogueID(ctx, dialogueID)
	if err != nil {
		logger.ContextError(ctx, "Chat.GetHistory", zap.String("dialogue_id", dialogueID), zap.Error(err))
		return nil, err
	}
	return records, nil
}

// ListDialogues 获取对话摘要列表
func (s *Chat) ListDialogues(ctx context.Context, userID string) ([]*model.DialogueSummary, error) {
	ctx, span := tracing.Start(ctx, "logic.Chat.ListDialogues")
	defer span.End()

	summaries, err := s.recordRepo.ListDialogues(ctx, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.ListDialogues", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	return summaries, nil
}

// SaveRecord 持久化对话记录，usage 来自 Eino 返回的 token 用量
func (s *Chat) SaveRecord(ctx context.Context, req *ChatRequest, answer string, usage *schema.TokenUsage) error {
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
		record.UserToken = int64(usage.PromptTokens)
		record.AgentToken = int64(usage.CompletionTokens)
		record.TotalToken = int64(usage.TotalTokens)
	}
	if err := s.recordRepo.Create(ctx, record); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.SaveRecord", zap.String("record_id", record.RecordID), zap.Error(err))
		return err
	}
	return nil
}

// buildMessages 组装 Eino schema.Message 列表
func buildMessages(question string, history []*model.DialogueRecord, fileData []byte, fileName string) []*schema.Message {
	var messages []*schema.Message

	if len(fileData) > 0 && fileName != "" {
		messages = append(messages, schema.SystemMessage(documentPromptPrefix+string(fileData)))
	}

	for _, record := range history {
		messages = append(messages, schema.UserMessage(record.UserContent), schema.AssistantMessage(record.AgentContent, nil))
	}

	messages = append(messages, schema.UserMessage(question))
	return messages
}

// callStream 调用 Eino ChatModel 流式接口并返回结果和 token 用量
func callStream(ctx context.Context, rule *model.StrategyRule, messages []*schema.Message, contentChan chan<- string) (*schema.TokenUsage, error) {
	cm, err := newChatModel(ctx, rule)
	if err != nil {
		return nil, err
	}

	reader, err := cm.Stream(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("chat model stream: %w", err)
	}
	defer reader.Close()

	var usage *schema.TokenUsage
	for {
		chunk, err := reader.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("stream recv: %w", err)
		}
		if chunk.ResponseMeta != nil && chunk.ResponseMeta.Usage != nil {
			usage = chunk.ResponseMeta.Usage
		}
		if chunk.Content != "" {
			select {
			case contentChan <- chunk.Content:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
	return usage, nil
}
