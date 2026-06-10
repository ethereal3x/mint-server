package logic

import (
	"context"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
)

// GetHistory 获取当前用户的对话历史记录
func (s *Chat) GetHistory(ctx context.Context, query *model.DialogueQuery) ([]*model.DialogueRecord, error) {
	records, err := s.recordRepo.ListByDialogue(ctx, query)
	if err != nil {
		logger.ContextError(ctx, "Chat.GetHistory", zap.String("dialogue_id", query.DialogueID), zap.String("user_id", query.UserID), zap.Error(err))
		return nil, mint_err.ErrGetHistory
	}
	return records, nil
}

// ListDialogues 获取对话摘要列表
func (s *Chat) ListDialogues(ctx context.Context, userID string) ([]*model.DialogueSummary, error) {
	ctx, span := tracing.Start(ctx, "logic.Chat.ListDialogues")
	defer span.End()

	summaries, err := s.recordRepo.ListDialogues(ctx, userID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Chat.ListDialogues", zap.String("user_id", userID), zap.Error(err))
		return nil, mint_err.ErrListDialogues
	}
	return summaries, nil
}

// AggregateByModel 按用户和模型聚合 token 消耗和费用
func (s *Chat) AggregateByModel(ctx context.Context, userID string) ([]*model.ModelStat, error) {
	stats, err := s.recordRepo.AggregateByModelForUser(ctx, userID)
	if err != nil {
		logger.ContextError(ctx, "Chat.AggregateByModel", zap.String("user_id", userID), zap.Error(err))
		return nil, mint_err.ErrModelStats
	}
	return stats, nil
}
