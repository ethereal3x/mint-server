package dto

import (
	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	"github.com/ethereal3x/mint-server/internal/model"
)

// ConfigToProto 转换模型配置为 proto
func ConfigToProto(c *model.ChatModelConfig) *agentpb.ModelConfig {
	if c == nil {
		return nil
	}
	return &agentpb.ModelConfig{
		Id:                 c.ID,
		ModelType:          c.ModelType,
		Manufacturer:       c.Manufacturer,
		Description:        c.Description,
		InputPrice:         c.InputPrice,
		OutputPrice:        c.OutputPrice,
		ApiKey:             c.APIKey,
		Url:                c.URL,
		MaxTokens:          c.MaxTokens,
		Stream:             c.Stream,
		Temperature:        c.Temperature,
		TopP:               c.TopP,
		N:                  c.N,
		PresencePenalty:    c.PresencePenalty,
		FrequencyPenalty:   c.FrequencyPenalty,
		AgentGenerateType:  c.AgentGenerateType,
		Route:              c.Route,
		IsEnabled:          c.IsEnabled,
		SupportsMultimodal: c.SupportsMultimodal,
	}
}

// CreateReqToModel 转换创建请求 proto 为模型
func CreateReqToModel(req *agentpb.CreateConfigRequest) *model.ChatModelConfig {
	return &model.ChatModelConfig{
		ModelType:          req.ModelType,
		Manufacturer:       req.Manufacturer,
		Description:        req.Description,
		InputPrice:         req.InputPrice,
		OutputPrice:        req.OutputPrice,
		APIKey:             req.ApiKey,
		URL:                req.Url,
		MaxTokens:          req.MaxTokens,
		Stream:             req.Stream,
		Temperature:        req.Temperature,
		TopP:               req.TopP,
		N:                  req.N,
		PresencePenalty:    req.PresencePenalty,
		FrequencyPenalty:   req.FrequencyPenalty,
		AgentGenerateType:  req.AgentGenerateType,
		Route:              req.Route,
		IsEnabled:          req.IsEnabled,
		SupportsMultimodal: req.SupportsMultimodal,
	}
}

// UpdateReqToModel 转换更新请求 proto 为模型
func UpdateReqToModel(req *agentpb.UpdateConfigRequest) *model.ChatModelConfig {
	return &model.ChatModelConfig{
		ID:                 req.Id,
		ModelType:          req.ModelType,
		Manufacturer:       req.Manufacturer,
		Description:        req.Description,
		InputPrice:         req.InputPrice,
		OutputPrice:        req.OutputPrice,
		APIKey:             req.ApiKey,
		URL:                req.Url,
		MaxTokens:          req.MaxTokens,
		Stream:             req.Stream,
		Temperature:        req.Temperature,
		TopP:               req.TopP,
		N:                  req.N,
		PresencePenalty:    req.PresencePenalty,
		FrequencyPenalty:   req.FrequencyPenalty,
		AgentGenerateType:  req.AgentGenerateType,
		Route:              req.Route,
		IsEnabled:          req.IsEnabled,
		SupportsMultimodal: req.SupportsMultimodal,
	}
}

// ModelStatToProto 转换模型统计为 proto
func ModelStatToProto(st *model.ModelStat) *agentpb.ModelStat {
	if st == nil {
		return nil
	}
	return &agentpb.ModelStat{
		Model:             st.Model,
		TotalInputTokens:  st.TotalInputTokens,
		TotalOutputTokens: st.TotalOutputTokens,
		TotalInputCost:    st.TotalInputCost,
		TotalOutputCost:   st.TotalOutputCost,
	}
}
