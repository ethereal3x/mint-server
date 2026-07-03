package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ethereal3x/apc/errs"
	"github.com/ethereal3x/apc/logger"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"go.uber.org/zap"
)

// HTTPResponse 统一 HTTP JSON 响应结构，与 gRPC proto 错误码对齐
type HTTPResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// writeHTTPJSON 写入 JSON 响应
func writeHTTPJSON(w http.ResponseWriter, status int, body HTTPResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.ContextError(context.Background(), "writeHTTPJSON", zap.Error(err))
	}
}

// writeHTTPSuccess 写入成功响应
func writeHTTPSuccess(w http.ResponseWriter, data any) {
	writeHTTPJSON(w, http.StatusOK, HTTPResponse{Code: 0, Message: "ok", Data: data})
}

// writeHTTPError 将 error 映射为 HTTP 响应，优先识别 BizError
func writeHTTPError(w http.ResponseWriter, status int, err error) {
	if bizErr, ok := errs.AsBizError(err); ok {
		writeHTTPJSON(w, status, HTTPResponse{Code: int32(bizErr.Code), Message: bizErr.Msg})
		return
	}
	writeHTTPJSON(w, status, HTTPResponse{Code: int32(mint_err.ERR_CODE_INTERNAL), Message: err.Error()})
}
