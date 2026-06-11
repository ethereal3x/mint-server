package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ethereal3x/mint-server/internal/model"
	"github.com/ethereal3x/mint-server/internal/util"
	"google.golang.org/grpc/metadata"
)

const authorizationHeader = "authorization"

// BearerTokenFromHeader 从 Authorization 头提取 Bearer token
func BearerTokenFromHeader(header string) (string, error) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], util.BearerTokenType) {
		return "", util.ErrInvalidToken
	}
	return parts[1], nil
}

// PrincipalFromRequest 从 HTTP 请求中解析认证主体
func PrincipalFromRequest(ctx context.Context, manager *util.TokenManager, request *http.Request) (*model.Principal, error) {
	tokenText, err := BearerTokenFromHeader(request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}
	claims, err := manager.ParseAccessToken(ctx, tokenText)
	if err != nil {
		return nil, fmt.Errorf("parse request token: %w", err)
	}
	return util.PrincipalFromClaims(claims), nil
}

// ContextWithMetadataPrincipal 从 gRPC metadata 解析认证主体并写入上下文
func ContextWithMetadataPrincipal(ctx context.Context, manager *util.TokenManager) (context.Context, error) {
	metadataMap, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, util.ErrMissingPrincipal
	}
	authValues := metadataMap.Get(authorizationHeader)
	if len(authValues) == 0 {
		authValues = metadataMap.Get("grpcgateway-" + authorizationHeader)
	}
	if len(authValues) == 0 {
		return ctx, util.ErrMissingPrincipal
	}
	tokenText, err := BearerTokenFromHeader(authValues[0])
	if err != nil {
		return ctx, err
	}
	claims, err := manager.ParseAccessToken(ctx, tokenText)
	if err != nil {
		return ctx, fmt.Errorf("parse metadata token: %w", err)
	}
	return WithPrincipal(ctx, util.PrincipalFromClaims(claims)), nil
}
