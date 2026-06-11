package main

import (
	"context"
	"errors"

	"github.com/ethereal3x/mint-server/internal/auth"
	"github.com/ethereal3x/mint-server/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Middleware JWT 认证拦截器
type Middleware struct {
	manager       *util.TokenManager
	publicMethods map[string]struct{}
}

// MiddlewareConfig JWT 认证拦截器配置
type MiddlewareConfig struct {
	Manager       *util.TokenManager
	PublicMethods []string
}

// NewMiddleware 创建 JWT 认证拦截器
func NewMiddleware(config *MiddlewareConfig) *Middleware {
	publicMethods := make(map[string]struct{}, len(config.PublicMethods))
	for _, method := range config.PublicMethods {
		publicMethods[method] = struct{}{}
	}
	return &Middleware{manager: config.Manager, publicMethods: publicMethods}
}

// UnaryInterceptor 校验 unary gRPC 请求登录态
func (m *Middleware) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if m.isPublicMethod(info.FullMethod) {
		return handler(ctx, req)
	}
	authCtx, err := auth.ContextWithMetadataPrincipal(ctx, m.manager)
	if err != nil {
		return nil, authStatusError(err)
	}
	return handler(authCtx, req)
}

// StreamInterceptor 校验 stream gRPC 请求登录态
func (m *Middleware) StreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if m.isPublicMethod(info.FullMethod) {
		return handler(srv, stream)
	}
	authCtx, err := auth.ContextWithMetadataPrincipal(stream.Context(), m.manager)
	if err != nil {
		return authStatusError(err)
	}
	return handler(srv, &principalServerStream{ServerStream: stream, ctx: authCtx})
}

// isPublicMethod 判断 gRPC 方法是否无需登录态
func (m *Middleware) isPublicMethod(fullMethod string) bool {
	_, ok := m.publicMethods[fullMethod]
	return ok
}

// principalServerStream 覆盖 stream context 以传递认证主体
type principalServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context 返回带认证主体的 stream context
func (s *principalServerStream) Context() context.Context {
	return s.ctx
}

// authStatusError 转换认证错误为 gRPC status error
func authStatusError(err error) error {
	if errors.Is(err, util.ErrExpiredToken) {
		return status.Error(codes.Unauthenticated, "登录已过期")
	}
	if errors.Is(err, util.ErrInvalidToken) {
		return status.Error(codes.Unauthenticated, "登录令牌无效")
	}
	return status.Error(codes.Unauthenticated, "未登录或登录已失效")
}
