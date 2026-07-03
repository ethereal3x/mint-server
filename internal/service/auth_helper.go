package service

import (
	"context"

	"github.com/ethereal3x/mint-server/internal/auth"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
)

// requireUserID 从上下文读取当前用户 ID
func requireUserID(ctx context.Context) (string, error) {
	userID, err := auth.RequireUserID(ctx)
	if err != nil {
		return "", mint_err.ErrUnauthorized
	}
	return userID, nil
}
