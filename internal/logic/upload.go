package logic

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/storage"
	"github.com/ethereal3x/apc/tracing"
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
	FindByID(ctx context.Context, id int32) (*model.FileUpload, error)
	ListByUserID(ctx context.Context, userID string) ([]*model.FileUpload, error)
}

// UploadLogic 文件上传业务逻辑
type UploadLogic struct {
	storage storage.ObjectStorage
	repo    UploadRepo
}

// UploadResult 上传结果
type UploadResult struct {
	ID       int32  `json:"id"`
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

// NewUploadLogic 创建文件上传业务逻辑
func NewUploadLogic(store storage.ObjectStorage, repo UploadRepo) *UploadLogic {
	return &UploadLogic{storage: store, repo: repo}
}

// Upload 简单文件上传：接收文件流 → 上传 rustfs → 记录 DB → 返回 URL
func (l *UploadLogic) Upload(ctx context.Context, userID string, fileName string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	ctx, span := tracing.Start(ctx, "logic.UploadLogic.Upload")
	defer span.End()

	if size > maxUploadSize {
		return nil, fmt.Errorf("file size %d exceeds max %d", size, maxUploadSize)
	}

	objectName := generateObjectName(userID, fileName)
	record := &model.FileUpload{
		ObjectName:   objectName,
		OriginalName: fileName,
		FileSize:     size,
		ContentType:  contentType,
		Status:       model.UPLOAD_STATUS_UPLOADING,
		UserID:       userID,
	}
	if err := l.repo.Create(ctx, record); err != nil {
		logger.ContextError(ctx, "UploadLogic.Upload", zap.String("object_name", objectName), zap.Error(err))
		return nil, mint_err.ErrDBCreate
	}

	opts := storage.UploadOptions{ContentType: contentType}
	if err := l.storage.Upload(ctx, objectName, reader, size, opts); err != nil {
		record.Status = model.UPLOAD_STATUS_FAILED
		_ = l.repo.Update(ctx, record)
		logger.ContextError(ctx, "UploadLogic.Upload", zap.String("object_name", objectName), zap.Error(err))
		return nil, fmt.Errorf("storage upload: %w", err)
	}

	record.Status = model.UPLOAD_STATUS_COMPLETED
	record.URL = l.storage.PublicURL(objectName)
	record.UploadedSize = size
	_ = l.repo.Update(ctx, record)

	return &UploadResult{ID: record.ID, URL: record.URL, FileName: fileName, FileSize: size}, nil
}

// GetUpload 查询上传记录
func (l *UploadLogic) GetUpload(ctx context.Context, id int32) (*model.FileUpload, error) {
	record, err := l.repo.FindByID(ctx, id)
	if err != nil {
		logger.ContextError(ctx, "UploadLogic.GetUpload", zap.Int32("id", id), zap.Error(err))
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
