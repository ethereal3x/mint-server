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

func (r *DialogueRepo) ListByDialogueID(ctx context.Context, dialogueID string) ([]*model.DialogueRecord, error) {
	var list []*model.DialogueRecord
	if err := r.db.WithContext(ctx).
		Where("dialogue_id = ?", dialogueID).
		Order("created_time ASC").
		Find(&list).Error; err != nil {
		logger.ContextError(ctx, "DialogueRepo.ListByDialogueID", zap.String("dialogue_id", dialogueID), zap.Error(err))
		return nil, fmt.Errorf("list dialogue by dialogue_id: %w", err)
	}
	return list, nil
}

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

	query := r.db.WithContext(ctx).Model(&model.DialogueRecord{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var summaries []*model.DialogueSummary
	sql := `SELECT
		dialogue_id,
		(SELECT user_content FROM tb_user_agent_dialogues d2
		 WHERE d2.dialogue_id = d1.dialogue_id ORDER BY created_time ASC LIMIT 1) AS title,
		COUNT(*) AS message_count,
		MAX(created_time) AS updated_time
	FROM tb_user_agent_dialogues d1`

	args := make([]interface{}, 0)
	if userID != "" {
		sql += " WHERE d1.user_id = ?"
		args = append(args, userID)
	}
	sql += " GROUP BY dialogue_id ORDER BY updated_time DESC"

	if err := query.Select("").Raw(sql, args...).Scan(&summaries).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "DialogueRepo.ListDialogues", zap.String("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("list dialogues: %w", err)
	}
	return summaries, nil
}
