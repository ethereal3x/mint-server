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

// FindByModelTypeForUser 按用户和模型标识查询启用的配置
func (r *ModelConfigRepo) FindByModelTypeForUser(ctx context.Context, modelType string, userID string) (*model.ChatModelConfig, error) {
	var config model.ChatModelConfig
	if err := r.db.WithContext(ctx).Where("model_type = ? AND user_id = ? AND is_enabled = 1", modelType, userID).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "ModelConfigRepo.FindByModelTypeForUser", zap.String("model_type", modelType), zap.String("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("find config by model_type: %w", err)
	}
	decrypted, err := crypto.DecryptAPIKey(config.APIKey, r.secretKey)
	if err != nil {
		logger.ContextError(ctx, "ModelConfigRepo.FindByModelTypeForUser", zap.String("model_type", modelType), zap.Error(err))
		return nil, fmt.Errorf("decrypt api key: %w", err)
	}
	config.APIKey = decrypted
	return &config, nil
}

// ListAllForUser 获取指定用户所有启用的配置
func (r *ModelConfigRepo) ListAllForUser(ctx context.Context, userID string) ([]*model.ChatModelConfig, error) {
	var list []*model.ChatModelConfig
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_enabled = 1", userID).Order("id ASC").Find(&list).Error; err != nil {
		logger.ContextError(ctx, "ModelConfigRepo.ListAllForUser", zap.String("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("list all configs: %w", err)
	}
	return list, nil
}

// ListForUser 分页获取用户配置
func (r *ModelConfigRepo) ListForUser(ctx context.Context, page, pageSize int32, userID string) ([]*model.ChatModelConfig, int64, error) {
	ctx, span := tracing.Start(ctx, "repo.ModelConfigRepo.ListForUser")
	defer span.End()

	var total int64
	query := r.db.WithContext(ctx).Model(&model.ChatModelConfig{}).Where("user_id = ?", userID)
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

// FindByIDForUser 按 ID 和用户ID查询配置，返回解密后的 API Key
func (r *ModelConfigRepo) FindByIDForUser(ctx context.Context, id int32, userID string) (*model.ChatModelConfig, error) {
	var config model.ChatModelConfig
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "ModelConfigRepo.FindByIDForUser", zap.Int32("id", id), zap.String("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("find config by id: %w", err)
	}
	decrypted, err := crypto.DecryptAPIKey(config.APIKey, r.secretKey)
	if err != nil {
		logger.ContextError(ctx, "ModelConfigRepo.FindByIDForUser", zap.Int32("id", id), zap.Error(err))
		return nil, fmt.Errorf("decrypt api key: %w", err)
	}
	config.APIKey = decrypted
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
		logger.ContextError(ctx, "ModelConfigRepo.Create", zap.String("model_type", config.ModelType), zap.String("user_id", config.UserID), zap.Error(err))
		return err
	}
	return nil
}

// UpdateForUser 更新用户配置，API Key 加密入库（空值不更新）
func (r *ModelConfigRepo) UpdateForUser(ctx context.Context, config *model.ChatModelConfig, userID string) error {
	ctx, span := tracing.Start(ctx, "repo.ModelConfigRepo.UpdateForUser")
	defer span.End()

	omitFields := []string{"created_time", "user_id"}
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

	if err := r.db.WithContext(ctx).Model(&model.ChatModelConfig{}).Where("id = ? AND user_id = ?", config.ID, userID).Select("*").Omit(omitFields...).Updates(config).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "ModelConfigRepo.UpdateForUser", zap.Int32("id", config.ID), zap.String("user_id", userID), zap.Error(err))
		return err
	}
	return nil
}

// DeleteForUser 删除用户配置
func (r *ModelConfigRepo) DeleteForUser(ctx context.Context, id int32, userID string) error {
	ctx, span := tracing.Start(ctx, "repo.ModelConfigRepo.DeleteForUser")
	defer span.End()

	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.ChatModelConfig{}).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "ModelConfigRepo.DeleteForUser", zap.Int32("id", id), zap.String("user_id", userID), zap.Error(err))
		return err
	}
	return nil
}
