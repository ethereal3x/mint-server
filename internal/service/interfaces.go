package service

import (
	"context"

	"github.com/ethereal3x/mint-server/internal/dto"
	"github.com/ethereal3x/mint-server/internal/model"
)

// AgentChatLogic AgentServer 需要的聊天业务接口
type AgentChatLogic interface {
	StreamChat(ctx context.Context, req *dto.ChatRequest, contentChan chan<- string) (*dto.ChatResult, error)
	GenerateChat(ctx context.Context, req *dto.ChatRequest) (*dto.ChatResult, error)
	SaveRecord(ctx context.Context, req *dto.SaveRecordRequest) error
	GetHistory(ctx context.Context, query *model.DialogueQuery) ([]*model.DialogueRecord, error)
	ListDialogues(ctx context.Context, userID string) ([]*model.DialogueSummary, error)
}

// AgentConfigLogic AgentServer 需要的模型配置查询接口
type AgentConfigLogic interface {
	ListAll(ctx context.Context, userID string) ([]*model.ChatModelConfig, error)
}

// ConfigCrudLogic ConfigServer 需要的配置 CRUD 接口
type ConfigCrudLogic interface {
	List(ctx context.Context, page, pageSize int32, userID string) ([]*model.ChatModelConfig, int64, error)
	GetByID(ctx context.Context, id int32, userID string) (*model.ChatModelConfig, error)
	Create(ctx context.Context, config *model.ChatModelConfig) error
	Update(ctx context.Context, config *model.ChatModelConfig, userID string) error
	Delete(ctx context.Context, id int32, userID string) error
}

// ConfigStatsLogic ConfigServer 需要的模型统计接口
type ConfigStatsLogic interface {
	AggregateByModel(ctx context.Context, userID string) ([]*model.ModelStat, error)
}

// AuthServiceLogic AuthServer 需要的认证业务接口
type AuthServiceLogic interface {
	RegisterAccount(ctx context.Context, req *dto.RegisterAccountRequest) (*dto.AuthResult, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResult, error)
}

// UserServiceLogic UserServer 需要的用户操作接口
type UserServiceLogic interface {
	GetMe(ctx context.Context, userID string) (*model.BaseUser, error)
	UpdateAvatar(ctx context.Context, userID string, avatarURL string) (*model.BaseUser, error)
	UpdatePassword(ctx context.Context, userID string, oldPassword string, newPassword string) error
	UpdateNickname(ctx context.Context, userID string, nickname string) (*model.BaseUser, error)
}

// UploadServiceLogic UploadHandler 需要的文件上传接口
type UploadServiceLogic interface {
	Upload(ctx context.Context, req *dto.UploadRequest) (*dto.UploadResult, error)
	GetUpload(ctx context.Context, query *model.FileUploadQuery) (*model.FileUpload, error)
	ListUploads(ctx context.Context, userID string) ([]*model.FileUpload, error)
}
