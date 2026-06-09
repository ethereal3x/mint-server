package logic

import (
	"context"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
)

// MappingRepo 模型映射数据访问接口
type MappingRepo interface {
	FindByID(ctx context.Context, id int32) (*model.ModelMapping, error)
	FindByModelType(ctx context.Context, modelType string) (*model.ModelMapping, error)
	ListByManufacturer(ctx context.Context, manufacturer string) ([]*model.ModelMapping, error)
	ListAll(ctx context.Context) ([]*model.ModelMapping, error)
	Create(ctx context.Context, mapping *model.ModelMapping) error
	Update(ctx context.Context, mapping *model.ModelMapping) error
	Delete(ctx context.Context, id int32) error
}

// Mapping 模型映射业务逻辑
type Mapping struct {
	repo MappingRepo
}

// NewMapping 创建模型映射业务逻辑
func NewMapping(repo MappingRepo) *Mapping {
	return &Mapping{repo: repo}
}

// GetByID 按 ID 获取模型映射
func (s *Mapping) GetByID(ctx context.Context, id int32) (*model.ModelMapping, error) {
	mapping, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logger.ContextError(ctx, "Mapping.GetByID", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}
	return mapping, nil
}

// GetByModelType 按模型类型获取映射
func (s *Mapping) GetByModelType(ctx context.Context, modelType string) (*model.ModelMapping, error) {
	mapping, err := s.repo.FindByModelType(ctx, modelType)
	if err != nil {
		logger.ContextError(ctx, "Mapping.GetByModelType", zap.String("model_type", modelType), zap.Error(err))
		return nil, err
	}
	return mapping, nil
}

// ListByManufacturer 按厂商获取映射列表
func (s *Mapping) ListByManufacturer(ctx context.Context, manufacturer string) ([]*model.ModelMapping, error) {
	list, err := s.repo.ListByManufacturer(ctx, manufacturer)
	if err != nil {
		logger.ContextError(ctx, "Mapping.ListByManufacturer", zap.String("manufacturer", manufacturer), zap.Error(err))
		return nil, err
	}
	return list, nil
}

// ListAll 获取所有映射
func (s *Mapping) ListAll(ctx context.Context) ([]*model.ModelMapping, error) {
	list, err := s.repo.ListAll(ctx)
	if err != nil {
		logger.ContextError(ctx, "Mapping.ListAll", zap.Error(err))
		return nil, err
	}
	return list, nil
}

// Create 创建模型映射
func (s *Mapping) Create(ctx context.Context, mapping *model.ModelMapping) error {
	ctx, span := tracing.Start(ctx, "logic.Mapping.Create")
	defer span.End()

	if err := s.repo.Create(ctx, mapping); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Mapping.Create", zap.String("model_type", mapping.ModelType), zap.Error(err))
		return err
	}
	return nil
}

// Update 更新模型映射
func (s *Mapping) Update(ctx context.Context, mapping *model.ModelMapping) error {
	ctx, span := tracing.Start(ctx, "logic.Mapping.Update")
	defer span.End()

	if err := s.repo.Update(ctx, mapping); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Mapping.Update", zap.Int32("id", mapping.ID), zap.Error(err))
		return err
	}
	return nil
}

// Delete 删除模型映射
func (s *Mapping) Delete(ctx context.Context, id int32) error {
	ctx, span := tracing.Start(ctx, "logic.Mapping.Delete")
	defer span.End()

	if err := s.repo.Delete(ctx, id); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Mapping.Delete", zap.Int32("id", id), zap.Error(err))
		return err
	}
	return nil
}
