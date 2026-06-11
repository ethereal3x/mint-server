package dto

import (
	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	"github.com/ethereal3x/mint-server/internal/model"
)

// configProtoMessage Create 和 Update 配置请求共有的 getter 方法集合
type configProtoMessage interface {
	GetModelType() string
	GetManufacturer() string
	GetDescription() string
	GetInputPrice() float64
	GetOutputPrice() float64
	GetApiKey() string
	GetUrl() string
	GetMaxTokens() int32
	GetStream() bool
	GetTemperature() float32
	GetTopP() float32
	GetN() int32
	GetPresencePenalty() float32
	GetFrequencyPenalty() float32
	GetAgentGenerateType() string
	GetRoute() string
	GetIsEnabled() int32
	GetSupportsMultimodal() bool
}

// baseConfigFromProto 从 proto 请求消息提取公共字段到 ChatModelConfig
func baseConfigFromProto(src configProtoMessage) *model.ChatModelConfig {
	return &model.ChatModelConfig{
		ModelType:          src.GetModelType(),
		Manufacturer:       src.GetManufacturer(),
		Description:        src.GetDescription(),
		InputPrice:         src.GetInputPrice(),
		OutputPrice:        src.GetOutputPrice(),
		APIKey:             src.GetApiKey(),
		URL:                src.GetUrl(),
		MaxTokens:          src.GetMaxTokens(),
		Stream:             src.GetStream(),
		Temperature:        src.GetTemperature(),
		TopP:               src.GetTopP(),
		N:                  src.GetN(),
		PresencePenalty:    src.GetPresencePenalty(),
		FrequencyPenalty:   src.GetFrequencyPenalty(),
		AgentGenerateType:  src.GetAgentGenerateType(),
		Route:              src.GetRoute(),
		IsEnabled:          src.GetIsEnabled(),
		SupportsMultimodal: src.GetSupportsMultimodal(),
	}
}

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
	return baseConfigFromProto(req)
}

// UpdateReqToModel 转换更新请求 proto 为模型
func UpdateReqToModel(req *agentpb.UpdateConfigRequest) *model.ChatModelConfig {
	cfg := baseConfigFromProto(req)
	cfg.ID = req.GetId()
	return cfg
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
