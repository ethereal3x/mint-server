package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/apc/tracing"
	"github.com/go-sql-driver/mysql"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserRepo 用户 GORM 实现
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建用户仓储
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// CreateBaseUser 创建账号密码用户
func (r *UserRepo) CreateBaseUser(ctx context.Context, user *model.BaseUser) error {
	ctx, span := tracing.Start(ctx, "repo.UserRepo.CreateBaseUser")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		tracing.RecordError(ctx, err)
		if isDuplicateEntry(err) {
			return mint_err.ErrUserExists
		}
		logger.ContextError(ctx, "UserRepo.CreateBaseUser", zap.String("username", user.Username), zap.Error(err))
		return fmt.Errorf("create base user: %w", err)
	}
	return nil
}

// FindBaseUserByUsername 按用户名查询基础用户
func (r *UserRepo) FindBaseUserByUsername(ctx context.Context, query *model.BaseUserQuery) (*model.BaseUser, error) {
	var user model.BaseUser
	if err := r.db.WithContext(ctx).Where("username = ?", query.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "UserRepo.FindBaseUserByUsername", zap.String("username", query.Username), zap.Error(err))
		return nil, fmt.Errorf("find base user by username: %w", err)
	}
	return &user, nil
}

// FindBaseUserByUserID 按用户ID查询基础用户
func (r *UserRepo) FindBaseUserByUserID(ctx context.Context, query *model.BaseUserQuery) (*model.BaseUser, error) {
	var user model.BaseUser
	if err := r.db.WithContext(ctx).Where("user_id = ?", query.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.ContextError(ctx, "UserRepo.FindBaseUserByUserID", zap.Int64("user_id", query.UserID), zap.Error(err))
		return nil, fmt.Errorf("find base user by user_id: %w", err)
	}
	return &user, nil
}

// UpdateBaseUser 按用户ID更新基础用户字段
func (r *UserRepo) UpdateBaseUser(ctx context.Context, userID int64, updates map[string]any) error {
	ctx, span := tracing.Start(ctx, "repo.UserRepo.UpdateBaseUser")
	defer span.End()

	if err := r.db.WithContext(ctx).Model(&model.BaseUser{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
		tracing.RecordError(ctx, err)
		logger.ContextError(ctx, "UserRepo.UpdateBaseUser", zap.Int64("user_id", userID), zap.Error(err))
		return fmt.Errorf("update base user: %w", err)
	}
	return nil
}

// isDuplicateEntry 判断数据库错误是否为唯一键冲突
func isDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
