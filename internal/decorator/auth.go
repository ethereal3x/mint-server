package decorator

import (
	"context"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
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
	ctx, span := tracing.Start(ctx, "Auth.RegisterAccount")
	defer span.End()
	logger.ContextInfo(ctx, "Auth.RegisterAccount", zap.String("account", req.Account))
	result, err := d.inner.RegisterAccount(ctx, req)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Auth.RegisterAccount", zap.Error(err))
	}
	return result, err
}

// Login 登录并装饰 tracing/logging
func (d *AuthDecorator) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResult, error) {
	ctx, span := tracing.Start(ctx, "Auth.Login")
	defer span.End()
	logger.ContextInfo(ctx, "Auth.Login", zap.String("provider", req.Provider), zap.String("identifier", req.Identifier))
	result, err := d.inner.Login(ctx, req)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Auth.Login", zap.Error(err))
	}
	return result, err
}

// GetMe 获取当前用户信息并装饰 tracing/logging
func (d *AuthDecorator) GetMe(ctx context.Context, userID string) (*model.BaseUser, error) {
	ctx, span := tracing.Start(ctx, "Auth.GetMe")
	defer span.End()
	user, err := d.inner.GetMe(ctx, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Auth.GetMe", zap.Error(err))
	}
	return user, err
}

// UpdateAvatar 更新头像并装饰 tracing/logging
func (d *AuthDecorator) UpdateAvatar(ctx context.Context, userID string, avatarURL string) (*model.BaseUser, error) {
	ctx, span := tracing.Start(ctx, "Auth.UpdateAvatar")
	defer span.End()
	user, err := d.inner.UpdateAvatar(ctx, userID, avatarURL)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Auth.UpdateAvatar", zap.Error(err))
	}
	return user, err
}

// UpdatePassword 更新密码并装饰 tracing/logging
func (d *AuthDecorator) UpdatePassword(ctx context.Context, userID string, oldPassword string, newPassword string) error {
	ctx, span := tracing.Start(ctx, "Auth.UpdatePassword")
	defer span.End()
	err := d.inner.UpdatePassword(ctx, userID, oldPassword, newPassword)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Auth.UpdatePassword", zap.Error(err))
	}
	return err
}

// UpdateNickname 更新昵称并装饰 tracing/logging
func (d *AuthDecorator) UpdateNickname(ctx context.Context, userID string, nickname string) (*model.BaseUser, error) {
	ctx, span := tracing.Start(ctx, "Auth.UpdateNickname")
	defer span.End()
	user, err := d.inner.UpdateNickname(ctx, userID, nickname)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Auth.UpdateNickname", zap.Error(err))
	}
	return user, err
}
