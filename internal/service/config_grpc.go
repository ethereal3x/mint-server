package service

import (
	"context"

	"github.com/ethereal3x/apc/errs"
	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	"github.com/ethereal3x/mint-server/internal/auth"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ConfigServer ModelConfigService gRPC 处理器
type ConfigServer struct {
	agentpb.UnimplementedModelConfigServiceServer
	logic *logic.Config
	chat  *logic.Chat
}

// NewConfigServer 创建模型配置 gRPC 处理器
func NewConfigServer(configLogic *logic.Config, chat *logic.Chat) *ConfigServer {
	return &ConfigServer{logic: configLogic, chat: chat}
}

func (s *ConfigServer) ListConfigs(ctx context.Context, req *agentpb.ListConfigsRequest) (*agentpb.ListConfigsResponse, error) {
	rsp := &agentpb.ListConfigsResponse{}
	userID, err := auth.RequireUserID(ctx)
	if err != nil {
		return errs.GenProtoReply(rsp, mint_err.ErrUnauthorized)
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	list, total, err := s.logic.List(ctx, page, pageSize, userID)
	if err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	configs := make([]*agentpb.ModelConfig, 0, len(list))
	for _, c := range list {
		configs = append(configs, configToProto(c))
	}
	rsp.Configs = configs
	rsp.Total = total
	return rsp, nil
}

func (s *ConfigServer) GetConfig(ctx context.Context, req *agentpb.GetConfigRequest) (*agentpb.ConfigReply, error) {
	rsp := &agentpb.ConfigReply{}
	userID, err := auth.RequireUserID(ctx)
	if err != nil {
		return errs.GenProtoReply(rsp, mint_err.ErrUnauthorized)
	}
	if req.Id <= 0 {
		return errs.GenProtoReply(rsp, mint_err.ErrParam)
	}
	config, err := s.logic.GetByID(ctx, req.Id, userID)
	if err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	if config == nil {
		return errs.GenProtoReply(rsp, mint_err.ErrConfigNotFound)
	}
	rsp.Config = configToProto(config)
	return rsp, nil
}

func (s *ConfigServer) CreateConfig(ctx context.Context, req *agentpb.CreateConfigRequest) (*agentpb.ConfigReply, error) {
	rsp := &agentpb.ConfigReply{}
	userID, err := auth.RequireUserID(ctx)
	if err != nil {
		return errs.GenProtoReply(rsp, mint_err.ErrUnauthorized)
	}
	if req.ModelType == "" || req.ApiKey == "" {
		return errs.GenProtoReply(rsp, mint_err.ErrParam)
	}
	config := createReqToModel(req)
	config.UserID = userID
	if err := s.logic.Create(ctx, config); err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	rsp.Config = configToProto(config)
	return rsp, nil
}

func (s *ConfigServer) UpdateConfig(ctx context.Context, req *agentpb.UpdateConfigRequest) (*agentpb.ConfigReply, error) {
	rsp := &agentpb.ConfigReply{}
	userID, err := auth.RequireUserID(ctx)
	if err != nil {
		return errs.GenProtoReply(rsp, mint_err.ErrUnauthorized)
	}
	if req.Id <= 0 {
		return errs.GenProtoReply(rsp, mint_err.ErrParam)
	}
	config := updateReqToModel(req)
	if err := s.logic.Update(ctx, config, userID); err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	updated, err := s.logic.GetByID(ctx, req.Id, userID)
	if err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	rsp.Config = configToProto(updated)
	return rsp, nil
}

func (s *ConfigServer) DeleteConfig(ctx context.Context, req *agentpb.DeleteConfigRequest) (*emptypb.Empty, error) {
	userID, err := auth.RequireUserID(ctx)
	if err != nil {
		return nil, nil
	}
	if req.Id <= 0 {
		return nil, nil
	}
	if err := s.logic.Delete(ctx, req.Id, userID); err != nil {
		return nil, nil
	}
	return &emptypb.Empty{}, nil
}

func configToProto(c *model.ChatModelConfig) *agentpb.ModelConfig {
	if c == nil {
		return nil
	}
	return &agentpb.ModelConfig{
		Id:                c.ID,
		ModelType:         c.ModelType,
		Manufacturer:      c.Manufacturer,
		Description:       c.Description,
		InputPrice:        c.InputPrice,
		OutputPrice:       c.OutputPrice,
		ApiKey:            c.APIKey,
		Url:               c.URL,
		MaxTokens:         c.MaxTokens,
		Stream:            c.Stream,
		Temperature:       c.Temperature,
		TopP:              c.TopP,
		N:                 c.N,
		PresencePenalty:   c.PresencePenalty,
		FrequencyPenalty:  c.FrequencyPenalty,
		AgentGenerateType: c.AgentGenerateType,
		Route:             c.Route,
		IsEnabled:         c.IsEnabled,
	}
}

func createReqToModel(req *agentpb.CreateConfigRequest) *model.ChatModelConfig {
	return &model.ChatModelConfig{
		ModelType:         req.ModelType,
		Manufacturer:      req.Manufacturer,
		Description:       req.Description,
		InputPrice:        req.InputPrice,
		OutputPrice:       req.OutputPrice,
		APIKey:            req.ApiKey,
		URL:               req.Url,
		MaxTokens:         req.MaxTokens,
		Stream:            req.Stream,
		Temperature:       req.Temperature,
		TopP:              req.TopP,
		N:                 req.N,
		PresencePenalty:   req.PresencePenalty,
		FrequencyPenalty:  req.FrequencyPenalty,
		AgentGenerateType: req.AgentGenerateType,
		Route:             req.Route,
		IsEnabled:         req.IsEnabled,
	}
}

func (s *ConfigServer) GetModelStats(ctx context.Context, _ *emptypb.Empty) (*agentpb.ModelStatsResponse, error) {
	rsp := &agentpb.ModelStatsResponse{}
	userID, err := auth.RequireUserID(ctx)
	if err != nil {
		return errs.GenProtoReply(rsp, mint_err.ErrUnauthorized)
	}
	stats, err := s.chat.AggregateByModel(ctx, userID)
	if err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	pbStats := make([]*agentpb.ModelStat, 0, len(stats))
	for _, st := range stats {
		pbStats = append(pbStats, &agentpb.ModelStat{
			Model:             st.Model,
			TotalInputTokens:  st.TotalInputTokens,
			TotalOutputTokens: st.TotalOutputTokens,
			TotalInputCost:    st.TotalInputCost,
			TotalOutputCost:   st.TotalOutputCost,
		})
	}
	rsp.Stats = pbStats
	return rsp, nil
}

func updateReqToModel(req *agentpb.UpdateConfigRequest) *model.ChatModelConfig {
	return &model.ChatModelConfig{
		ID:                req.Id,
		ModelType:         req.ModelType,
		Manufacturer:      req.Manufacturer,
		Description:       req.Description,
		InputPrice:        req.InputPrice,
		OutputPrice:       req.OutputPrice,
		APIKey:            req.ApiKey,
		URL:               req.Url,
		MaxTokens:         req.MaxTokens,
		Stream:            req.Stream,
		Temperature:       req.Temperature,
		TopP:              req.TopP,
		N:                 req.N,
		PresencePenalty:   req.PresencePenalty,
		FrequencyPenalty:  req.FrequencyPenalty,
		AgentGenerateType: req.AgentGenerateType,
		Route:             req.Route,
		IsEnabled:         req.IsEnabled,
	}
}
