package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"github.com/ethereal3x/mint-server/internal/auth"
	"github.com/ethereal3x/mint-server/internal/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthMiddleware JWT 认证拦截器
type AuthMiddleware struct {
	manager       *util.TokenManager
	publicMethods map[string]struct{}
}

// AuthMiddlewareConfig JWT 认证拦截器配置
type AuthMiddlewareConfig struct {
	Manager       *util.TokenManager
	PublicMethods []string
}

// NewAuthMiddleware 创建 JWT 认证拦截器
func NewAuthMiddleware(config AuthMiddlewareConfig) *AuthMiddleware {
	publicMethods := make(map[string]struct{}, len(config.PublicMethods))
	for _, method := range config.PublicMethods {
		publicMethods[method] = struct{}{}
	}
	return &AuthMiddleware{manager: config.Manager, publicMethods: publicMethods}
}

// UnaryInterceptor 校验 unary gRPC 请求登录态
func (middleware *AuthMiddleware) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if middleware.isPublicMethod(info.FullMethod) {
		return handler(ctx, req)
	}
	authCtx, err := auth.ContextWithMetadataPrincipal(ctx, middleware.manager)
	if err != nil {
		return nil, authStatusError(err)
	}
	return handler(authCtx, req)
}

// StreamInterceptor 校验 stream gRPC 请求登录态
func (middleware *AuthMiddleware) StreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if middleware.isPublicMethod(info.FullMethod) {
		return handler(srv, stream)
	}
	authCtx, err := auth.ContextWithMetadataPrincipal(stream.Context(), middleware.manager)
	if err != nil {
		return authStatusError(err)
	}
	return handler(srv, &principalServerStream{ServerStream: stream, ctx: authCtx})
}

// isPublicMethod 判断 gRPC 方法是否无需登录态
func (middleware *AuthMiddleware) isPublicMethod(fullMethod string) bool {
	_, ok := middleware.publicMethods[fullMethod]
	return ok
}

// principalServerStream 覆盖 stream context 以传递认证主体
type principalServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context 返回带认证主体的 stream context
func (stream *principalServerStream) Context() context.Context {
	return stream.ctx
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

// UnaryObservabilityInterceptor 为 unary RPC 记录 tracing 与结构化日志
func UnaryObservabilityInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx, span := tracing.Start(ctx, info.FullMethod)
	defer span.End()
	startTime := time.Now()
	logger.ContextInfo(ctx, "grpc request", zap.String("method", info.FullMethod))
	response, err := handler(ctx, req)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "grpc request failed", zap.String("method", info.FullMethod), zap.Duration("elapsed", time.Since(startTime)), zap.Error(err))
		return response, err
	}
	logger.ContextInfo(ctx, "grpc request done", zap.String("method", info.FullMethod), zap.Duration("elapsed", time.Since(startTime)))
	return response, nil
}

// StreamObservabilityInterceptor 为 stream RPC 记录 tracing 与结构化日志
func StreamObservabilityInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx, span := tracing.Start(stream.Context(), info.FullMethod)
	defer span.End()
	startTime := time.Now()
	logger.ContextInfo(ctx, "grpc stream start", zap.String("method", info.FullMethod))
	err := handler(srv, &observabilityServerStream{ServerStream: stream, ctx: ctx})
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "grpc stream failed", zap.String("method", info.FullMethod), zap.Duration("elapsed", time.Since(startTime)), zap.Error(err))
		return err
	}
	logger.ContextInfo(ctx, "grpc stream done", zap.String("method", info.FullMethod), zap.Duration("elapsed", time.Since(startTime)))
	return nil
}

// observabilityServerStream 覆盖 stream context 以传递 tracing context
type observabilityServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context 返回带 tracing 的 stream context
func (stream *observabilityServerStream) Context() context.Context {
	return stream.ctx
}
