package model

import "time"

// 上传状态常量
const (
	UPLOAD_STATUS_PENDING   = "pending"
	UPLOAD_STATUS_UPLOADING = "uploading"
	UPLOAD_STATUS_COMPLETED = "completed"
	UPLOAD_STATUS_CANCELED  = "canceled"
	UPLOAD_STATUS_FAILED    = "failed"
)

// FileUploadQuery 文件上传查询条件
type FileUploadQuery struct {
	ID     int32
	UserID string
}

// FileUpload 对应 tb_file_uploads 表
type FileUpload struct {
	ID           int32     `gorm:"column:id;primaryKey;autoIncrement"`
	ObjectName   string    `gorm:"column:object_name;uniqueIndex;size:255;not null"`
	OriginalName string    `gorm:"column:original_name;size:255;not null"`
	FileSize     int64     `gorm:"column:file_size;default:0"`
	ContentType  string    `gorm:"column:content_type;size:100;not null"`
	URL          string    `gorm:"column:url;size:512;not null"`
	UploadID     string    `gorm:"column:upload_id;size:255"`
	Status       string    `gorm:"column:status;size:20;default:pending"`
	UploadedSize int64     `gorm:"column:uploaded_size;default:0"`
	UserID       string    `gorm:"column:user_id;size:255;not null"`
	CreatedTime  time.Time `gorm:"column:created_time;autoCreateTime"`
	UpdatedTime  time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

// TableName 指定表名
func (FileUpload) TableName() string {
	return "tb_file_uploads"
}
