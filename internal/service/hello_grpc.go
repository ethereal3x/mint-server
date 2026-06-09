package service

import (
	"context"

	hellopb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/hello"
	"github.com/ethereal3x/mint-server/internal/logic"
)

// HelloServer HelloService gRPC 处理器
type HelloServer struct {
	hellopb.UnimplementedHelloServiceServer
	logic *logic.Hello
}

// NewHelloServer 创建 Hello gRPC 处理器
func NewHelloServer(helloLogic *logic.Hello) *HelloServer {
	return &HelloServer{logic: helloLogic}
}

func (s *HelloServer) HelloCheck(ctx context.Context, req *hellopb.HelloCheckRequest) (*hellopb.HelloCheckResponse, error) {
	return s.logic.Check(ctx, req)
}
