package auth

import "context"

type principalContextKey struct{}

// Principal 表示当前认证主体
type Principal struct {
	UserID     string
	Provider   string
	Identifier string
}

// WithPrincipal 将认证主体写入上下文
func WithPrincipal(ctx context.Context, principal *Principal) context.Context {
	return context.WithValue(ctx, principalContextKey{}, principal)
}

// PrincipalFromContext 从上下文读取认证主体
func PrincipalFromContext(ctx context.Context) (*Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(*Principal)
	return principal, ok
}

// RequireUserID 从上下文读取当前用户ID
func RequireUserID(ctx context.Context) (string, error) {
	principal, ok := PrincipalFromContext(ctx)
	if !ok || principal == nil || principal.UserID == "" {
		return "", ErrMissingPrincipal
	}
	return principal.UserID, nil
}
