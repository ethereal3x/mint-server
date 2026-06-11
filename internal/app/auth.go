package app

import (
	"github.com/ethereal3x/mint-server/internal/decorator"
	"github.com/ethereal3x/mint-server/internal/util"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/repo"
	"github.com/ethereal3x/mint-server/internal/service"
	"gorm.io/gorm"
)

// AuthModule 聚合认证和用户的 gRPC 服务处理器
type AuthModule struct {
	AuthServer *service.AuthServer
	UserServer *service.UserServer
}

// newAuthModule 装配 Auth 领域全部依赖，通过 decorator 注入 tracing/logging
func newAuthModule(db *gorm.DB, tokenManager *util.TokenManager) *AuthModule {
	userRepo := repo.NewUserRepo(db)
	authLogic := logic.NewAuth(userRepo, tokenManager)
	authDecorated := decorator.NewAuthDecorator(authLogic)

	return &AuthModule{
		AuthServer: service.NewAuthServer(authDecorated),
		UserServer: service.NewUserServer(authDecorated),
	}
}
