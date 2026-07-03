package logic

import (
	"context"
)

// HelloResult 健康检查结果
type HelloResult struct {
	Message string
}

// Hello 基础健康检查业务逻辑
type Hello struct{}

// NewHello 创建健康检查业务逻辑
func NewHello() *Hello {
	return &Hello{}
}

// Check 执行健康检查
func (hello *Hello) Check(ctx context.Context, input string) (*HelloResult, error) {
	return &HelloResult{Message: "ok: " + input}, nil
}
