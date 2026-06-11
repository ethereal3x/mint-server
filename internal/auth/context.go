package auth

import (
	"context"

	"github.com/ethereal3x/mint-server/internal/model"
	"github.com/ethereal3x/mint-server/internal/util"
)

type principalContextKey struct{}

// WithPrincipal 将认证主体写入上下文
func WithPrincipal(ctx context.Context, principal *model.Principal) context.Context {
	return context.WithValue(ctx, principalContextKey{}, principal)
}

// PrincipalFromContext 从上下文读取认证主体
func PrincipalFromContext(ctx context.Context) (*model.Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(*model.Principal)
	return principal, ok
}

// RequireUserID 从上下文读取当前用户ID
func RequireUserID(ctx context.Context) (string, error) {
	principal, ok := PrincipalFromContext(ctx)
	if !ok || principal == nil || principal.UserID == "" {
		return "", util.ErrMissingPrincipal
	}
	return principal.UserID, nil
}
