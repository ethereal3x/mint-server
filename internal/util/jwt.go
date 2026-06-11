package util

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereal3x/mint-server/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultIssuer          = "mint-server"
	defaultAccessTokenTTL  = 24 * time.Hour
	BearerTokenType        = "Bearer"
	accessTokenSigningName = "access_token"
)

// ErrMissingPrincipal 表示上下文缺少认证主体
var ErrMissingPrincipal = errors.New("missing principal")

// ErrInvalidToken 表示 JWT 令牌无效
var ErrInvalidToken = errors.New("invalid token")

// ErrExpiredToken 表示 JWT 令牌已过期
var ErrExpiredToken = errors.New("expired token")

// TokenManagerConfig JWT 管理器配置
type TokenManagerConfig struct {
	Secret string
	Issuer string
	TTL    time.Duration
}

// TokenInput JWT 签发参数
type TokenInput struct {
	UserID     string
	Provider   string
	Identifier string
}

// TokenResult JWT 签发结果
type TokenResult struct {
	AccessToken string
	TokenType   string
	ExpiresAt   time.Time
	ExpiresIn   int64
}

// Claims JWT 业务声明
type Claims struct {
	Provider   string `json:"provider"`
	Identifier string `json:"identifier"`
	jwt.RegisteredClaims
}

// TokenManager 负责 JWT 签发和解析
type TokenManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

// NewTokenManager 创建 JWT 管理器
func NewTokenManager(config TokenManagerConfig) *TokenManager {
	issuer := config.Issuer
	if issuer == "" {
		issuer = defaultIssuer
	}
	ttl := config.TTL
	if ttl <= 0 {
		ttl = defaultAccessTokenTTL
	}
	return &TokenManager{secret: []byte(config.Secret), issuer: issuer, ttl: ttl}
}

// IssueAccessToken 签发访问令牌
func (m *TokenManager) IssueAccessToken(ctx context.Context, input *TokenInput) (*TokenResult, error) {
	now := time.Now()
	expiresAt := now.Add(m.ttl)
	claims := &Claims{
		Provider:   input.Provider,
		Identifier: input.Identifier,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   input.UserID,
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        accessTokenSigningName,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(m.secret)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}
	return &TokenResult{AccessToken: accessToken, TokenType: BearerTokenType, ExpiresAt: expiresAt, ExpiresIn: int64(m.ttl.Seconds())}, nil
}

// ParseAccessToken 解析并校验访问令牌
func (m *TokenManager) ParseAccessToken(ctx context.Context, tokenText string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenText, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	}, jwt.WithIssuer(m.issuer))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("parse access token: %w", err)
	}
	if !token.Valid || claims.Subject == "" {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// PrincipalFromClaims 根据 JWT 声明构建认证主体
func PrincipalFromClaims(claims *Claims) *model.Principal {
	return &model.Principal{UserID: claims.Subject, Provider: claims.Provider, Identifier: claims.Identifier}
}
