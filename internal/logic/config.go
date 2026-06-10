package logic

import (
	"context"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
)

// ConfigRepo 模型配置数据访问接口
type ConfigRepo interface {
	FindByModelType(ctx context.Context, modelType string) (*model.ChatModelConfig, error)
	FindByID(ctx context.Context, id int32) (*model.ChatModelConfig, error)
	ListAll(ctx context.Context) ([]*model.ChatModelConfig, error)
	List(ctx context.Context, page, pageSize int32) ([]*model.ChatModelConfig, int64, error)
	Create(ctx context.Context, config *model.ChatModelConfig) error
	Update(ctx context.Context, config *model.ChatModelConfig) error
	Delete(ctx context.Context, id int32) error
}

// Config 模型配置业务逻辑
type Config struct {
	repo ConfigRepo
}

// NewConfig 创建模型配置业务逻辑
func NewConfig(repo ConfigRepo) *Config {
	return &Config{repo: repo}
}

// GetByModelType 按模型标识获取配置
func (s *Config) GetByModelType(ctx context.Context, modelType string) (*model.ChatModelConfig, error) {
	config, err := s.repo.FindByModelType(ctx, modelType)
	if err != nil {
		logger.ContextError(ctx, "Config.GetByModelType", zap.String("model_type", modelType), zap.Error(err))
		return nil, mint_err.ErrDBQuery
	}
	return config, nil
}

// GetByID 按 ID 获取配置
func (s *Config) GetByID(ctx context.Context, id int32) (*model.ChatModelConfig, error) {
	config, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logger.ContextError(ctx, "Config.GetByID", zap.Int32("id", id), zap.Error(err))
		return nil, mint_err.ErrDBQuery
	}
	return config, nil
}

// ListAll 获取所有配置
func (s *Config) ListAll(ctx context.Context) ([]*model.ChatModelConfig, error) {
	list, err := s.repo.ListAll(ctx)
	if err != nil {
		logger.ContextError(ctx, "Config.ListAll", zap.Error(err))
		return nil, mint_err.ErrListModels
	}
	return list, nil
}

// List 分页获取配置
func (s *Config) List(ctx context.Context, page, pageSize int32) ([]*model.ChatModelConfig, int64, error) {
	list, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		logger.ContextError(ctx, "Config.List", zap.Int32("page", page), zap.Error(err))
		return nil, 0, mint_err.ErrListConfigs
	}
	return list, total, nil
}

// Create 创建模型配置
func (s *Config) Create(ctx context.Context, config *model.ChatModelConfig) error {
	ctx, span := tracing.Start(ctx, "logic.Config.Create")
	defer span.End()

	if err := s.repo.Create(ctx, config); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Config.Create", zap.String("model_type", config.ModelType), zap.Error(err))
		return mint_err.ErrCreateConfig
	}
	return nil
}

// Update 更新模型配置
func (s *Config) Update(ctx context.Context, config *model.ChatModelConfig) error {
	ctx, span := tracing.Start(ctx, "logic.Config.Update")
	defer span.End()

	if err := s.repo.Update(ctx, config); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Config.Update", zap.Int32("id", config.ID), zap.Error(err))
		return mint_err.ErrUpdateConfig
	}
	return nil
}

// Delete 删除模型配置
func (s *Config) Delete(ctx context.Context, id int32) error {
	ctx, span := tracing.Start(ctx, "logic.Config.Delete")
	defer span.End()

	if err := s.repo.Delete(ctx, id); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Config.Delete", zap.Int32("id", id), zap.Error(err))
		return mint_err.ErrDeleteConfig
	}
	return nil
}
