package service

import (
	"context"

	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// StrategyServer StrategyService gRPC 处理器
type StrategyServer struct {
	agentpb.UnimplementedStrategyServiceServer
	logic *logic.Strategy
}

// NewStrategyServer 创建 Strategy gRPC 处理器
func NewStrategyServer(strategyLogic *logic.Strategy) *StrategyServer {
	return &StrategyServer{logic: strategyLogic}
}

func (s *StrategyServer) ListStrategies(ctx context.Context, req *agentpb.ListStrategiesRequest) (*agentpb.ListStrategiesResponse, error) {
	page := max(int(req.Page), 1)
	pageSize := clamp(int(req.PageSize), 1, 100, 20)
	offset := (page - 1) * pageSize

	list, total, err := s.logic.List(ctx, offset, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list strategies: %v", err)
	}
	strategies := make([]*agentpb.Strategy, 0, len(list))
	for _, rule := range list {
		strategies = append(strategies, protoFromRule(rule))
	}
	return &agentpb.ListStrategiesResponse{Strategies: strategies, Total: total}, nil
}

func (s *StrategyServer) GetStrategy(ctx context.Context, req *agentpb.GetStrategyRequest) (*agentpb.Strategy, error) {
	if req.RuleId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "get strategy: missing rule_id")
	}
	rule, err := s.logic.GetByRuleID(ctx, req.RuleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get strategy: %v", err)
	}
	if rule == nil {
		return nil, status.Errorf(codes.NotFound, "get strategy: not found")
	}
	return protoFromRule(rule), nil
}

func (s *StrategyServer) CreateStrategy(ctx context.Context, req *agentpb.CreateStrategyRequest) (*agentpb.Strategy, error) {
	if req.ApiKey == "" || req.AgentManufacturer == "" {
		return nil, status.Errorf(codes.InvalidArgument, "create strategy: missing api_key or agent_manufacturer")
	}
	rule := ruleFromCreateReq(req)
	if err := s.logic.Create(ctx, rule); err != nil {
		return nil, status.Errorf(codes.Internal, "create strategy: %v", err)
	}
	return protoFromRule(rule), nil
}

func (s *StrategyServer) UpdateStrategy(ctx context.Context, req *agentpb.UpdateStrategyRequest) (*agentpb.Strategy, error) {
	if req.RuleId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "update strategy: missing rule_id")
	}
	rule := ruleFromUpdateReq(req)
	if err := s.logic.Update(ctx, rule); err != nil {
		return nil, status.Errorf(codes.Internal, "update strategy: %v", err)
	}
	return protoFromRule(rule), nil
}

func (s *StrategyServer) DeleteStrategy(ctx context.Context, req *agentpb.DeleteStrategyRequest) (*emptypb.Empty, error) {
	if req.RuleId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "delete strategy: missing rule_id")
	}
	if err := s.logic.Delete(ctx, req.RuleId); err != nil {
		return nil, status.Errorf(codes.Internal, "delete strategy: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func protoFromRule(rule *model.StrategyRule) *agentpb.Strategy {
	apiKey := rule.APIKey
	if len(apiKey) > 4 {
		apiKey = apiKey[:4] + "****"
	} else if apiKey != "" {
		apiKey = "****"
	}
	return &agentpb.Strategy{
		Id:                 rule.ID,
		RuleId:             rule.RuleID,
		ApiKey:             apiKey,
		AgentModel:         rule.AgentModel,
		AgentManufacturer:  rule.AgentManufacturer,
		AgentGenerateType:  rule.AgentGenerateType,
		Url:                rule.URL,
		MaxTokens:          rule.MaxTokens,
		Stream:             rule.Stream,
		Temperature:        rule.Temperature,
		TopP:               rule.TopP,
		N:                  rule.N,
		PresencePenalty:    rule.PresencePenalty,
		FrequencyPenalty:   rule.FrequencyPenalty,
		Route:              rule.Route,
		IsEnabled:          rule.IsEnabled,
	}
}

func ruleFromCreateReq(req *agentpb.CreateStrategyRequest) *model.StrategyRule {
	return &model.StrategyRule{
		APIKey:            req.ApiKey,
		AgentModel:        req.AgentModel,
		AgentManufacturer: req.AgentManufacturer,
		AgentGenerateType: req.AgentGenerateType,
		URL:               req.Url,
		MaxTokens:         req.MaxTokens,
		Stream:            req.Stream,
		Temperature:       req.Temperature,
		TopP:              req.TopP,
		N:                 req.N,
		PresencePenalty:   req.PresencePenalty,
		FrequencyPenalty:  req.FrequencyPenalty,
		Route:             req.Route,
		IsEnabled:         1,
	}
}

func ruleFromUpdateReq(req *agentpb.UpdateStrategyRequest) *model.StrategyRule {
	return &model.StrategyRule{
		RuleID:            req.RuleId,
		APIKey:            req.ApiKey,
		AgentModel:        req.AgentModel,
		AgentManufacturer: req.AgentManufacturer,
		AgentGenerateType: req.AgentGenerateType,
		URL:               req.Url,
		MaxTokens:         req.MaxTokens,
		Stream:            req.Stream,
		Temperature:       req.Temperature,
		TopP:              req.TopP,
		N:                 req.N,
		PresencePenalty:   req.PresencePenalty,
		FrequencyPenalty:  req.FrequencyPenalty,
		Route:             req.Route,
		IsEnabled:         req.IsEnabled,
	}
}

func clamp(v, min, max, def int) int {
	if v < min || v > max {
		return def
	}
	return v
}
