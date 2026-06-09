package model

import "time"

// DialogueRecord 对应 tb_user_agent_dialogues 表
type DialogueRecord struct {
	DialogueID   string    `gorm:"column:dialogue_id;size:255;not null;index:idx_tb_user_agent_dialogues_dialogue_id"`
	RecordID     string    `gorm:"column:record_id;size:255;primaryKey"`
	UserID       string    `gorm:"column:user_id;size:255;not null;index:idx_tb_user_agent_dialogues_user_id"`
	Model        string    `gorm:"column:model;size:255;not null"`
	UserContent  string    `gorm:"column:user_content;type:text"`
	AgentContent string    `gorm:"column:agent_content;type:text"`
	TotalToken   int64     `gorm:"column:total_tokens;default:0"`
	UserToken    int64     `gorm:"column:user_tokens;default:0"`
	AgentToken   int64     `gorm:"column:agent_tokens;default:0"`
	CreatedTime  time.Time `gorm:"column:created_time;autoCreateTime"`
	UpdatedTime  time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

// TableName 指定表名
func (DialogueRecord) TableName() string {
	return "tb_user_agent_dialogues"
}

// DialogueSummary 对话摘要，用于侧边栏列表展示
type DialogueSummary struct {
	DialogueID   string
	Title        string
	MessageCount int32
	UpdatedTime  string
}
