package middleware

import (
	"net/http"

	"github.com/ethereal3x/mint-server/internal/auth"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/util"
)

// RequireUserIDFromRequest 从 HTTP 请求解析当前用户 ID
func RequireUserIDFromRequest(r *http.Request, tokenManager *util.TokenManager) (string, error) {
	principal, err := auth.PrincipalFromRequest(r.Context(), tokenManager, r)
	if err != nil || principal.UserID == "" {
		return "", mint_err.ErrUnauthorized
	}
	return principal.UserID, nil
}
