package logic

import (
	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	"github.com/ethereal3x/mint-server/internal/dto"
	"github.com/ethereal3x/mint-server/internal/model"
)

// GroupModelsByManufacturer 按厂商分组模型列表
func GroupModelsByManufacturer(configs []*model.ChatModelConfig) []*agentpb.ManufacturerGroup {
	groupMap := make(map[string][]*agentpb.ModelInfo)
	for _, config := range configs {
		pbConfig := dto.ConfigToProto(config)
		groupMap[config.Manufacturer] = append(groupMap[config.Manufacturer], &agentpb.ModelInfo{
			Model:              config.ModelType,
			Description:        config.Description,
			ModelCapabilities:  pbConfig.ModelCapabilities,
			SupportsMultimodal: config.SupportsMultimodal,
		})
	}
	manufacturers := make([]*agentpb.ManufacturerGroup, 0, len(groupMap))
	for name, models := range groupMap {
		manufacturers = append(manufacturers, &agentpb.ManufacturerGroup{
			Manufacturer: name,
			Models:       models,
		})
	}
	return manufacturers
}
