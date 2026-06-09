package model

import "time"

// ModelMapping 对应 tb_model_manufacturer_mapping 表
type ModelMapping struct {
	ID           int32     `gorm:"column:id;primaryKey;autoIncrement"`
	ModelType    string    `gorm:"column:model_type;uniqueIndex;size:50;not null"`
	Manufacturer string    `gorm:"column:manufacturer;size:50;not null"`
	Description  string    `gorm:"column:description;type:text"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (ModelMapping) TableName() string {
	return "tb_model_manufacturer_mapping"
}
