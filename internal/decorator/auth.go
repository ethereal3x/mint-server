package decorator

import (
	"context"

	"github.com/ethereal3x/mint-server/internal/dto"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
)

// authLogic AuthDecorator 需要的认证业务接口
type authLogic interface {
	RegisterAccount(ctx context.Context, req *dto.RegisterAccountRequest) (*dto.AuthResult, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResult, error)
	GetMe(ctx context.Context, userID string) (*model.BaseUser, error)
	UpdateAvatar(ctx context.Context, userID string, avatarURL string) (*model.BaseUser, error)
	UpdatePassword(ctx context.Context, userID string, oldPassword string, newPassword string) error
	UpdateNickname(ctx context.Context, userID string, nickname string) (*model.BaseUser, error)
}

// AuthDecorator 为认证业务逻辑添加 tracing 和 logging 横切关注点
type AuthDecorator struct {
	inner authLogic
}

// NewAuthDecorator 创建认证逻辑装饰器
func NewAuthDecorator(inner authLogic) *AuthDecorator {
	return &AuthDecorator{inner: inner}
}

// RegisterAccount 注册账号并装饰 tracing/logging
func (d *AuthDecorator) RegisterAccount(ctx context.Context, req *dto.RegisterAccountRequest) (*dto.AuthResult, error) {
	return wrap(ctx, "Auth.RegisterAccount", func(ctx context.Context) (*dto.AuthResult, error) {
		return d.inner.RegisterAccount(ctx, req)
	}, zap.String("account", req.Account))
}

// Login 登录并装饰 tracing/logging
func (d *AuthDecorator) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResult, error) {
	return wrap(ctx, "Auth.Login", func(ctx context.Context) (*dto.AuthResult, error) {
		return d.inner.Login(ctx, req)
	}, zap.String("provider", req.Provider), zap.String("identifier", req.Identifier))
}

// GetMe 获取当前用户信息并装饰 tracing/logging
func (d *AuthDecorator) GetMe(ctx context.Context, userID string) (*model.BaseUser, error) {
	return wrap(ctx, "Auth.GetMe", func(ctx context.Context) (*model.BaseUser, error) {
		return d.inner.GetMe(ctx, userID)
	})
}

// UpdateAvatar 更新头像并装饰 tracing/logging
func (d *AuthDecorator) UpdateAvatar(ctx context.Context, userID string, avatarURL string) (*model.BaseUser, error) {
	return wrap(ctx, "Auth.UpdateAvatar", func(ctx context.Context) (*model.BaseUser, error) {
		return d.inner.UpdateAvatar(ctx, userID, avatarURL)
	})
}

// UpdatePassword 更新密码并装饰 tracing/logging
func (d *AuthDecorator) UpdatePassword(ctx context.Context, userID string, oldPassword string, newPassword string) error {
	return wrapErr(ctx, "Auth.UpdatePassword", func(ctx context.Context) error {
		return d.inner.UpdatePassword(ctx, userID, oldPassword, newPassword)
	})
}

// UpdateNickname 更新昵称并装饰 tracing/logging
func (d *AuthDecorator) UpdateNickname(ctx context.Context, userID string, nickname string) (*model.BaseUser, error) {
	return wrap(ctx, "Auth.UpdateNickname", func(ctx context.Context) (*model.BaseUser, error) {
		return d.inner.UpdateNickname(ctx, userID, nickname)
	})
}
