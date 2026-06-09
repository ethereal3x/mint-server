package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MappingRepo 模型映射 GORM 实现
type MappingRepo struct {
	db *gorm.DB
}

// NewMappingRepo 创建模型映射仓储
func NewMappingRepo(db *gorm.DB) *MappingRepo {
	return &MappingRepo{db: db}
}

func (r *MappingRepo) FindByID(ctx context.Context, id int32) (*model.ModelMapping, error) {
	var mapping model.ModelMapping
	if err := r.db.WithContext(ctx).First(&mapping, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "MappingRepo.FindByID", zap.Int32("id", id), zap.Error(err))
		return nil, fmt.Errorf("find mapping by id: %w", err)
	}
	return &mapping, nil
}

func (r *MappingRepo) FindByModelType(ctx context.Context, modelType string) (*model.ModelMapping, error) {
	var mapping model.ModelMapping
	if err := r.db.WithContext(ctx).Where("model_type = ?", modelType).First(&mapping).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "MappingRepo.FindByModelType", zap.String("model_type", modelType), zap.Error(err))
		return nil, fmt.Errorf("find mapping by model type: %w", err)
	}
	return &mapping, nil
}

func (r *MappingRepo) ListByManufacturer(ctx context.Context, manufacturer string) ([]*model.ModelMapping, error) {
	var list []*model.ModelMapping
	if err := r.db.WithContext(ctx).Where("manufacturer = ?", manufacturer).Find(&list).Error; err != nil {
		logger.ContextError(ctx, "MappingRepo.ListByManufacturer", zap.String("manufacturer", manufacturer), zap.Error(err))
		return nil, fmt.Errorf("list mappings by manufacturer: %w", err)
	}
	return list, nil
}

func (r *MappingRepo) ListAll(ctx context.Context) ([]*model.ModelMapping, error) {
	var list []*model.ModelMapping
	if err := r.db.WithContext(ctx).Find(&list).Error; err != nil {
		logger.ContextError(ctx, "MappingRepo.ListAll", zap.Error(err))
		return nil, fmt.Errorf("list all mappings: %w", err)
	}
	return list, nil
}

func (r *MappingRepo) Create(ctx context.Context, mapping *model.ModelMapping) error {
	ctx, span := tracing.Start(ctx, "repo.MappingRepo.Create")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(mapping).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "MappingRepo.Create", zap.String("model_type", mapping.ModelType), zap.Error(err))
		return err
	}
	return nil
}

func (r *MappingRepo) Update(ctx context.Context, mapping *model.ModelMapping) error {
	ctx, span := tracing.Start(ctx, "repo.MappingRepo.Update")
	defer span.End()

	if err := r.db.WithContext(ctx).Save(mapping).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "MappingRepo.Update", zap.Int32("id", mapping.ID), zap.Error(err))
		return err
	}
	return nil
}

func (r *MappingRepo) Delete(ctx context.Context, id int32) error {
	ctx, span := tracing.Start(ctx, "repo.MappingRepo.Delete")
	defer span.End()

	if err := r.db.WithContext(ctx).Delete(&model.ModelMapping{}, id).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "MappingRepo.Delete", zap.Int32("id", id), zap.Error(err))
		return err
	}
	return nil
}
