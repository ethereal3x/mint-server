package model

import "time"

// ChatModelConfig 对应 tb_model_config 表，合并策略规则与模型映射
type ChatModelConfig struct {
	ID                int32     `gorm:"column:id;primaryKey;autoIncrement"`
	ModelType         string    `gorm:"column:model_type;uniqueIndex;size:50;not null"`
	Manufacturer      string    `gorm:"column:manufacturer;size:50;not null"`
	Description       string    `gorm:"column:description;type:text"`
	InputPrice        float64   `gorm:"column:input_price;type:decimal(10,8);default:0"`
	OutputPrice       float64   `gorm:"column:output_price;type:decimal(10,8);default:0"`
	APIKey            string    `gorm:"column:api_key;size:512;not null"`
	URL               string    `gorm:"column:url;size:128"`
	MaxTokens         int32     `gorm:"column:max_tokens;default:0"`
	Stream            bool      `gorm:"column:stream;default:true"`
	Temperature       float32   `gorm:"column:temperature;default:0"`
	TopP              float32   `gorm:"column:top_p;default:0"`
	N                 int32     `gorm:"column:n;default:1"`
	PresencePenalty   float32   `gorm:"column:presence_penalty;default:0"`
	FrequencyPenalty  float32   `gorm:"column:frequency_penalty;default:0"`
	AgentGenerateType string    `gorm:"column:agent_generate_type;size:50"`
	Route             string    `gorm:"column:route;size:128"`
	IsEnabled         int32     `gorm:"column:is_enabled;default:1"`
	CreatedTime       time.Time `gorm:"column:created_time;autoCreateTime"`
	UpdatedTime       time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

// TableName 指定表名
func (ChatModelConfig) TableName() string {
	return "tb_model_config"
}
