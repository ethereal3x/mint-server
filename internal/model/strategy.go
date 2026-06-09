package model

import "time"

// StrategyRule 对应 tb_agent_strategy_rules 表
type StrategyRule struct {
	ID               int32     `gorm:"column:id;primaryKey;autoIncrement"`
	RuleID           string    `gorm:"column:rule_id;uniqueIndex;size:50;not null"`
	APIKey           string    `gorm:"column:api_key;size:128;not null"`
	AgentModel       string    `gorm:"column:agent_model;size:50"`
	AgentManufacturer string   `gorm:"column:agent_manufacturer;size:50"`
	AgentGenerateType string   `gorm:"column:agent_generate_type;size:50"`
	URL              string    `gorm:"column:url;size:128"`
	MaxTokens        int32     `gorm:"column:max_tokens;default:0"`
	Stream           bool      `gorm:"column:stream;default:true"`
	Temperature      float32   `gorm:"column:temperature;default:0"`
	TopP             float32   `gorm:"column:top_p;default:0"`
	N                int32     `gorm:"column:n;default:1"`
	PresencePenalty  float32   `gorm:"column:presence_penalty;default:0"`
	FrequencyPenalty float32   `gorm:"column:frequency_penalty;default:0"`
	Route            string    `gorm:"column:route;size:128"`
	IsEnabled        int32     `gorm:"column:is_enabled;default:1"`
	CreatedTime      time.Time `gorm:"column:created_time;autoCreateTime"`
	UpdatedTime      time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

// TableName 指定表名
func (StrategyRule) TableName() string {
	return "tb_agent_strategy_rules"
}
