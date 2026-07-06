package app

import (
	"context"
	"fmt"
	"strings"

	apccfg "github.com/ethereal3x/apc/config"
	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/orm"
	"github.com/ethereal3x/apc/storage"
	"github.com/ethereal3x/mint-server/internal/config"
	"github.com/ethereal3x/mint-server/internal/util"
	"gorm.io/gorm"
)

// App 持有应用全部基础设施和领域模块
type App struct {
	DB           *gorm.DB
	Storage      storage.ObjectStorage
	TokenManager *util.TokenManager
	BizConfig    config.BusinessConfig

	Agent  *AgentModule
	Authx  *AuthModule
	Hello  *HelloModule
	Upload *UploadModule
}

// New 初始化全部依赖并返回装配好的 App
func New(cfgPath string) (*App, error) {
	if err := apccfg.Load(apccfg.LoadOptions{Path: cfgPath}); err != nil {
		return nil, fmt.Errorf("load apc config: %w", err)
	}
	applyTracingEnvOverrides(apccfg.GetConf())
	if err := config.Load(config.LoadOptions{Path: cfgPath}); err != nil {
		return nil, fmt.Errorf("load business config: %w", err)
	}
	bizCfg := config.GetBusinessConfig()
	if err := config.Validate(bizCfg); err != nil {
		return nil, fmt.Errorf("validate business config: %w", err)
	}

	apcConf := apccfg.GetConf()
	if err := validateInfraConfig(apcConf); err != nil {
		return nil, fmt.Errorf("validate infra config: %w", err)
	}

	logger.SetLogger(logger.NewLogger(&apcConf.Plugin.Log))

	db, err := initDB()
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	store, err := initStorage()
	if err != nil {
		return nil, fmt.Errorf("init storage: %w", err)
	}

	tokenManager := util.NewTokenManager(util.TokenManagerConfig{
		Secret: bizCfg.JWT.Secret,
	})
	secretKey := []byte(bizCfg.SecretKey.Encryption)

	application := &App{
		DB:           db,
		Storage:      store,
		TokenManager: tokenManager,
		BizConfig:    bizCfg,
	}
	application.Agent = newAgentModule(db, secretKey)
	application.Authx = newAuthModule(db, tokenManager)
	application.Hello = newHelloModule()
	application.Upload = newUploadModule(db, store, tokenManager)

	return application, nil
}

// validateInfraConfig 校验基础设施配置必填项
func validateInfraConfig(cfg *apccfg.Config) error {
	if cfg == nil {
		return fmt.Errorf("apc config is nil")
	}
	if strings.TrimSpace(cfg.Server.GrpcAddr) == "" {
		return fmt.Errorf("server.grpc_addr is required")
	}
	if strings.TrimSpace(cfg.Server.GatewayAddr) == "" {
		return fmt.Errorf("server.gateway_addr is required")
	}
	if apccfg.GetClientConf("mysql") == nil {
		return fmt.Errorf("client mysql config not found")
	}
	rustfs := cfg.Plugin.RustFS
	if strings.TrimSpace(rustfs.Endpoint) == "" {
		return fmt.Errorf("plugin.rustfs.endpoint is required")
	}
	if strings.TrimSpace(rustfs.AccessKey) == "" || strings.TrimSpace(rustfs.SecretKey) == "" {
		return fmt.Errorf("plugin.rustfs access_key and secret_key are required")
	}
	if strings.TrimSpace(rustfs.BucketName) == "" {
		return fmt.Errorf("plugin.rustfs.bucket_name is required")
	}
	return nil
}

// initDB 根据配置初始化 GORM 数据库连接
func initDB() (*gorm.DB, error) {
	clientConf := apccfg.GetClientConf("mysql")
	gormCfg := apccfg.GenGormConfig(clientConf)
	db, err := orm.NewGormInstance(gormCfg)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// initStorage 初始化 rustfs 对象存储客户端
func initStorage() (storage.ObjectStorage, error) {
	cfg := apccfg.GetConf()
	storageClient, err := storage.NewS3Client(context.Background(), storage.NewS3ClientParams{
		Provider: storage.STORAGE_PROVIDER_RUSTFS,
		Config:   &cfg.Plugin.RustFS,
	})
	if err != nil {
		return nil, err
	}
	return storageClient, nil
}
