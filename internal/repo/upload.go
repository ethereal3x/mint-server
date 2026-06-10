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

// FileUploadRepo 文件上传记录 GORM 实现
type FileUploadRepo struct {
	db *gorm.DB
}

// NewFileUploadRepo 创建文件上传记录仓储
func NewFileUploadRepo(db *gorm.DB) *FileUploadRepo {
	return &FileUploadRepo{db: db}
}

func (r *FileUploadRepo) FindByID(ctx context.Context, id int32) (*model.FileUpload, error) {
	var record model.FileUpload
	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "FileUploadRepo.FindByID", zap.Int32("id", id), zap.Error(err))
		return nil, fmt.Errorf("find upload by id: %w", err)
	}
	return &record, nil
}

func (r *FileUploadRepo) FindByObjectName(ctx context.Context, objectName string) (*model.FileUpload, error) {
	var record model.FileUpload
	if err := r.db.WithContext(ctx).Where("object_name = ?", objectName).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "FileUploadRepo.FindByObjectName", zap.String("object_name", objectName), zap.Error(err))
		return nil, fmt.Errorf("find upload by object_name: %w", err)
	}
	return &record, nil
}

func (r *FileUploadRepo) ListByUserID(ctx context.Context, userID string) ([]*model.FileUpload, error) {
	var list []*model.FileUpload
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_time DESC").Limit(50).Find(&list).Error; err != nil {
		logger.ContextError(ctx, "FileUploadRepo.ListByUserID", zap.String("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("list uploads: %w", err)
	}
	return list, nil
}

func (r *FileUploadRepo) Create(ctx context.Context, record *model.FileUpload) error {
	ctx, span := tracing.Start(ctx, "repo.FileUploadRepo.Create")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "FileUploadRepo.Create", zap.String("object_name", record.ObjectName), zap.Error(err))
		return err
	}
	return nil
}

func (r *FileUploadRepo) Update(ctx context.Context, record *model.FileUpload) error {
	ctx, span := tracing.Start(ctx, "repo.FileUploadRepo.Update")
	defer span.End()

	if err := r.db.WithContext(ctx).Model(&model.FileUpload{}).Where("id = ?", record.ID).
		Select("*").Omit("created_time").Updates(record).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "FileUploadRepo.Update", zap.Int32("id", record.ID), zap.Error(err))
		return err
	}
	return nil
}
