package logic

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"github.com/ethereal3x/mint-server/internal/auth"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/idgen"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
)

// UserStore 用户认证场景需要的数据访问接口
type UserStore interface {
	CreateBaseUser(ctx context.Context, user *model.BaseUser) error
	FindBaseUserByUsername(ctx context.Context, query *model.BaseUserQuery) (*model.BaseUser, error)
	FindBaseUserByUserID(ctx context.Context, query *model.BaseUserQuery) (*model.BaseUser, error)
	UpdateBaseUser(ctx context.Context, userID int64, updates map[string]any) error
}

// Auth 认证业务逻辑
type Auth struct {
	repo         UserStore
	tokenManager *auth.TokenManager
}

// RegisterAccountRequest 账号密码注册参数
type RegisterAccountRequest struct {
	Account     string
	Password    string
	DisplayName string
	AvatarURL   string
}

// LoginRequest 登录参数
type LoginRequest struct {
	Provider   string
	Identifier string
	Credential string
}

// AuthResult 登录态结果
type AuthResult struct {
	User  *model.BaseUser
	Token *auth.TokenResult
}

// issueTokenRequest 签发登录态参数
type issueTokenRequest struct {
	User       *model.BaseUser
	Provider   string
	Identifier string
}

// NewAuth 创建认证业务逻辑
func NewAuth(repo UserStore, tokenManager *auth.TokenManager) *Auth {
	return &Auth{repo: repo, tokenManager: tokenManager}
}

// RegisterAccount 执行账号密码注册并签发登录令牌
func (s *Auth) RegisterAccount(ctx context.Context, req *RegisterAccountRequest) (*AuthResult, error) {
	ctx, span := tracing.Start(ctx, "logic.Auth.RegisterAccount")
	defer span.End()

	account := normalizeIdentifier(req.Account)
	if account == "" {
		return nil, mint_err.ErrParam
	}
	credentialHash, err := auth.HashPassword(req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrPasswordTooShort) {
			return nil, mint_err.ErrPasswordWeak
		}
		logger.ContextError(ctx, "Auth.RegisterAccount", zap.Error(err))
		return nil, mint_err.ErrInternal
	}

	nickname := strings.TrimSpace(req.DisplayName)
	if nickname == "" {
		nickname = account
	}
	now := time.Now().Unix()
	user := &model.BaseUser{
		UserID:   idgen.GenIntUUID(),
		Username:  account,
		Nickname:  nickname,
		AvatarURL: strings.TrimSpace(req.AvatarURL),
		Password:  credentialHash,
		RegTime:  now,
	}
	if err := s.repo.CreateBaseUser(ctx, user); err != nil {
		tracing.RecordError(ctx, err)
		if errors.Is(err, mint_err.ErrUserExists) {
			return nil, mint_err.ErrUserExists
		}
		logger.ContextError(ctx, "Auth.RegisterAccount", zap.String("account", account), zap.Error(err))
		return nil, mint_err.ErrDBCreate
	}
	return s.issueToken(ctx, &issueTokenRequest{User: user, Provider: model.AUTH_PROVIDER_ACCOUNT_PASSWORD, Identifier: account})
}

// Login 校验登录凭证并签发登录令牌
func (s *Auth) Login(ctx context.Context, req *LoginRequest) (*AuthResult, error) {
	ctx, span := tracing.Start(ctx, "logic.Auth.Login")
	defer span.End()

	provider := req.Provider
	if provider == "" {
		provider = model.AUTH_PROVIDER_ACCOUNT_PASSWORD
	}
	identifier := normalizeIdentifier(req.Identifier)
	if provider != model.AUTH_PROVIDER_ACCOUNT_PASSWORD || identifier == "" || req.Credential == "" {
		return nil, mint_err.ErrParam
	}
	user, err := s.repo.FindBaseUserByUsername(ctx, &model.BaseUserQuery{Username: identifier})
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Auth.Login", zap.String("provider", provider), zap.String("identifier", identifier), zap.Error(err))
		return nil, mint_err.ErrDBQuery
	}
	if user == nil || !auth.ComparePassword(user.Password, req.Credential) {
		return nil, mint_err.ErrInvalidCredential
	}
	return s.issueToken(ctx, &issueTokenRequest{User: user, Provider: provider, Identifier: identifier})
}

