package main

import (
	"context"
	"log"
	"strings"

	apccfg "github.com/ethereal3x/apc/config"
	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/orm"
	"github.com/ethereal3x/apc/server"
	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	hellopb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/hello"
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
	globalDB       *gorm.DB
	globalServices *serviceSet
)

type serviceSet struct {
	agentService    *service.AgentServer
	strategyService *service.StrategyServer
	mappingService  *service.MappingServer
	helloService    *service.HelloServer
}

// initApp 初始化应用依赖
func initApp() {
	if err := config.InitBusinessConfig("config.yaml"); err != nil {
		log.Fatalf("load business config: %v", err)
	}

	logger.SetLogger(logger.NewLogger())

	globalDB = initGorm()
	if globalDB == nil {
		log.Fatal("init gorm: db is nil")
	}

	secretKey := []byte(config.GetBusinessConfig().SecretKey.Encryption)
	strategyRepo := repo.NewStrategyRepo(globalDB)
	mappingRepo := repo.NewMappingRepo(globalDB)
	recordRepo := repo.NewDialogueRepo(globalDB)

	strategyLogic := logic.NewStrategy(strategyRepo, secretKey)
	mappingLogic := logic.NewMapping(mappingRepo)
	chatLogic := logic.NewChat(mappingRepo, recordRepo, strategyLogic)
	helloLogic := logic.NewHello()

	globalServices = &serviceSet{
		agentService:    service.NewAgentServer(chatLogic, mappingLogic),
		strategyService: service.NewStrategyServer(strategyLogic),
		mappingService:  service.NewMappingServer(mappingLogic),
		helloService:    service.NewHelloServer(helloLogic),
	}
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
func newGrpcServer() *server.GrpcServer {
	rs := server.NewRpcServer()
	rs.SetRegisterFunc(func(s *grpc.Server) {
		if globalServices == nil {
			return
		}
		hellopb.RegisterHelloServiceServer(s, globalServices.helloService)
		agentpb.RegisterAgentServiceServer(s, globalServices.agentService)
		agentpb.RegisterStrategyServiceServer(s, globalServices.strategyService)
		agentpb.RegisterMappingServiceServer(s, globalServices.mappingService)
	})
	return rs
}

// newGatewayServer 创建 HTTP gateway 服务器并注册所有路由
func newGatewayServer() *server.HttpServer {
	hs := server.NewHttpServer()
	hs.SetRegisterFunc(func(ctx context.Context, mux *runtime.ServeMux) error {
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		addr := ":" + grpcPort()

		if err := hellopb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
			return err
		}
		if err := agentpb.RegisterAgentServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
			return err
		}
		if err := agentpb.RegisterStrategyServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
			return err
		}
		if err := agentpb.RegisterMappingServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
			return err
		}
		return nil
	})
	return hs
}

// grpcPort 从配置中提取 gRPC 端口号
func grpcPort() string {
	return strings.TrimPrefix(apccfg.GetConf().Server.GrpcAddr, ":")
}
