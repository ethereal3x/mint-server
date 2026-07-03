package service

import (
	"context"

	"github.com/ethereal3x/apc/errs"
	authpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/auth"
	"github.com/ethereal3x/mint-server/internal/dto"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
)

// AuthServer AuthService gRPC 处理器
type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	logic AuthServiceLogic
}

// NewAuthServer 创建认证 gRPC 处理器
func NewAuthServer(authLogic AuthServiceLogic) *AuthServer {
	return &AuthServer{logic: authLogic}
}

// RegisterAccount 处理账号密码注册请求
func (s *AuthServer) RegisterAccount(ctx context.Context, req *authpb.RegisterAccountRequest) (*authpb.RegisterAccountResponse, error) {
	return errs.Handle(&authpb.RegisterAccountResponse{}, func(rsp *authpb.RegisterAccountResponse) error {
		if req.Account == "" || req.Password == "" {
			return mint_err.ErrParam
		}
		result, err := s.logic.RegisterAccount(ctx, &dto.RegisterAccountRequest{
			Account:     req.Account,
			Password:    req.Password,
			DisplayName: req.DisplayName,
			AvatarURL:   req.AvatarUrl,
		})
		if err != nil {
			return err
		}
		rsp.AccessToken = result.Token.AccessToken
		rsp.TokenType = result.Token.TokenType
		rsp.ExpiresIn = result.Token.ExpiresIn
		rsp.User = dto.UserToAuthProto(result.User)
		return nil
	})
}

// Login 处理登录请求
func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	return errs.Handle(&authpb.LoginResponse{}, func(rsp *authpb.LoginResponse) error {
		if req.Identifier == "" || req.Credential == "" {
			return mint_err.ErrParam
		}
		result, err := s.logic.Login(ctx, &dto.LoginRequest{
			Provider:   req.Provider,
			Identifier: req.Identifier,
			Credential: req.Credential,
		})
		if err != nil {
			return err
		}
		rsp.AccessToken = result.Token.AccessToken
		rsp.TokenType = result.Token.TokenType
		rsp.ExpiresIn = result.Token.ExpiresIn
		rsp.User = dto.UserToAuthProto(result.User)
		return nil
	})
}
