package model

import "time"

const (
	AUTH_PROVIDER_ACCOUNT_PASSWORD = "account_password"
)

// BaseUser 对应 tbl_base_user 表
type BaseUser struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID        int64     `gorm:"column:user_id;default:0;not null;uniqueIndex"`
	Username      string    `gorm:"column:username;size:64;not null;uniqueIndex"`
	Nickname      string    `gorm:"column:nickname;size:20;not null"`
	AvatarURL     string    `gorm:"column:avatar_url;size:512;not null"`
	Realname      string    `gorm:"column:realname;size:20;not null"`
	Password      string    `gorm:"column:password;type:char(60);not null"`
	Password2     string    `gorm:"column:password2;size:60;not null"`
	RegTime       int64     `gorm:"column:reg_time;default:0;not null"`
	Mobile        string    `gorm:"column:mobile;size:20;not null"`
	LastLoginIP   string    `gorm:"column:last_login_ip;size:45;not null"`
	LastLoginTime int64     `gorm:"column:last_login_time;default:0;not null"`
	CreateAt      time.Time `gorm:"column:create_at;autoCreateTime"`
	UpdateAt      time.Time `gorm:"column:update_at;autoUpdateTime"`
}

// TableName 指定基础用户表名
func (BaseUser) TableName() string {
	return "tbl_base_user"
}

// ThirdUser 对应 tbl_third_user 表
type ThirdUser struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID      int64     `gorm:"column:user_id;default:0;not null"`
	ChannelCode string    `gorm:"column:channel_code;size:20;not null"`
	UnionID     string    `gorm:"column:union_id;size:256"`
	OpenID      string    `gorm:"column:open_id;size:256"`
	CreateAt    time.Time `gorm:"column:create_at;autoCreateTime"`
	UpdateAt    time.Time `gorm:"column:update_at;autoUpdateTime"`
}

// TableName 指定第三方用户关联表名
func (ThirdUser) TableName() string {
	return "tbl_third_user"
}

// BaseUserQuery 基础用户查询条件
type BaseUserQuery struct {
	UserID   int64
	Username string
}
