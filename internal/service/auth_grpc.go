package service

import (
	"context"
	"strconv"

	"github.com/ethereal3x/apc/errs"
	authpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/auth"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/model"
)

// AuthServer AuthService gRPC 处理器
type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	logic *logic.Auth
}

// NewAuthServer 创建认证 gRPC 处理器
func NewAuthServer(authLogic *logic.Auth) *AuthServer {
	return &AuthServer{logic: authLogic}
}

// RegisterAccount 处理账号密码注册请求
func (s *AuthServer) RegisterAccount(ctx context.Context, req *authpb.RegisterAccountRequest) (*authpb.RegisterAccountResponse, error) {
	rsp := &authpb.RegisterAccountResponse{}
	if req.Account == "" || req.Password == "" {
		return errs.GenProtoReply(rsp, mint_err.ErrParam)
	}
	result, err := s.logic.RegisterAccount(ctx, &logic.RegisterAccountRequest{
		Account:     req.Account,
		Password:    req.Password,
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarUrl,
	})
	if err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	rsp.AccessToken = result.Token.AccessToken
	rsp.TokenType = result.Token.TokenType
	rsp.ExpiresIn = result.Token.ExpiresIn
	rsp.User = userToProto(result.User)
	return rsp, nil
}

// Login 处理登录请求
func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	rsp := &authpb.LoginResponse{}
	if req.Identifier == "" || req.Credential == "" {
		return errs.GenProtoReply(rsp, mint_err.ErrParam)
	}
	result, err := s.logic.Login(ctx, &logic.LoginRequest{
		Provider:   req.Provider,
		Identifier: req.Identifier,
		Credential: req.Credential,
	})
	if err != nil {
		return errs.GenProtoReply(rsp, err)
	}
	rsp.AccessToken = result.Token.AccessToken
	rsp.TokenType = result.Token.TokenType
	rsp.ExpiresIn = result.Token.ExpiresIn
	rsp.User = userToProto(result.User)
	return rsp, nil
}

// userToProto 转换用户模型为 proto
func userToProto(user *model.BaseUser) *authpb.UserInfo {
	if user == nil {
		return nil
	}
	return &authpb.UserInfo{
		UserId:      strconv.FormatInt(user.UserID, 10),
		DisplayName: user.Nickname,
		AvatarUrl:   user.AvatarURL,
		Status:      1,
	}
}
