package service

import (
	"context"

	"github.com/ethereal3x/apc/errs"
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

// HelloCheck 健康检查
func (server *HelloServer) HelloCheck(ctx context.Context, req *hellopb.HelloCheckRequest) (*hellopb.HelloCheckResponse, error) {
	return errs.Handle(&hellopb.HelloCheckResponse{}, func(rsp *hellopb.HelloCheckResponse) error {
		result, err := server.logic.Check(ctx, req.Check)
		if err != nil {
			return err
		}
		rsp.CheckDown = result.Message
		return nil
	})
}
