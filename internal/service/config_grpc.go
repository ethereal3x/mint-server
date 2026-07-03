package service

import (
	"context"

	"github.com/ethereal3x/apc/errs"
	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
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

// ListConfigs 分页查询模型配置
func (s *ConfigServer) ListConfigs(ctx context.Context, req *agentpb.ListConfigsRequest) (*agentpb.ListConfigsResponse, error) {
	return errs.Handle(&agentpb.ListConfigsResponse{}, func(rsp *agentpb.ListConfigsResponse) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
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
			return err
		}
		configs := make([]*agentpb.ModelConfig, 0, len(list))
		for _, item := range list {
			configs = append(configs, dto.ConfigToProto(item))
		}
		rsp.Configs = configs
		rsp.Total = total
		return nil
	})
}

// GetConfig 按 ID 查询模型配置
func (s *ConfigServer) GetConfig(ctx context.Context, req *agentpb.GetConfigRequest) (*agentpb.ConfigReply, error) {
	return errs.Handle(&agentpb.ConfigReply{}, func(rsp *agentpb.ConfigReply) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		if req.Id <= 0 {
			return mint_err.ErrParam
		}
		config, err := s.logic.GetByID(ctx, req.Id, userID)
		if err != nil {
			return err
		}
		if config == nil {
			return mint_err.ErrConfigNotFound
		}
		rsp.Config = dto.ConfigToProto(config)
		return nil
	})
}

// CreateConfig 创建模型配置
func (s *ConfigServer) CreateConfig(ctx context.Context, req *agentpb.CreateConfigRequest) (*agentpb.ConfigReply, error) {
	return errs.Handle(&agentpb.ConfigReply{}, func(rsp *agentpb.ConfigReply) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		if req.ModelType == "" || req.ApiKey == "" {
			return mint_err.ErrParam
		}
		config := dto.CreateReqToModel(req)
		config.UserID = userID
		if err := s.logic.Create(ctx, config); err != nil {
			return err
		}
		rsp.Config = dto.ConfigToProto(config)
		return nil
	})
}

// UpdateConfig 更新模型配置
func (s *ConfigServer) UpdateConfig(ctx context.Context, req *agentpb.UpdateConfigRequest) (*agentpb.ConfigReply, error) {
	return errs.Handle(&agentpb.ConfigReply{}, func(rsp *agentpb.ConfigReply) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		if req.Id <= 0 {
			return mint_err.ErrParam
		}
		config := dto.UpdateReqToModel(req)
		if err := s.logic.Update(ctx, config, userID); err != nil {
			return err
		}
		updated, err := s.logic.GetByID(ctx, req.Id, userID)
		if err != nil {
			return err
		}
		rsp.Config = dto.ConfigToProto(updated)
		return nil
	})
}

// DeleteConfig 删除模型配置
func (s *ConfigServer) DeleteConfig(ctx context.Context, req *agentpb.DeleteConfigRequest) (*emptypb.Empty, error) {
	return errs.Handle(&emptypb.Empty{}, func(rsp *emptypb.Empty) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		if req.Id <= 0 {
			return mint_err.ErrParam
		}
		return s.logic.Delete(ctx, req.Id, userID)
	})
}

// GetModelStats 按模型聚合 token 消耗统计
func (s *ConfigServer) GetModelStats(ctx context.Context, _ *emptypb.Empty) (*agentpb.ModelStatsResponse, error) {
	return errs.Handle(&agentpb.ModelStatsResponse{}, func(rsp *agentpb.ModelStatsResponse) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		stats, err := s.chat.AggregateByModel(ctx, userID)
		if err != nil {
			return err
		}
		pbStats := make([]*agentpb.ModelStat, 0, len(stats))
		for _, stat := range stats {
			pbStats = append(pbStats, dto.ModelStatToProto(stat))
		}
		rsp.Stats = pbStats
		return nil
	})
}
