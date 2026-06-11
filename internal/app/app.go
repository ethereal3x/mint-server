package app

import (
	"context"
	"fmt"

	apccfg "github.com/ethereal3x/apc/config"
	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/orm"
	"github.com/ethereal3x/apc/storage"
	"github.com/ethereal3x/mint-server/internal/config"
	"github.com/ethereal3x/mint-server/internal/util"
	"gorm.io/gorm"
)

// App 持有应用全部基础设施和领域模块，替代旧的 globalDB、globalServices 等全局变量。
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

// New 初始化全部依赖并返回装配好的 App。调用方通过此入口替代旧的 initApp。
func New(cfgPath string) (*App, error) {
	if err := config.InitBusinessConfig(cfgPath); err != nil {
		return nil, fmt.Errorf("load business config: %w", err)
	}
	bizCfg := config.GetBusinessConfig()

	apcConf := apccfg.GetConf()
	logger.SetLogger(logger.NewLogger(&apcConf.Plugin.Log))

	db, err := initDB()
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}
	if db == nil {
		return nil, fmt.Errorf("init db: db is nil")
	}

	store, err := initStorage()
	if err != nil {
		return nil, fmt.Errorf("init storage: %w", err)
	}

	tokenManager := util.NewTokenManager(util.TokenManagerConfig{
		Secret: bizCfg.JWT.Secret,
	})
	secretKey := []byte(bizCfg.SecretKey.Encryption)

	app := &App{
		DB:           db,
		Storage:      store,
		TokenManager: tokenManager,
		BizConfig:    bizCfg,
	}
	app.Agent = newAgentModule(db, secretKey)
	app.Authx = newAuthModule(db, tokenManager)
	app.Hello = newHelloModule()
	app.Upload = newUploadModule(db, store, tokenManager)

	return app, nil
}

// initDB 根据配置初始化 GORM 数据库连接
func initDB() (*gorm.DB, error) {
	clientConf := apccfg.GetClientConf("mysql")
	if clientConf == nil {
		return nil, fmt.Errorf("mysql client config not found")
	}
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
