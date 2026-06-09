package logic

import (
	"context"
	"fmt"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"github.com/ethereal3x/mint-server/internal/crypto"
	"github.com/ethereal3x/mint-server/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// StrategyRepo 策略规则数据访问接口
type StrategyRepo interface {
	FindByID(ctx context.Context, id int32) (*model.StrategyRule, error)
	FindByRuleID(ctx context.Context, ruleID string) (*model.StrategyRule, error)
	FindByManufacturer(ctx context.Context, manufacturer string) (*model.StrategyRule, error)
	List(ctx context.Context, offset, limit int) ([]*model.StrategyRule, int64, error)
	Create(ctx context.Context, rule *model.StrategyRule) error
	Update(ctx context.Context, rule *model.StrategyRule) error
	Delete(ctx context.Context, ruleID string) error
}

// Strategy 策略规则业务逻辑
type Strategy struct {
	repo      StrategyRepo
	secretKey []byte
}

// NewStrategy 创建策略规则业务逻辑
func NewStrategy(repo StrategyRepo, secretKey []byte) *Strategy {
	return &Strategy{repo: repo, secretKey: secretKey}
}

// GetByManufacturer 按厂商获取启用的策略，API Key 已解密（供内部 chat 使用）
func (s *Strategy) GetByManufacturer(ctx context.Context, manufacturer string) (*model.StrategyRule, error) {
	ctx, span := tracing.Start(ctx, "logic.Strategy.GetByManufacturer")
	defer span.End()

	rule, err := s.repo.FindByManufacturer(ctx, manufacturer)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Strategy.GetByManufacturer", zap.String("manufacturer", manufacturer), zap.Error(err))
		return nil, err
	}
	return decryptRule(rule, s.secretKey), nil
}

// GetByRuleID 按 RuleID 获取策略（API Key 保持加密，供管理 API 使用）
func (s *Strategy) GetByRuleID(ctx context.Context, ruleID string) (*model.StrategyRule, error) {
	ctx, span := tracing.Start(ctx, "logic.Strategy.GetByRuleID")
	defer span.End()

	rule, err := s.repo.FindByRuleID(ctx, ruleID)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Strategy.GetByRuleID", zap.String("rule_id", ruleID), zap.Error(err))
		return nil, err
	}
	return rule, nil
}

// List 分页获取策略列表（API Key 保持加密，供管理 API 使用）
func (s *Strategy) List(ctx context.Context, offset, limit int) ([]*model.StrategyRule, int64, error) {
	ctx, span := tracing.Start(ctx, "logic.Strategy.List")
	defer span.End()

	list, total, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Strategy.List", zap.Int("offset", offset), zap.Int("limit", limit), zap.Error(err))
		return nil, 0, err
	}
	return list, total, nil
}

// Create 创建策略（自动生成 RuleID、加密 API Key）
func (s *Strategy) Create(ctx context.Context, rule *model.StrategyRule) error {
	ctx, span := tracing.Start(ctx, "logic.Strategy.Create")
	defer span.End()

	rule.RuleID = uuid.NewString()
	encryptedKey, err := crypto.EncryptAPIKey(rule.APIKey, s.secretKey)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Strategy.Create", zap.Error(err))
		return fmt.Errorf("encrypt api key: %w", err)
	}
	rule.APIKey = encryptedKey
	if err := s.repo.Create(ctx, rule); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Strategy.Create", zap.String("rule_id", rule.RuleID), zap.Error(err))
		return err
	}
	return nil
}

// Update 更新策略
func (s *Strategy) Update(ctx context.Context, rule *model.StrategyRule) error {
	ctx, span := tracing.Start(ctx, "logic.Strategy.Update")
	defer span.End()

	if _, err := s.repo.FindByRuleID(ctx, rule.RuleID); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Strategy.Update", zap.String("rule_id", rule.RuleID), zap.Error(err))
		return fmt.Errorf("strategy not found: %s", rule.RuleID)
	}
	encryptedKey, err := crypto.EncryptAPIKey(rule.APIKey, s.secretKey)
	if err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Strategy.Update", zap.Error(err))
		return fmt.Errorf("encrypt api key: %w", err)
	}
	rule.APIKey = encryptedKey
	if err := s.repo.Update(ctx, rule); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Strategy.Update", zap.String("rule_id", rule.RuleID), zap.Error(err))
		return err
	}
	return nil
}

// Delete 删除策略
func (s *Strategy) Delete(ctx context.Context, ruleID string) error {
	ctx, span := tracing.Start(ctx, "logic.Strategy.Delete")
	defer span.End()

	if err := s.repo.Delete(ctx, ruleID); err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "Strategy.Delete", zap.String("rule_id", ruleID), zap.Error(err))
		return err
	}
	return nil
}

// decryptRule 返回 API Key 已解密的新副本，不修改原始数据
func decryptRule(rule *model.StrategyRule, secretKey []byte) *model.StrategyRule {
	if rule == nil {
		return nil
	}
	cp := *rule
	decrypted, err := crypto.DecryptAPIKey(cp.APIKey, secretKey)
	if err != nil {
		return rule
	}
	cp.APIKey = decrypted
	return &cp
}
