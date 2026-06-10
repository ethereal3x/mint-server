package decorator

import (
	"context"

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
	return wrap(ctx, "Chat.StreamChat", func(ctx context.Context) (*dto.ChatResult, error) {
		return d.inner.StreamChat(ctx, req, contentChan)
	}, zap.String("model", req.Model), zap.String("dialogue_id", req.DialogueID))
}

// GenerateChat 执行非流式聊天并装饰 tracing/logging
func (d *ChatDecorator) GenerateChat(ctx context.Context, req *dto.ChatRequest) (*dto.ChatResult, error) {
	return wrap(ctx, "Chat.GenerateChat", func(ctx context.Context) (*dto.ChatResult, error) {
		return d.inner.GenerateChat(ctx, req)
	}, zap.String("model", req.Model))
}

// SaveRecord 持久化对话记录并装饰 tracing/logging
func (d *ChatDecorator) SaveRecord(ctx context.Context, req *dto.SaveRecordRequest) error {
	return wrapErr(ctx, "Chat.SaveRecord", func(ctx context.Context) error {
		return d.inner.SaveRecord(ctx, req)
	})
}

// GetHistory 获取对话历史并装饰 tracing/logging
func (d *ChatDecorator) GetHistory(ctx context.Context, query *model.DialogueQuery) ([]*model.DialogueRecord, error) {
	return wrap(ctx, "Chat.GetHistory", func(ctx context.Context) ([]*model.DialogueRecord, error) {
		return d.inner.GetHistory(ctx, query)
	})
}

// ListDialogues 获取对话列表并装饰 tracing/logging
func (d *ChatDecorator) ListDialogues(ctx context.Context, userID string) ([]*model.DialogueSummary, error) {
	return wrap(ctx, "Chat.ListDialogues", func(ctx context.Context) ([]*model.DialogueSummary, error) {
		return d.inner.ListDialogues(ctx, userID)
	})
}

// AggregateByModel 按模型聚合统计并装饰 tracing/logging
func (d *ChatDecorator) AggregateByModel(ctx context.Context, userID string) ([]*model.ModelStat, error) {
	return wrap(ctx, "Chat.AggregateByModel", func(ctx context.Context) ([]*model.ModelStat, error) {
		return d.inner.AggregateByModel(ctx, userID)
	})
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
	return wrap(ctx, "Config.ListAll", func(ctx context.Context) ([]*model.ChatModelConfig, error) {
		return d.inner.ListAll(ctx, userID)
	})
}

// List 分页查询配置并装饰 tracing/logging
func (d *ConfigDecorator) List(ctx context.Context, page, pageSize int32, userID string) ([]*model.ChatModelConfig, int64, error) {
	ctx, span := tracing.Start(ctx, "Config.List")
	defer span.End()
	configs, total, err := d.inner.List(ctx, page, pageSize, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
	}
	return configs, total, err
}

// GetByID 按 ID 查询配置并装饰 tracing/logging
func (d *ConfigDecorator) GetByID(ctx context.Context, id int32, userID string) (*model.ChatModelConfig, error) {
	return wrap(ctx, "Config.GetByID", func(ctx context.Context) (*model.ChatModelConfig, error) {
		return d.inner.GetByID(ctx, id, userID)
	})
}

// Create 创建配置并装饰 tracing/logging
func (d *ConfigDecorator) Create(ctx context.Context, config *model.ChatModelConfig) error {
	return wrapErr(ctx, "Config.Create", func(ctx context.Context) error {
		return d.inner.Create(ctx, config)
	})
}

// Update 更新配置并装饰 tracing/logging
func (d *ConfigDecorator) Update(ctx context.Context, config *model.ChatModelConfig, userID string) error {
	return wrapErr(ctx, "Config.Update", func(ctx context.Context) error {
		return d.inner.Update(ctx, config, userID)
	})
}

// Delete 删除配置并装饰 tracing/logging
func (d *ConfigDecorator) Delete(ctx context.Context, id int32, userID string) error {
	return wrapErr(ctx, "Config.Delete", func(ctx context.Context) error {
		return d.inner.Delete(ctx, id, userID)
	})
}
