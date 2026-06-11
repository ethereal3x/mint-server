package app

import (
	"github.com/ethereal3x/mint-server/internal/decorator"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/repo"
	"github.com/ethereal3x/mint-server/internal/service"
	"gorm.io/gorm"
)

// AgentModule 聚合聊天和模型配置的 gRPC 服务处理器
type AgentModule struct {
	AgentServer  *service.AgentServer
	ConfigServer *service.ConfigServer
}

// newAgentModule 装配 Agent 领域全部依赖，通过 decorator 注入 tracing/logging
func newAgentModule(db *gorm.DB, secretKey []byte) *AgentModule {
	configRepo := repo.NewModelConfigRepo(db, secretKey)
	recordRepo := repo.NewDialogueRepo(db)

	configLogic := logic.NewConfig(configRepo)
	chatLogic := logic.NewChat(configRepo, recordRepo, logic.NewEinoAdapter())

	chatDecorated := decorator.NewChatDecorator(chatLogic)
	configDecorated := decorator.NewConfigDecorator(configLogic)

	return &AgentModule{
		AgentServer:  service.NewAgentServer(chatDecorated, configDecorated),
		ConfigServer: service.NewConfigServer(configDecorated, chatDecorated),
	}
}
