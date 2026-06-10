package logic

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/ethereal3x/mint-server/internal/model"
)

// wrapLLMError 包装 LLM 调用错误
func wrapLLMError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("模型调用异常: %w", err)
}

// LLMAdapter 模型调用抽象接口，后续新增模型只需实现此接口
type LLMAdapter interface {
	Generate(ctx context.Context, config *model.ChatModelConfig, messages []*schema.Message) (string, *schema.TokenUsage, error)
	Stream(ctx context.Context, config *model.ChatModelConfig, messages []*schema.Message, contentChan chan<- string) (*schema.TokenUsage, error)
}

// EinoAdapter 基于 CloudWeGo Eino 的 LLM 适配器
type EinoAdapter struct{}

// NewEinoAdapter 创建 Eino 适配器
func NewEinoAdapter() *EinoAdapter {
	return &EinoAdapter{}
}

// createModel 根据模型配置创建 Eino ChatModel
func (a *EinoAdapter) createModel(ctx context.Context, config *model.ChatModelConfig) (*openai.ChatModel, error) {
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  config.APIKey,
		Model:   config.ModelType,
		BaseURL: config.URL,
	})
}

// Generate 非流式生成
func (a *EinoAdapter) Generate(ctx context.Context, config *model.ChatModelConfig, messages []*schema.Message) (string, *schema.TokenUsage, error) {
	cm, err := a.createModel(ctx, config)
	if err != nil {
		return "", nil, fmt.Errorf("create model: %w", err)
	}

	msg, err := cm.Generate(ctx, messages)
	if err != nil {
		return "", nil, wrapLLMError(err)
	}
	return msg.Content, msg.ResponseMeta.Usage, nil
}

// Stream 流式生成，增量内容通过 contentChan 返回
func (a *EinoAdapter) Stream(ctx context.Context, config *model.ChatModelConfig, messages []*schema.Message, contentChan chan<- string) (*schema.TokenUsage, error) {
	cm, err := a.createModel(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create model: %w", err)
	}

	reader, err := cm.Stream(ctx, messages)
	if err != nil {
		return nil, wrapLLMError(err)
	}
	defer reader.Close()

	var usage *schema.TokenUsage
	for {
		chunk, err := reader.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("stream recv: %w", err)
		}
		if chunk.ResponseMeta != nil && chunk.ResponseMeta.Usage != nil {
			usage = chunk.ResponseMeta.Usage
		}
		if chunk.Content != "" {
			select {
			case contentChan <- chunk.Content:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
	return usage, nil
}
