package app

import (
	"github.com/ethereal3x/apc/storage"
	"github.com/ethereal3x/mint-server/internal/auth"
	"github.com/ethereal3x/mint-server/internal/decorator"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/repo"
	"github.com/ethereal3x/mint-server/internal/service"
	"gorm.io/gorm"
)

// UploadModule 聚合文件上传的 HTTP 处理器
type UploadModule struct {
	UploadHandler *service.UploadHandler
}

// newUploadModule 装配 Upload 领域全部依赖，通过 decorator 注入 tracing/logging
func newUploadModule(db *gorm.DB, store storage.ObjectStorage, tokenManager *auth.TokenManager) *UploadModule {
	uploadRepo := repo.NewFileUploadRepo(db)
	uploadLogic := logic.NewUploadLogic(store, uploadRepo)
	uploadDecorated := decorator.NewUploadDecorator(uploadLogic)

	return &UploadModule{
		UploadHandler: service.NewUploadHandler(uploadDecorated, tokenManager),
	}
}
