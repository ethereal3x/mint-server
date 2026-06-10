package logic

import (
	"github.com/ethereal3x/mint-server/internal/model"
	"github.com/cloudwego/eino/schema"
)

const (
	documentPromptPrefix = "以下是用户传输文件内容，根据用户提问结合文件内容回答："
	defaultMaxHistory    = 20 // 滑动窗口默认保留最近 20 轮对话
)

// ContextBuilder 对话上下文构建器，支持链式调用组装 prompt
type ContextBuilder struct {
	systemPrompt     string
	history          []*model.DialogueRecord
	question         string
	imageURLs        []string
	fileData         []byte
	fileName         string
	maxHistoryRounds int
}

// NewContextBuilder 创建上下文构建器，默认保留最近 20 轮历史
func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{maxHistoryRounds: defaultMaxHistory}
}

// WithMaxHistory 设置滑动窗口大小（轮数），0 表示不限制
func (b *ContextBuilder) WithMaxHistory(rounds int) *ContextBuilder {
	b.maxHistoryRounds = rounds
	return b
}

// AddSystem 添加系统提示词
func (b *ContextBuilder) AddSystem(prompt string) *ContextBuilder {
	b.systemPrompt = prompt
	return b
}

// AddHistory 添加历史对话记录
func (b *ContextBuilder) AddHistory(records []*model.DialogueRecord) *ContextBuilder {
	b.history = records
	return b
}

// AddUserQuestion 添加当前用户问题
func (b *ContextBuilder) AddUserQuestion(q string) *ContextBuilder {
	b.question = q
	return b
}

// AddImages 添加图片 URL 列表，用于多模态对话
func (b *ContextBuilder) AddImages(urls []string) *ContextBuilder {
	b.imageURLs = urls
	return b
}

// AddFile 添加文件上下文
func (b *ContextBuilder) AddFile(data []byte, name string) *ContextBuilder {
	b.fileData = data
	b.fileName = name
	return b
}

// Build 组装最终的 Eino Message 列表
func (b *ContextBuilder) Build() []*schema.Message {
	var messages []*schema.Message

	if b.systemPrompt != "" {
		messages = append(messages, schema.SystemMessage(b.systemPrompt))
	}
	if len(b.fileData) > 0 && b.fileName != "" {
		messages = append(messages, schema.SystemMessage(documentPromptPrefix+string(b.fileData)))
	}

	history := b.history
	if b.maxHistoryRounds > 0 && len(history) > b.maxHistoryRounds {
		history = history[len(history)-b.maxHistoryRounds:]
	}
	for _, record := range history {
		messages = append(messages, schema.UserMessage(record.UserContent), schema.AssistantMessage(record.AgentContent, nil))
	}

	if b.question != "" {
		messages = append(messages, buildUserMessage(b.question, b.imageURLs))
	}
	return messages
}

// buildUserMessage 根据是否有图片构建纯文本或多模态用户消息
func buildUserMessage(text string, imageURLs []string) *schema.Message {
	if len(imageURLs) == 0 {
		return schema.UserMessage(text)
	}
	parts := make([]schema.MessageInputPart, 0, len(imageURLs)+1)
	if text != "" {
		parts = append(parts, schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeText,
			Text: text,
		})
	}
	for _, url := range imageURLs {
		parts = append(parts, schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeImageURL,
			Image: &schema.MessageInputImage{
				MessagePartCommon: schema.MessagePartCommon{URL: strPtr(url)},
			},
		})
	}
	return &schema.Message{Role: schema.User, UserInputMultiContent: parts}
}

func strPtr(s string) *string {
	return &s
}
