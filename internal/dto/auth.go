package dto

import (
	"strconv"

	authpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/auth"
	userpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/user"
	"github.com/ethereal3x/mint-server/internal/model"
	"github.com/ethereal3x/mint-server/internal/util"
)

// RegisterAccountRequest 账号密码注册参数
type RegisterAccountRequest struct {
	Account     string
	Password    string
	DisplayName string
	AvatarURL   string
}

// LoginRequest 登录参数
type LoginRequest struct {
	Provider   string
	Identifier string
	Credential string
}

// AuthResult 登录态结果
type AuthResult struct {
	User  *model.BaseUser
	Token *util.TokenResult
}

// UserToAuthProto 转换用户模型为 authpb.UserInfo
func UserToAuthProto(user *model.BaseUser) *authpb.UserInfo {
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

// UserToUserProto 转换用户模型为 userpb.UserInfo
func UserToUserProto(user *model.BaseUser) *userpb.UserInfo {
	if user == nil {
		return nil
	}
	return &userpb.UserInfo{
		UserId:      strconv.FormatInt(user.UserID, 10),
		DisplayName: user.Nickname,
		AvatarUrl:   user.AvatarURL,
		Status:      1,
	}
}