// GetMe 获取当前登录用户信息
func (s *Auth) GetMe(ctx context.Context, userID string) (*model.BaseUser, error) {
	numericUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, mint_err.ErrTokenInvalid
	}
	user, err := s.repo.FindBaseUserByUserID(ctx, &model.BaseUserQuery{UserID: numericUserID})
	if err != nil {
		logger.ContextError(ctx, "Auth.GetMe", zap.String("user_id", userID), zap.Error(err))
		return nil, mint_err.ErrDBQuery
	}
	if user == nil {
		return nil, mint_err.ErrUserNotFound
	}
	return user, nil
}

// UpdateAvatar 更新当前用户的头像
func (s *Auth) UpdateAvatar(ctx context.Context, userID string, avatarURL string) (*model.BaseUser, error) {
	numericUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, mint_err.ErrTokenInvalid
	}
	avatarURL = strings.TrimSpace(avatarURL)
	if avatarURL == "" {
		return nil, mint_err.ErrParam
	}
	if err := s.repo.UpdateBaseUser(ctx, numericUserID, map[string]any{"avatar_url": avatarURL}); err != nil {
		logger.ContextError(ctx, "Auth.UpdateAvatar", zap.String("user_id", userID), zap.Error(err))
		return nil, mint_err.ErrDBUpdate
	}
	user, err := s.repo.FindBaseUserByUserID(ctx, &model.BaseUserQuery{UserID: numericUserID})
	if err != nil {
		return nil, mint_err.ErrDBQuery
	}
	if user == nil {
		return nil, mint_err.ErrUserNotFound
	}
	return user, nil
}

// UpdatePassword 更新当前用户的登录密码
func (s *Auth) UpdatePassword(ctx context.Context, userID string, oldPassword string, newPassword string) error {
	numericUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return mint_err.ErrTokenInvalid
	}
	if oldPassword == "" || newPassword == "" {
		return mint_err.ErrParam
	}
	user, err := s.repo.FindBaseUserByUserID(ctx, &model.BaseUserQuery{UserID: numericUserID})
	if err != nil {
		logger.ContextError(ctx, "Auth.UpdatePassword", zap.String("user_id", userID), zap.Error(err))
		return mint_err.ErrDBQuery
	}
	if user == nil || !auth.ComparePassword(user.Password, oldPassword) {
		return mint_err.ErrInvalidCredential
	}
	credentialHash, err := auth.HashPassword(newPassword)
	if err != nil {
		if errors.Is(err, auth.ErrPasswordTooShort) {
			return mint_err.ErrPasswordWeak
		}
		logger.ContextError(ctx, "Auth.UpdatePassword", zap.Error(err))
		return mint_err.ErrInternal
	}
	if err := s.repo.UpdateBaseUser(ctx, numericUserID, map[string]any{"password": credentialHash}); err != nil {
		logger.ContextError(ctx, "Auth.UpdatePassword", zap.String("user_id", userID), zap.Error(err))
		return mint_err.ErrDBUpdate
	}
	return nil
}

// UpdateNickname 更新当前用户的展示名称
func (s *Auth) UpdateNickname(ctx context.Context, userID string, nickname string) (*model.BaseUser, error) {
	numericUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, mint_err.ErrTokenInvalid
	}
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return nil, mint_err.ErrParam
	}
	if err := s.repo.UpdateBaseUser(ctx, numericUserID, map[string]any{"nickname": nickname}); err != nil {
		logger.ContextError(ctx, "Auth.UpdateNickname", zap.String("user_id", userID), zap.Error(err))
		return nil, mint_err.ErrDBUpdate
	}
	user, err := s.repo.FindBaseUserByUserID(ctx, &model.BaseUserQuery{UserID: numericUserID})
	if err != nil {
		return nil, mint_err.ErrDBQuery
	}
	if user == nil {
		return nil, mint_err.ErrUserNotFound
	}
	return user, nil
}

// issueToken 为用户签发访问令牌
func (s *Auth) issueToken(ctx context.Context, req *issueTokenRequest) (*AuthResult, error) {
	token, err := s.tokenManager.IssueAccessToken(ctx, &auth.TokenInput{UserID: strconv.FormatInt(req.User.UserID, 10), Provider: req.Provider, Identifier: req.Identifier})
	if err != nil {
		logger.ContextError(ctx, "Auth.issueToken", zap.Int64("user_id", req.User.UserID), zap.Error(err))
		return nil, mint_err.ErrInternal
	}
	return &AuthResult{User: req.User, Token: token}, nil
}

// normalizeIdentifier 规范化登录标识
func normalizeIdentifier(identifier string) string {
	return strings.ToLower(strings.TrimSpace(identifier))
}
