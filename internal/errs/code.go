package errs

import "github.com/ethereal3x/apc/errs"

// 业务错误码
const (
	ERR_CODE_PARAM              errs.ErrorCode = -1000
	ERR_CODE_MODEL_NOT_FOUND    errs.ErrorCode = -1001
	ERR_CODE_CONFIG_NOT_FOUND   errs.ErrorCode = -1002
	ERR_CODE_CHAT_STREAM        errs.ErrorCode = -1003
	ERR_CODE_CHAT_GENERATE      errs.ErrorCode = -1004
	ERR_CODE_SAVE_RECORD        errs.ErrorCode = -1005
	ERR_CODE_LIST_DIALOGUES     errs.ErrorCode = -1006
	ERR_CODE_GET_HISTORY        errs.ErrorCode = -1007
	ERR_CODE_LIST_MODELS        errs.ErrorCode = -1008
	ERR_CODE_LIST_CONFIGS       errs.ErrorCode = -1009
	ERR_CODE_CREATE_CONFIG      errs.ErrorCode = -1010
	ERR_CODE_UPDATE_CONFIG      errs.ErrorCode = -1011
	ERR_CODE_DELETE_CONFIG      errs.ErrorCode = -1012
	ERR_CODE_MODEL_STATS        errs.ErrorCode = -1013
	ERR_CODE_ENCRYPT_API_KEY    errs.ErrorCode = -1014
	ERR_CODE_DECRYPT_API_KEY    errs.ErrorCode = -1015
	ERR_CODE_DB_QUERY           errs.ErrorCode = -1016
	ERR_CODE_DB_CREATE          errs.ErrorCode = -1017
	ERR_CODE_DB_UPDATE          errs.ErrorCode = -1018
	ERR_CODE_DB_DELETE          errs.ErrorCode = -1019
	ERR_CODE_INTERNAL           errs.ErrorCode = -1020
	ERR_CODE_UNAUTHORIZED       errs.ErrorCode = -1021
	ERR_CODE_FORBIDDEN          errs.ErrorCode = -1022
	ERR_CODE_USER_EXISTS        errs.ErrorCode = -1023
	ERR_CODE_USER_NOT_FOUND     errs.ErrorCode = -1024
	ERR_CODE_INVALID_CREDENTIAL errs.ErrorCode = -1025
	ERR_CODE_TOKEN_INVALID      errs.ErrorCode = -1026
	ERR_CODE_TOKEN_EXPIRED      errs.ErrorCode = -1027
	ERR_CODE_PASSWORD_WEAK      errs.ErrorCode = -1028
)

// 预定义业务错误实例
var (
	ErrParam             = errs.New(ERR_CODE_PARAM, "参数错误")
	ErrModelNotFound     = errs.New(ERR_CODE_MODEL_NOT_FOUND, "模型配置不存在")
	ErrConfigNotFound    = errs.New(ERR_CODE_CONFIG_NOT_FOUND, "配置不存在")
	ErrChatStream        = errs.New(ERR_CODE_CHAT_STREAM, "流式对话失败")
	ErrChatGenerate      = errs.New(ERR_CODE_CHAT_GENERATE, "对话生成失败")
	ErrSaveRecord        = errs.New(ERR_CODE_SAVE_RECORD, "保存记录失败")
	ErrListDialogues     = errs.New(ERR_CODE_LIST_DIALOGUES, "获取对话列表失败")
	ErrGetHistory        = errs.New(ERR_CODE_GET_HISTORY, "获取历史记录失败")
	ErrListModels        = errs.New(ERR_CODE_LIST_MODELS, "获取模型列表失败")
	ErrListConfigs       = errs.New(ERR_CODE_LIST_CONFIGS, "获取配置列表失败")
	ErrCreateConfig      = errs.New(ERR_CODE_CREATE_CONFIG, "创建配置失败")
	ErrUpdateConfig      = errs.New(ERR_CODE_UPDATE_CONFIG, "更新配置失败")
	ErrDeleteConfig      = errs.New(ERR_CODE_DELETE_CONFIG, "删除配置失败")
	ErrModelStats        = errs.New(ERR_CODE_MODEL_STATS, "获取模型统计失败")
	ErrEncryptAPIKey     = errs.New(ERR_CODE_ENCRYPT_API_KEY, "加密API密钥失败")
	ErrDecryptAPIKey     = errs.New(ERR_CODE_DECRYPT_API_KEY, "解密API密钥失败")
	ErrDBQuery           = errs.New(ERR_CODE_DB_QUERY, "数据库查询失败")
	ErrDBCreate          = errs.New(ERR_CODE_DB_CREATE, "数据库创建失败")
	ErrDBUpdate          = errs.New(ERR_CODE_DB_UPDATE, "数据库更新失败")
	ErrDBDelete          = errs.New(ERR_CODE_DB_DELETE, "数据库删除失败")
	ErrInternal          = errs.New(ERR_CODE_INTERNAL, "内部服务错误")
	ErrUnauthorized      = errs.New(ERR_CODE_UNAUTHORIZED, "未登录或登录已失效")
	ErrForbidden         = errs.New(ERR_CODE_FORBIDDEN, "无权限访问")
	ErrUserExists        = errs.New(ERR_CODE_USER_EXISTS, "账号已存在")
	ErrUserNotFound      = errs.New(ERR_CODE_USER_NOT_FOUND, "用户不存在")
	ErrInvalidCredential = errs.New(ERR_CODE_INVALID_CREDENTIAL, "账号或密码错误")
	ErrTokenInvalid      = errs.New(ERR_CODE_TOKEN_INVALID, "登录令牌无效")
	ErrTokenExpired      = errs.New(ERR_CODE_TOKEN_EXPIRED, "登录已过期")
	ErrPasswordWeak      = errs.New(ERR_CODE_PASSWORD_WEAK, "密码强度不足")
)
