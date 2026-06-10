package service

import (
	"context"

	"github.com/ethereal3x/apc/errs"
	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	"github.com/ethereal3x/mint-server/internal/auth"
	"github.com/ethereal3x/mint-server/internal/dto"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ConfigServer ModelConfigService gRPC 处理器
type ConfigServer struct {
	agentpb.UnimplementedModelConfigServiceServer
	logic ConfigCrudLogic
	chat  ConfigStatsLogic
}

// NewConfigServer 创建模型配置 gRPC 处理器
func NewConfigServer(configLogic ConfigCrudLogic, chat ConfigStatsLogic) *ConfigServer {
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
		configs = append(configs, dto.ConfigToProto(c))
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
	rsp.Config = dto.ConfigToProto(config)
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
	config := dto.CreateReqToModel(req)
	config.UserID = userID
	if err := s.logic.Create(ctx, config); err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	rsp.Config = dto.ConfigToProto(config)
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
	config := dto.UpdateReqToModel(req)
	if err := s.logic.Update(ctx, config, userID); err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	updated, err := s.logic.GetByID(ctx, req.Id, userID)
	if err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	rsp.Config = dto.ConfigToProto(updated)
	return rsp, nil
}

func (s *ConfigServer) DeleteConfig(ctx context.Context, req *agentpb.DeleteConfigRequest) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	userID, err := auth.RequireUserID(ctx)
	if err != nil {
		return errs.GenProtoReply(rsp, mint_err.ErrUnauthorized)
	}
	if req.Id <= 0 {
		return errs.GenProtoReply(rsp, mint_err.ErrParam)
	}
	if err := s.logic.Delete(ctx, req.Id, userID); err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	return rsp, nil
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
		pbStats = append(pbStats, dto.ModelStatToProto(st))
	}
	rsp.Stats = pbStats
	return rsp, nil
}
