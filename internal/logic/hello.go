package logic

import (
	"context"

	hellopb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/hello"
)

// Hello 基础健康检查业务逻辑
type Hello struct{}

// NewHello 创建健康检查业务逻辑
func NewHello() *Hello {
	return &Hello{}
}

// Check 执行健康检查
func (s *Hello) Check(ctx context.Context, req *hellopb.HelloCheckRequest) (*hellopb.HelloCheckResponse, error) {
	return &hellopb.HelloCheckResponse{
		CheckDown: "ok: " + req.Check,
	}, nil
}
