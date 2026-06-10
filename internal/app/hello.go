package app

import (
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/service"
)

// HelloModule 聚合健康检查的 gRPC 服务处理器
type HelloModule struct {
	HelloServer *service.HelloServer
}

// newHelloModule 装配 Hello 领域全部依赖
func newHelloModule() *HelloModule {
	return &HelloModule{
		HelloServer: service.NewHelloServer(logic.NewHello()),
	}
}
