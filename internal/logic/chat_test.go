package logic

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/ethereal3x/mint-server/internal/dto"
	"github.com/ethereal3x/mint-server/internal/model"
)

type stubChatRecordRepo struct {
	created *model.DialogueRecord
}

// Create 记录创建调用参数
func (repo *stubChatRecordRepo) Create(_ context.Context, record *model.DialogueRecord) error {
	repo.created = record
	return nil
}

// ListByDialogue 桩实现
func (repo *stubChatRecordRepo) ListByDialogue(_ context.Context, _ *model.DialogueQuery) ([]*model.DialogueRecord, error) {
	return nil, nil
}

// ListDialogues 桩实现
func (repo *stubChatRecordRepo) ListDialogues(_ context.Context, _ string) ([]*model.DialogueSummary, error) {
	return nil, nil
}

// AggregateByModelForUser 桩实现
func (repo *stubChatRecordRepo) AggregateByModelForUser(_ context.Context, _ string) ([]*model.ModelStat, error) {
	return nil, nil
}

// TestSaveRecordPricing 校验 SaveRecord 按 token 单价计算费用
func TestSaveRecordPricing(t *testing.T) {
	repo := &stubChatRecordRepo{}
	chat := NewChat(nil, repo, nil)

	err := chat.SaveRecord(context.Background(), &dto.SaveRecordRequest{
		ChatRequest: &dto.ChatRequest{
			DialogueID: "dialogue-1",
			RecordID:   "record-1",
			UserID:     "user-1",
			Model:      "gpt-4o",
			Question:   "hello",
		},
		Answer: "world",
		Config: &model.ChatModelConfig{
			InputPrice:  2_000_000,
			OutputPrice: 4_000_000,
		},
		Usage: &schema.TokenUsage{
			PromptTokens:     1000,
			CompletionTokens: 500,
			TotalTokens:      1500,
		},
	})
	if err != nil {
		t.Fatalf("SaveRecord: %v", err)
	}
	if repo.created == nil {
		t.Fatal("expected record to be created")
	}
	if repo.created.InputCost != 2000.0 {
		t.Fatalf("input cost = %v, want 2000.0", repo.created.InputCost)
	}
	if repo.created.OutputCost != 2000.0 {
		t.Fatalf("output cost = %v, want 2000.0", repo.created.OutputCost)
	}
}
