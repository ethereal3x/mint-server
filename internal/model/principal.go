package model

// Principal 表示当前认证主体
type Principal struct {
	UserID     string
	Provider   string
	Identifier string
}
