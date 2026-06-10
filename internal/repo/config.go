package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"github.com/ethereal3x/mint-server/internal/crypto"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ModelConfigRepo 模型配置 GORM 实现
type ModelConfigRepo struct {
	db        *gorm.DB
	secretKey []byte
}

// NewModelConfigRepo 创建模型配置仓储
func NewModelConfigRepo(db *gorm.DB, secretKey []byte) *ModelConfigRepo {
	return &ModelConfigRepo{db: db, secretKey: secretKey}
}

func (r *ModelConfigRepo) FindByModelType(ctx context.Context, modelType string) (*model.ChatModelConfig, error) {
	var config model.ChatModelConfig
	if err := r.db.WithContext(ctx).Where("model_type = ? AND is_enabled = 1", modelType).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "ModelConfigRepo.FindByModelType", zap.String("model_type", modelType), zap.Error(err))
		return nil, fmt.Errorf("find config by model_type: %w", err)
	}
	decrypted, err := crypto.DecryptAPIKey(config.APIKey, r.secretKey)
	if err != nil {
		logger.ContextError(ctx, "ModelConfigRepo.FindByModelType", zap.String("model_type", modelType), zap.Error(err))
		return nil, fmt.Errorf("decrypt api key: %w", err)
	}
	config.APIKey = decrypted
	return &config, nil
}

func (r *ModelConfigRepo) ListAll(ctx context.Context) ([]*model.ChatModelConfig, error) {
	var list []*model.ChatModelConfig
	if err := r.db.WithContext(ctx).Where("is_enabled = 1").Order("id ASC").Find(&list).Error; err != nil {
		logger.ContextError(ctx, "ModelConfigRepo.ListAll", zap.Error(err))
		return nil, fmt.Errorf("list all configs: %w", err)
	}
	return list, nil
}

func (r *ModelConfigRepo) List(ctx context.Context, page, pageSize int32) ([]*model.ChatModelConfig, int64, error) {
	ctx, span := tracing.Start(ctx, "repo.ModelConfigRepo.List")
	defer span.End()

	var total int64
	query := r.db.WithContext(ctx).Model(&model.ChatModelConfig{})
	if err := query.Count(&total).Error; err != nil {
		tracing.RecordError(ctx, err)
		return nil, 0, fmt.Errorf("count configs: %w", err)
	}

	var list []*model.ChatModelConfig
	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("id ASC").Find(&list).Error; err != nil {
		tracing.RecordError(ctx, err)
		return nil, 0, fmt.Errorf("list configs: %w", err)
	}
	return list, total, nil
}

func (r *ModelConfigRepo) FindByID(ctx context.Context, id int32) (*model.ChatModelConfig, error) {
	var config model.ChatModelConfig
	if err := r.db.WithContext(ctx).First(&config, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "ModelConfigRepo.FindByID", zap.Int32("id", id), zap.Error(err))
		return nil, fmt.Errorf("find config by id: %w", err)
	}
	return &config, nil
}

// Create 创建配置，API Key 加密入库
func (r *ModelConfigRepo) Create(ctx context.Context, config *model.ChatModelConfig) error {
	ctx, span := tracing.Start(ctx, "repo.ModelConfigRepo.Create")
	defer span.End()

	encrypted, err := crypto.EncryptAPIKey(config.APIKey, r.secretKey)
	if err != nil {
		tracing.RecordError(ctx, err)
		return fmt.Errorf("encrypt api key: %w", err)
	}
	config.APIKey = encrypted

	if err := r.db.WithContext(ctx).Create(config).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "ModelConfigRepo.Create", zap.String("model_type", config.ModelType), zap.Error(err))
		return err
	}
	return nil
}

// Update 更新配置，API Key 加密入库（空值不更新）
func (r *ModelConfigRepo) Update(ctx context.Context, config *model.ChatModelConfig) error {
	ctx, span := tracing.Start(ctx, "repo.ModelConfigRepo.Update")
	defer span.End()

	omitFields := []string{"created_time"}
	if config.APIKey == "" {
		omitFields = append(omitFields, "api_key")
	} else {
		var err error
		config.APIKey, err = crypto.EncryptAPIKey(config.APIKey, r.secretKey)
		if err != nil {
			tracing.RecordError(ctx, err)
			return fmt.Errorf("encrypt api key: %w", err)
		}
	}

	if err := r.db.WithContext(ctx).Model(&model.ChatModelConfig{}).Where("id = ?", config.ID).Select("*").Omit(omitFields...).Updates(config).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "ModelConfigRepo.Update", zap.Int32("id", config.ID), zap.Error(err))
		return err
	}
	return nil
}

func (r *ModelConfigRepo) Delete(ctx context.Context, id int32) error {
	ctx, span := tracing.Start(ctx, "repo.ModelConfigRepo.Delete")
	defer span.End()

	if err := r.db.WithContext(ctx).Delete(&model.ChatModelConfig{}, id).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "ModelConfigRepo.Delete", zap.Int32("id", id), zap.Error(err))
		return err
	}
	return nil
}
