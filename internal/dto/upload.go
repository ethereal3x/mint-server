package dto

import "io"

// UploadRequest 文件上传请求参数
type UploadRequest struct {
	UserID      string
	FileName    string
	Reader      io.Reader
	Size        int64
	ContentType string
}

// UploadResult 上传结果
type UploadResult struct {
	ID       int32  `json:"id"`
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}
