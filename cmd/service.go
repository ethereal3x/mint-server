package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	apccfg "github.com/ethereal3x/apc/config"
	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/orm"
	"github.com/ethereal3x/apc/server"
	"github.com/ethereal3x/apc/storage"
	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	authpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/auth"
	hellopb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/hello"
	userpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/user"
	"github.com/ethereal3x/mint-server/internal/auth"
	"github.com/ethereal3x/mint-server/internal/config"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/repo"
	"github.com/ethereal3x/mint-server/internal/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

var (
	globalDB           *gorm.DB
	globalServices     *serviceSet
	globalTokenManager *auth.TokenManager
)

type serviceSet struct {
	agentService  *service.AgentServer
	authService   *service.AuthServer
	configService *service.ConfigServer
	helloService  *service.HelloServer
	userService   *service.UserServer
	uploadHandler *service.UploadHandler
}

// initApp 初始化应用依赖
func initApp() {
	if err := config.InitBusinessConfig("config.yaml"); err != nil {
		log.Fatalf("load business config: %v", err)
	}

	apcConf := apccfg.GetConf()
	logger.SetLogger(logger.NewLogger(&apcConf.Plugin.Log))

	globalDB = initGorm()
	if globalDB == nil {
		log.Fatal("init gorm: db is nil")
	}

	secretKey := []byte(config.GetBusinessConfig().SecretKey.Encryption)
	configRepo := repo.NewModelConfigRepo(globalDB, secretKey)
	recordRepo := repo.NewDialogueRepo(globalDB)
	uploadRepo := repo.NewFileUploadRepo(globalDB)
	userRepo := repo.NewUserRepo(globalDB)

	storageClient := initStorage()
	globalTokenManager = auth.NewTokenManager(auth.TokenManagerConfig{Secret: config.GetBusinessConfig().JWT.Secret})

	configLogic := logic.NewConfig(configRepo)
	chatAdapter := logic.NewEinoAdapter()
	chatLogic := logic.NewChat(configRepo, recordRepo, chatAdapter)
	helloLogic := logic.NewHello()
	uploadLogic := logic.NewUploadLogic(storageClient, uploadRepo)
	authLogic := logic.NewAuth(userRepo, globalTokenManager)

	globalServices = &serviceSet{
		agentService:  service.NewAgentServer(chatLogic, configLogic),
		authService:   service.NewAuthServer(authLogic),
		configService: service.NewConfigServer(configLogic, chatLogic),
		helloService:  service.NewHelloServer(helloLogic),
		userService:   service.NewUserServer(authLogic),
		uploadHandler: service.NewUploadHandler(uploadLogic, globalTokenManager),
	}
}

// initStorage 初始化 rustfs 对象存储客户端
func initStorage() storage.ObjectStorage {
	cfg := apccfg.GetConf()
	storageClient, err := storage.NewS3Client(context.Background(), storage.NewS3ClientParams{
		Provider: storage.STORAGE_PROVIDER_RUSTFS,
		Config:   &cfg.Plugin.RustFS,
	})
	if err != nil {
		log.Fatalf("init storage client: %v", err)
	}
	return storageClient
}

// initGorm 根据配置初始化 GORM 数据库连接
func initGorm() *gorm.DB {
	clientConf := apccfg.GetClientConf("mysql")
	if clientConf == nil {
		return nil
	}
	gormCfg := apccfg.GenGormConfig(clientConf)
	db, err := orm.NewGormInstance(gormCfg)
	if err != nil {
		return nil
	}
	return db
}

// newGrpcServer 创建 gRPC 服务器并注册所有服务
func newGrpcServer(globalTokenManager *auth.TokenManager) *server.GrpcServer {
	rs := server.NewRpcServer()
	authMiddleware := auth.NewMiddleware(&auth.MiddlewareConfig{
		Manager:       globalTokenManager,
		PublicMethods: publicMethods(),
	})
	rs.SetInterceptors(
		[]grpc.StreamServerInterceptor{authMiddleware.StreamInterceptor},
		[]grpc.UnaryServerInterceptor{authMiddleware.UnaryInterceptor},
	)
	rs.SetRegisterFunc(func(s *grpc.Server) {
		if globalServices == nil {
			return
		}
		hellopb.RegisterHelloServiceServer(s, globalServices.helloService)
		authpb.RegisterAuthServiceServer(s, globalServices.authService)
		userpb.RegisterUserServiceServer(s, globalServices.userService)
		agentpb.RegisterAgentServiceServer(s, globalServices.agentService)
		agentpb.RegisterModelConfigServiceServer(s, globalServices.configService)
	})
	return rs
}

// newGatewayServer 创建 HTTP gateway 服务器并注册所有路由
func newGatewayServer() *server.HttpServer {
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
			globalServices.uploadHandler.HandleUpload(w, r)
		})); err != nil {
			return err
		}
		if err := mux.HandlePath("GET", "/api/v1/files/{id}", runtime.HandlerFunc(func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			r.SetPathValue("id", pathParams["id"])
			globalServices.uploadHandler.HandleGetFile(w, r)
		})); err != nil {
			return err
		}
		if err := mux.HandlePath("GET", "/api/v1/files", runtime.HandlerFunc(func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
			globalServices.uploadHandler.HandleListFiles(w, r)
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
