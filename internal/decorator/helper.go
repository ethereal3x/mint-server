package decorator

import (
	"context"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"go.uber.org/zap"
)

// wrap 为业务调用添加 tracing span 和结构化日志
func wrap[R any](ctx context.Context, name string, fn func(context.Context) (R, error), fields ...zap.Field) (R, error) {
	ctx, span := tracing.Start(ctx, name)
	defer span.End()
	logger.ContextInfo(ctx, name, fields...)
	result, err := fn(ctx)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, name, append(fields, zap.Error(err))...)
	}
	return result, err
}

// wrapErr 为仅返回 error 的业务调用添加 tracing span 和结构化日志
func wrapErr(ctx context.Context, name string, fn func(context.Context) error, fields ...zap.Field) error {
	_, err := wrap(ctx, name, func(ctx context.Context) (struct{}, error) {
		return struct{}{}, fn(ctx)
	}, fields...)
	return err
}
