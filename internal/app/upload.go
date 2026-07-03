package app

import (
	"github.com/ethereal3x/apc/storage"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/repo"
	"github.com/ethereal3x/mint-server/internal/service"
	"github.com/ethereal3x/mint-server/internal/util"
	"gorm.io/gorm"
)

// UploadModule 聚合文件上传的 HTTP 处理器
type UploadModule struct {
	UploadHandler *service.UploadHandler
}

// newUploadModule 装配 Upload 领域全部依赖
func newUploadModule(db *gorm.DB, store storage.ObjectStorage, tokenManager *util.TokenManager) *UploadModule {
	uploadRepo := repo.NewFileUploadRepo(db)
	uploadLogic := logic.NewUploadLogic(store, uploadRepo)
	return &UploadModule{
		UploadHandler: service.NewUploadHandler(uploadLogic, tokenManager),
	}
}
