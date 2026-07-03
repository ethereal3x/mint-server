package app

import (
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/repo"
	"github.com/ethereal3x/mint-server/internal/service"
	"github.com/ethereal3x/mint-server/internal/util"
	"gorm.io/gorm"
)

// AuthModule 聚合认证和用户的 gRPC 服务处理器
type AuthModule struct {
	AuthServer *service.AuthServer
	UserServer *service.UserServer
}

// newAuthModule 装配 Auth 领域全部依赖
func newAuthModule(db *gorm.DB, tokenManager *util.TokenManager) *AuthModule {
	userRepo := repo.NewUserRepo(db)
	authLogic := logic.NewAuth(userRepo, tokenManager)
	return &AuthModule{
		AuthServer: service.NewAuthServer(authLogic),
		UserServer: service.NewUserServer(authLogic),
	}
}
