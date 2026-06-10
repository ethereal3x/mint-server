package dto

import (
	"github.com/cloudwego/eino/schema"
	"github.com/ethereal3x/mint-server/internal/model"
)

// ChatRequest 聊天请求参数
type ChatRequest struct {
	UserID     string
	Question   string
	Model      string
	DialogueID string
	RecordID   string
	FileData   []byte
	FileName   string
	ImageURLs  []string
}

// SaveRecordRequest 保存对话记录参数
type SaveRecordRequest struct {
	ChatRequest *ChatRequest
	Answer      string
	Config      *model.ChatModelConfig
	Usage       *schema.TokenUsage
}

// ChatResult 聊天调用结果
type ChatResult struct {
	Config  *model.ChatModelConfig
	Content string
	Usage   *schema.TokenUsage
}
