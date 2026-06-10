package decorator

import (
	"context"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"github.com/ethereal3x/mint-server/internal/dto"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
)

// uploadLogic UploadDecorator 需要的上传业务接口
type uploadLogic interface {
	Upload(ctx context.Context, req *dto.UploadRequest) (*dto.UploadResult, error)
	GetUpload(ctx context.Context, query *model.FileUploadQuery) (*model.FileUpload, error)
	ListUploads(ctx context.Context, userID string) ([]*model.FileUpload, error)
}

// UploadDecorator 为文件上传业务逻辑添加 tracing 和 logging 横切关注点
type UploadDecorator struct {
	inner uploadLogic
}

// NewUploadDecorator 创建上传逻辑装饰器
func NewUploadDecorator(inner uploadLogic) *UploadDecorator {
	return &UploadDecorator{inner: inner}
}

// Upload 文件上传并装饰 tracing/logging
func (d *UploadDecorator) Upload(ctx context.Context, req *dto.UploadRequest) (*dto.UploadResult, error) {
	ctx, span := tracing.Start(ctx, "Upload.Upload")
	defer span.End()
	logger.ContextInfo(ctx, "Upload.Upload", zap.String("filename", req.FileName))
	result, err := d.inner.Upload(ctx, req)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Upload.Upload", zap.Error(err))
	}
	return result, err
}

// GetUpload 查询上传记录并装饰 tracing/logging
func (d *UploadDecorator) GetUpload(ctx context.Context, query *model.FileUploadQuery) (*model.FileUpload, error) {
	ctx, span := tracing.Start(ctx, "Upload.GetUpload")
	defer span.End()
	record, err := d.inner.GetUpload(ctx, query)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Upload.GetUpload", zap.Error(err))
	}
	return record, err
}

// ListUploads 查询上传记录列表并装饰 tracing/logging
func (d *UploadDecorator) ListUploads(ctx context.Context, userID string) ([]*model.FileUpload, error) {
	ctx, span := tracing.Start(ctx, "Upload.ListUploads")
	defer span.End()
	list, err := d.inner.ListUploads(ctx, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Upload.ListUploads", zap.Error(err))
	}
	return list, err
}
