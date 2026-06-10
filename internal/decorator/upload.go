package decorator

import (
	"context"

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
	return wrap(ctx, "Upload.Upload", func(ctx context.Context) (*dto.UploadResult, error) {
		return d.inner.Upload(ctx, req)
	}, zap.String("filename", req.FileName))
}

// GetUpload 查询上传记录并装饰 tracing/logging
func (d *UploadDecorator) GetUpload(ctx context.Context, query *model.FileUploadQuery) (*model.FileUpload, error) {
	return wrap(ctx, "Upload.GetUpload", func(ctx context.Context) (*model.FileUpload, error) {
		return d.inner.GetUpload(ctx, query)
	})
}

// ListUploads 查询上传记录列表并装饰 tracing/logging
func (d *UploadDecorator) ListUploads(ctx context.Context, userID string) ([]*model.FileUpload, error) {
	return wrap(ctx, "Upload.ListUploads", func(ctx context.Context) ([]*model.FileUpload, error) {
		return d.inner.ListUploads(ctx, userID)
	})
}
