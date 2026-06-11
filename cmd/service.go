package main

import (
	"context"
	"net/http"
	"strings"

	apccfg "github.com/ethereal3x/apc/config"
	"github.com/ethereal3x/apc/server"
	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	authpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/auth"
	hellopb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/hello"
	userpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/user"
	"github.com/ethereal3x/mint-server/internal/app"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// newGrpcServer 创建 gRPC 服务器并注册所有服务
func newGrpcServer(app *app.App) *server.GrpcServer {
	rs := server.NewRpcServer()
	authMiddleware := NewMiddleware(&MiddlewareConfig{
		Manager:       app.TokenManager,
		PublicMethods: publicMethods(),
	})
	rs.SetInterceptors(
		[]grpc.StreamServerInterceptor{authMiddleware.StreamInterceptor},
		[]grpc.UnaryServerInterceptor{authMiddleware.UnaryInterceptor},
	)
	rs.SetRegisterFunc(func(s *grpc.Server) {
		hellopb.RegisterHelloServiceServer(s, app.Hello.HelloServer)
		authpb.RegisterAuthServiceServer(s, app.Authx.AuthServer)
		userpb.RegisterUserServiceServer(s, app.Authx.UserServer)
		agentpb.RegisterAgentServiceServer(s, app.Agent.AgentServer)
		agentpb.RegisterModelConfigServiceServer(s, app.Agent.ConfigServer)
	})
	return rs
}

// newGatewayServer 创建 HTTP gateway 服务器并注册所有路由
func newGatewayServer(app *app.App) *server.HttpServer {
	hs := server.NewHttpServer()
	hs.SetWriteTimeout(0) // 流式聊天响应不设超时
	hs.SetServeMuxOpts([]runtime.ServeMuxOption{runtime.WithIncomingHeaderMatcher(authHeaderMatcher)})
	hs.SetRegisterFunc(func(ctx context.Context, mux *runtime.ServeMux) error {
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		addr := ":" + grpcPort()

		if err := hellopb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
			return err
		}
		if err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
			return err
		}
		if err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
			return err
		}
		if err := agentpb.RegisterAgentServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
			return err
		}
		if err := agentpb.RegisterModelConfigServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
			return err
		}

		// 注册文件上传 HTTP 路由（multipart 不适合 gRPC-gateway）
		if err := mux.HandlePath("POST", "/api/v1/files/upload", runtime.HandlerFunc(func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
			app.Upload.UploadHandler.HandleUpload(w, r)
		})); err != nil {
			return err
		}
		if err := mux.HandlePath("GET", "/api/v1/files/{id}", runtime.HandlerFunc(func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			r.SetPathValue("id", pathParams["id"])
			app.Upload.UploadHandler.HandleGetFile(w, r)
		})); err != nil {
			return err
		}
		if err := mux.HandlePath("GET", "/api/v1/files", runtime.HandlerFunc(func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
			app.Upload.UploadHandler.HandleListFiles(w, r)
		})); err != nil {
			return err
		}

		return nil
	})
	return hs
}

// authHeaderMatcher 透传认证头到 gRPC metadata
func authHeaderMatcher(key string) (string, bool) {
	if strings.EqualFold(key, "Authorization") {
		return "authorization", true
	}
	return runtime.DefaultHeaderMatcher(key)
}

// publicMethods 无需登录态的 gRPC 方法列表
func publicMethods() []string {
	return []string{
		"/mint_server.auth.AuthService/RegisterAccount",
		"/mint_server.auth.AuthService/Login",
		"/mint_server.hello.HelloService/HelloCheck",
	}
}

// grpcPort 从配置中提取 gRPC 端口号
func grpcPort() string {
	return strings.TrimPrefix(apccfg.GetConf().Server.GrpcAddr, ":")
}
