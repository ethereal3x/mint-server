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

// StrategyRepo 策略规则 GORM 实现
type StrategyRepo struct {
	db *gorm.DB
}

// NewStrategyRepo 创建策略规则仓储
func NewStrategyRepo(db *gorm.DB) *StrategyRepo {
	return &StrategyRepo{db: db}
}

func (r *StrategyRepo) FindByID(ctx context.Context, id int32) (*model.StrategyRule, error) {
	var rule model.StrategyRule
	if err := r.db.WithContext(ctx).First(&rule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "StrategyRepo.FindByID", zap.Int32("id", id), zap.Error(err))
		return nil, fmt.Errorf("find strategy by id: %w", err)
	}
	return &rule, nil
}

func (r *StrategyRepo) FindByRuleID(ctx context.Context, ruleID string) (*model.StrategyRule, error) {
	var rule model.StrategyRule
	if err := r.db.WithContext(ctx).Where("rule_id = ?", ruleID).First(&rule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "StrategyRepo.FindByRuleID", zap.String("rule_id", ruleID), zap.Error(err))
		return nil, fmt.Errorf("find strategy by rule_id: %w", err)
	}
	return &rule, nil
}

func (r *StrategyRepo) FindByManufacturer(ctx context.Context, manufacturer string) (*model.StrategyRule, error) {
	var rule model.StrategyRule
	if err := r.db.WithContext(ctx).
		Where("agent_manufacturer = ? AND is_enabled = 1", manufacturer).
		First(&rule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "StrategyRepo.FindByManufacturer", zap.String("manufacturer", manufacturer), zap.Error(err))
		return nil, fmt.Errorf("find strategy by manufacturer: %w", err)
	}
	return &rule, nil
}

func (r *StrategyRepo) List(ctx context.Context, offset, limit int) ([]*model.StrategyRule, int64, error) {
	ctx, span := tracing.Start(ctx, "repo.StrategyRepo.List")
	defer span.End()

	var list []*model.StrategyRule
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.StrategyRule{}).Count(&total).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "StrategyRepo.List", zap.Error(err))
		return nil, 0, fmt.Errorf("count strategies: %w", err)
	}
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "StrategyRepo.List", zap.Int("offset", offset), zap.Int("limit", limit), zap.Error(err))
		return nil, 0, fmt.Errorf("list strategies: %w", err)
	}
	return list, total, nil
}

func (r *StrategyRepo) Create(ctx context.Context, rule *model.StrategyRule) error {
	ctx, span := tracing.Start(ctx, "repo.StrategyRepo.Create")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(rule).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "StrategyRepo.Create", zap.String("rule_id", rule.RuleID), zap.Error(err))
		return err
	}
	return nil
}

func (r *StrategyRepo) Update(ctx context.Context, rule *model.StrategyRule) error {
	ctx, span := tracing.Start(ctx, "repo.StrategyRepo.Update")
	defer span.End()

	if err := r.db.WithContext(ctx).Save(rule).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "StrategyRepo.Update", zap.String("rule_id", rule.RuleID), zap.Error(err))
		return err
	}
	return nil
}

func (r *StrategyRepo) Delete(ctx context.Context, ruleID string) error {
	ctx, span := tracing.Start(ctx, "repo.StrategyRepo.Delete")
	defer span.End()

	if err := r.db.WithContext(ctx).Where("rule_id = ?", ruleID).Delete(&model.StrategyRule{}).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "StrategyRepo.Delete", zap.String("rule_id", ruleID), zap.Error(err))
		return err
	}
	return nil
}
