package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DialogueRepo 对话记录 GORM 实现
type DialogueRepo struct {
	db *gorm.DB
}

// NewDialogueRepo 创建对话记录仓储
func NewDialogueRepo(db *gorm.DB) *DialogueRepo {
	return &DialogueRepo{db: db}
}

// FindByRecordID 按记录ID查询对话记录
func (r *DialogueRepo) FindByRecordID(ctx context.Context, recordID string) (*model.DialogueRecord, error) {
	var record model.DialogueRecord
	if err := r.db.WithContext(ctx).First(&record, "record_id = ?", recordID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "DialogueRepo.FindByRecordID", zap.String("record_id", recordID), zap.Error(err))
		return nil, fmt.Errorf("find dialogue by record_id: %w", err)
	}
	return &record, nil
}

// ListByDialogue 按对话ID和用户ID查询对话记录
func (r *DialogueRepo) ListByDialogue(ctx context.Context, query *model.DialogueQuery) ([]*model.DialogueRecord, error) {
	var list []*model.DialogueRecord
	if err := r.db.WithContext(ctx).
		Where("dialogue_id = ? AND user_id = ?", query.DialogueID, query.UserID).
		Order("created_time ASC").
		Find(&list).Error; err != nil {
		logger.ContextError(ctx, "DialogueRepo.ListByDialogue", zap.String("dialogue_id", query.DialogueID), zap.String("user_id", query.UserID), zap.Error(err))
		return nil, fmt.Errorf("list dialogue by dialogue_id and user_id: %w", err)
	}
	return list, nil
}

// ListByUserID 按用户ID分页查询对话记录
func (r *DialogueRepo) ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.DialogueRecord, int64, error) {
	ctx, span := tracing.Start(ctx, "repo.DialogueRepo.ListByUserID")
	defer span.End()

	var list []*model.DialogueRecord
	var total int64
	query := r.db.WithContext(ctx).Model(&model.DialogueRecord{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "DialogueRepo.ListByUserID", zap.String("user_id", userID), zap.Error(err))
		return nil, 0, fmt.Errorf("count dialogues: %w", err)
	}
	if err := query.Offset(offset).Limit(limit).Order("created_time DESC").Find(&list).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "DialogueRepo.ListByUserID", zap.String("user_id", userID), zap.Int("offset", offset), zap.Int("limit", limit), zap.Error(err))
		return nil, 0, fmt.Errorf("list dialogues: %w", err)
	}
	return list, total, nil
}

// Create 创建对话记录
func (r *DialogueRepo) Create(ctx context.Context, record *model.DialogueRecord) error {
	ctx, span := tracing.Start(ctx, "repo.DialogueRepo.Create")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "DialogueRepo.Create", zap.String("record_id", record.RecordID), zap.Error(err))
		return err
	}
	return nil
}

// ListDialogues 获取对话摘要列表，按最后更新时间倒序
func (r *DialogueRepo) ListDialogues(ctx context.Context, userID string) ([]*model.DialogueSummary, error) {
	ctx, span := tracing.Start(ctx, "repo.DialogueRepo.ListDialogues")
	defer span.End()

	if userID == "" {
		return []*model.DialogueSummary{}, nil
	}

	var summaries []*model.DialogueSummary
	sql := `SELECT
		dialogue_id,
		(SELECT user_content FROM tb_user_agent_dialogues d2
		 WHERE d2.dialogue_id = d1.dialogue_id AND d2.user_id = ? ORDER BY created_time ASC LIMIT 1) AS title,
		COUNT(*) AS message_count,
		MAX(created_time) AS updated_time
	FROM tb_user_agent_dialogues d1`

	args := []interface{}{userID, userID}
	sql += " WHERE d1.user_id = ? GROUP BY dialogue_id ORDER BY updated_time DESC"

	if err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&summaries).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "DialogueRepo.ListDialogues", zap.String("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("list dialogues: %w", err)
	}
	return summaries, nil
}

// AggregateByModelForUser 按用户和模型聚合 token 消耗和费用
func (r *DialogueRepo) AggregateByModelForUser(ctx context.Context, userID string) ([]*model.ModelStat, error) {
	ctx, span := tracing.Start(ctx, "repo.DialogueRepo.AggregateByModelForUser")
	defer span.End()

	var stats []*model.ModelStat
	sql := `SELECT
		model,
		SUM(user_tokens) AS total_input_tokens,
		SUM(agent_tokens) AS total_output_tokens,
		SUM(input_cost) AS total_input_cost,
		SUM(output_cost) AS total_output_cost
	FROM tb_user_agent_dialogues
	WHERE user_id = ?
	GROUP BY model
	ORDER BY total_input_tokens + total_output_tokens DESC`

	if err := r.db.WithContext(ctx).Raw(sql, userID).Scan(&stats).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "DialogueRepo.AggregateByModelForUser", zap.String("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("aggregate by model: %w", err)
	}
	return stats, nil
}
