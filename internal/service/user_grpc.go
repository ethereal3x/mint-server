package service

import (
	"context"

	"github.com/ethereal3x/apc/errs"
	userpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/user"
	"github.com/ethereal3x/mint-server/internal/dto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserServer UserService gRPC 处理器
type UserServer struct {
	userpb.UnimplementedUserServiceServer
	logic UserServiceLogic
}

// NewUserServer 创建用户 gRPC 处理器
func NewUserServer(authLogic UserServiceLogic) *UserServer {
	return &UserServer{logic: authLogic}
}

// GetMe 获取当前登录用户信息
func (s *UserServer) GetMe(ctx context.Context, _ *emptypb.Empty) (*userpb.GetMeResponse, error) {
	return errs.Handle(&userpb.GetMeResponse{}, func(rsp *userpb.GetMeResponse) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		user, err := s.logic.GetMe(ctx, userID)
		if err != nil {
			return err
		}
		rsp.User = dto.UserToUserProto(user)
		return nil
	})
}

// UpdateAvatar 更新当前用户头像
func (s *UserServer) UpdateAvatar(ctx context.Context, req *userpb.UpdateAvatarRequest) (*userpb.UpdateAvatarResponse, error) {
	return errs.Handle(&userpb.UpdateAvatarResponse{}, func(rsp *userpb.UpdateAvatarResponse) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		user, err := s.logic.UpdateAvatar(ctx, userID, req.AvatarUrl)
		if err != nil {
			return err
		}
		rsp.User = dto.UserToUserProto(user)
		return nil
	})
}

// UpdatePassword 更新当前用户密码
func (s *UserServer) UpdatePassword(ctx context.Context, req *userpb.UpdatePasswordRequest) (*userpb.UpdatePasswordResponse, error) {
	return errs.Handle(&userpb.UpdatePasswordResponse{}, func(rsp *userpb.UpdatePasswordResponse) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		return s.logic.UpdatePassword(ctx, userID, req.OldPassword, req.NewPassword)
	})
}

// UpdateNickname 更新当前用户展示名称
func (s *UserServer) UpdateNickname(ctx context.Context, req *userpb.UpdateNicknameRequest) (*userpb.UpdateNicknameResponse, error) {
	return errs.Handle(&userpb.UpdateNicknameResponse{}, func(rsp *userpb.UpdateNicknameResponse) error {
		userID, err := requireUserID(ctx)
		if err != nil {
			return err
		}
		user, err := s.logic.UpdateNickname(ctx, userID, req.Nickname)
		if err != nil {
			return err
		}
		rsp.User = dto.UserToUserProto(user)
		return nil
	})
}
