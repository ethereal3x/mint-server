package logic

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/storage"
	"github.com/ethereal3x/mint-server/internal/dto"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const maxUploadSize = 50 << 20 // 50MB

// UploadRepo 文件上传记录数据访问接口
type UploadRepo interface {
	Create(ctx context.Context, record *model.FileUpload) error
	Update(ctx context.Context, record *model.FileUpload) error
	FindByIDForUser(ctx context.Context, query *model.FileUploadQuery) (*model.FileUpload, error)
	ListByUserID(ctx context.Context, userID string) ([]*model.FileUpload, error)
}

// UploadLogic 文件上传业务逻辑
type UploadLogic struct {
	storage storage.ObjectStorage
	repo    UploadRepo
}

// NewUploadLogic 创建文件上传业务逻辑
func NewUploadLogic(store storage.ObjectStorage, repo UploadRepo) *UploadLogic {
	return &UploadLogic{storage: store, repo: repo}
}

// Upload 简单文件上传：接收文件流、上传对象存储、记录 DB 并返回 URL
func (l *UploadLogic) Upload(ctx context.Context, req *dto.UploadRequest) (*dto.UploadResult, error) {
	if req.Size > maxUploadSize {
		return nil, fmt.Errorf("file size %d exceeds max %d", req.Size, maxUploadSize)
	}

	objectName := generateObjectName(req.UserID, req.FileName)
	record := &model.FileUpload{
		ObjectName:   objectName,
		OriginalName: req.FileName,
		FileSize:     req.Size,
		ContentType:  req.ContentType,
		Status:       model.UPLOAD_STATUS_UPLOADING,
		UserID:       req.UserID,
	}
	if err := l.repo.Create(ctx, record); err != nil {
		logger.ContextError(ctx, "UploadLogic.Upload", zap.String("object_name", objectName), zap.Error(err))
		return nil, mint_err.ErrDBCreate
	}

	opts := storage.UploadOptions{ContentType: req.ContentType}
	if err := l.storage.Upload(ctx, objectName, req.Reader, req.Size, opts); err != nil {
		record.Status = model.UPLOAD_STATUS_FAILED
		if updateErr := l.repo.Update(ctx, record); updateErr != nil {
			logger.ContextError(ctx, "UploadLogic.Upload", zap.String("object_name", objectName), zap.Error(updateErr))
		}
		logger.ContextError(ctx, "UploadLogic.Upload", zap.String("object_name", objectName), zap.Error(err))
		return nil, fmt.Errorf("storage upload: %w", err)
	}

	record.Status = model.UPLOAD_STATUS_COMPLETED
	record.URL = l.storage.PublicURL(objectName)
	record.UploadedSize = req.Size
	if err := l.repo.Update(ctx, record); err != nil {
		logger.ContextError(ctx, "UploadLogic.Upload", zap.String("object_name", objectName), zap.Error(err))
		return nil, mint_err.ErrDBUpdate
	}

	return &dto.UploadResult{ID: record.ID, URL: record.URL, FileName: req.FileName, FileSize: req.Size}, nil
}

// GetUpload 查询当前用户的上传记录
func (l *UploadLogic) GetUpload(ctx context.Context, query *model.FileUploadQuery) (*model.FileUpload, error) {
	record, err := l.repo.FindByIDForUser(ctx, query)
	if err != nil {
		logger.ContextError(ctx, "UploadLogic.GetUpload", zap.Int32("id", query.ID), zap.String("user_id", query.UserID), zap.Error(err))
		return nil, mint_err.ErrDBQuery
	}
	return record, nil
}

// ListUploads 查询用户上传记录
func (l *UploadLogic) ListUploads(ctx context.Context, userID string) ([]*model.FileUpload, error) {
	list, err := l.repo.ListByUserID(ctx, userID)
	if err != nil {
		logger.ContextError(ctx, "UploadLogic.ListUploads", zap.String("user_id", userID), zap.Error(err))
		return nil, mint_err.ErrDBQuery
	}
	return list, nil
}

// generateObjectName 生成唯一对象名：uploads/{yyyyMMdd}/{uuid}_{filename}
func generateObjectName(userID string, fileName string) string {
	ext := filepath.Ext(fileName)
	dateDir := time.Now().Format("20060102")
	return fmt.Sprintf("uploads/%s/%s_%s%s", dateDir, userID, uuid.NewString()[:8], ext)
}
