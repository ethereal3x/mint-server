package decorator

import (
	"context"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"github.com/ethereal3x/mint-server/internal/dto"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
)

// chatLogic ChatDecorator 需要的聊天业务接口
type chatLogic interface {
	StreamChat(ctx context.Context, req *dto.ChatRequest, contentChan chan<- string) (*dto.ChatResult, error)
	GenerateChat(ctx context.Context, req *dto.ChatRequest) (*dto.ChatResult, error)
	SaveRecord(ctx context.Context, req *dto.SaveRecordRequest) error
	GetHistory(ctx context.Context, query *model.DialogueQuery) ([]*model.DialogueRecord, error)
	ListDialogues(ctx context.Context, userID string) ([]*model.DialogueSummary, error)
	AggregateByModel(ctx context.Context, userID string) ([]*model.ModelStat, error)
}

// ChatDecorator 为 Chat 业务逻辑添加 tracing 和 logging 横切关注点
type ChatDecorator struct {
	inner chatLogic
}

// NewChatDecorator 创建聊天逻辑装饰器
func NewChatDecorator(inner chatLogic) *ChatDecorator {
	return &ChatDecorator{inner: inner}
}

// StreamChat 执行流式聊天并装饰 tracing/logging
func (d *ChatDecorator) StreamChat(ctx context.Context, req *dto.ChatRequest, contentChan chan<- string) (*dto.ChatResult, error) {
	ctx, span := tracing.Start(ctx, "Chat.StreamChat")
	defer span.End()
	logger.ContextInfo(ctx, "Chat.StreamChat", zap.String("model", req.Model), zap.String("dialogue_id", req.DialogueID))
	result, err := d.inner.StreamChat(ctx, req, contentChan)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.StreamChat", zap.Error(err))
	}
	return result, err
}

// GenerateChat 执行非流式聊天并装饰 tracing/logging
func (d *ChatDecorator) GenerateChat(ctx context.Context, req *dto.ChatRequest) (*dto.ChatResult, error) {
	ctx, span := tracing.Start(ctx, "Chat.GenerateChat")
	defer span.End()
	logger.ContextInfo(ctx, "Chat.GenerateChat", zap.String("model", req.Model))
	result, err := d.inner.GenerateChat(ctx, req)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.GenerateChat", zap.Error(err))
	}
	return result, err
}

// SaveRecord 持久化对话记录并装饰 tracing/logging
func (d *ChatDecorator) SaveRecord(ctx context.Context, req *dto.SaveRecordRequest) error {
	ctx, span := tracing.Start(ctx, "Chat.SaveRecord")
	defer span.End()
	err := d.inner.SaveRecord(ctx, req)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.SaveRecord", zap.Error(err))
	}
	return err
}

// GetHistory 获取对话历史并装饰 tracing/logging
func (d *ChatDecorator) GetHistory(ctx context.Context, query *model.DialogueQuery) ([]*model.DialogueRecord, error) {
	ctx, span := tracing.Start(ctx, "Chat.GetHistory")
	defer span.End()
	records, err := d.inner.GetHistory(ctx, query)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.GetHistory", zap.Error(err))
	}
	return records, err
}

// ListDialogues 获取对话列表并装饰 tracing/logging
func (d *ChatDecorator) ListDialogues(ctx context.Context, userID string) ([]*model.DialogueSummary, error) {
	ctx, span := tracing.Start(ctx, "Chat.ListDialogues")
	defer span.End()
	summaries, err := d.inner.ListDialogues(ctx, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.ListDialogues", zap.Error(err))
	}
	return summaries, err
}

// AggregateByModel 按模型聚合统计并装饰 tracing/logging
func (d *ChatDecorator) AggregateByModel(ctx context.Context, userID string) ([]*model.ModelStat, error) {
	ctx, span := tracing.Start(ctx, "Chat.AggregateByModel")
	defer span.End()
	stats, err := d.inner.AggregateByModel(ctx, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.AggregateByModel", zap.Error(err))
	}
	return stats, err
}

// configLogic ConfigDecorator 需要的配置业务接口
type configLogic interface {
	ListAll(ctx context.Context, userID string) ([]*model.ChatModelConfig, error)
	List(ctx context.Context, page, pageSize int32, userID string) ([]*model.ChatModelConfig, int64, error)
	GetByID(ctx context.Context, id int32, userID string) (*model.ChatModelConfig, error)
	Create(ctx context.Context, config *model.ChatModelConfig) error
	Update(ctx context.Context, config *model.ChatModelConfig, userID string) error
	Delete(ctx context.Context, id int32, userID string) error
}

// ConfigDecorator 为 Config 业务逻辑添加 tracing 和 logging 横切关注点
type ConfigDecorator struct {
	inner configLogic
}

// NewConfigDecorator 创建配置逻辑装饰器
func NewConfigDecorator(inner configLogic) *ConfigDecorator {
	return &ConfigDecorator{inner: inner}
}

// ListAll 获取全部模型配置并装饰 tracing/logging
func (d *ConfigDecorator) ListAll(ctx context.Context, userID string) ([]*model.ChatModelConfig, error) {
	ctx, span := tracing.Start(ctx, "Config.ListAll")
	defer span.End()
	configs, err := d.inner.ListAll(ctx, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Config.ListAll", zap.Error(err))
	}
	return configs, err
}

// List 分页查询配置并装饰 tracing/logging
func (d *ConfigDecorator) List(ctx context.Context, page, pageSize int32, userID string) ([]*model.ChatModelConfig, int64, error) {
	ctx, span := tracing.Start(ctx, "Config.List")
	defer span.End()
	configs, total, err := d.inner.List(ctx, page, pageSize, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Config.List", zap.Error(err))
	}
	return configs, total, err
}

// GetByID 按 ID 查询配置并装饰 tracing/logging
func (d *ConfigDecorator) GetByID(ctx context.Context, id int32, userID string) (*model.ChatModelConfig, error) {
	ctx, span := tracing.Start(ctx, "Config.GetByID")
	defer span.End()
	config, err := d.inner.GetByID(ctx, id, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Config.GetByID", zap.Error(err))
	}
	return config, err
}

// Create 创建配置并装饰 tracing/logging
func (d *ConfigDecorator) Create(ctx context.Context, config *model.ChatModelConfig) error {
	ctx, span := tracing.Start(ctx, "Config.Create")
	defer span.End()
	err := d.inner.Create(ctx, config)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Config.Create", zap.Error(err))
	}
	return err
}

// Update 更新配置并装饰 tracing/logging
func (d *ConfigDecorator) Update(ctx context.Context, config *model.ChatModelConfig, userID string) error {
	ctx, span := tracing.Start(ctx, "Config.Update")
	defer span.End()
	err := d.inner.Update(ctx, config, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Config.Update", zap.Error(err))
	}
	return err
}

// Delete 删除配置并装饰 tracing/logging
func (d *ConfigDecorator) Delete(ctx context.Context, id int32, userID string) error {
	ctx, span := tracing.Start(ctx, "Config.Delete")
	defer span.End()
	err := d.inner.Delete(ctx, id, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Config.Delete", zap.Error(err))
	}
	return err
}
