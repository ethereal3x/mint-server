package service

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/mint-server/internal/auth"
	"github.com/ethereal3x/mint-server/internal/dto"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/model"
	"go.uber.org/zap"
)

const httpMaxUploadSize = 50 << 20 // 50MB

// UploadHandler 文件上传 HTTP 处理器
type UploadHandler struct {
	logic        UploadServiceLogic
	tokenManager *auth.TokenManager
}

// NewUploadHandler 创建上传 HTTP 处理器
func NewUploadHandler(uploadLogic UploadServiceLogic, tokenManager *auth.TokenManager) *UploadHandler {
	return &UploadHandler{logic: uploadLogic, tokenManager: tokenManager}
}

// HandleUpload 处理文件上传请求
func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authUserID(w, r)
	if !ok {
		return
	}

	if err := r.ParseMultipartForm(httpMaxUploadSize); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"code": mint_err.ERR_CODE_PARAM, "message": "文件过大或格式错误"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"code": mint_err.ERR_CODE_PARAM, "message": "未找到上传文件"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	result, err := h.logic.Upload(r.Context(), &dto.UploadRequest{
		UserID:      userID,
		FileName:    header.Filename,
		Reader:      file,
		Size:        header.Size,
		ContentType: contentType,
	})
	if err != nil {
		logger.ContextError(r.Context(), "UploadHandler.HandleUpload", zap.String("filename", header.Filename), zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]any{"code": mint_err.ERR_CODE_INTERNAL, "message": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"code": 0, "message": "ok", "data": result})
}

// HandleGetFile 处理文件详情查询请求
func (h *UploadHandler) HandleGetFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authUserID(w, r)
	if !ok {
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"code": mint_err.ERR_CODE_PARAM, "message": "无效的文件 ID"})
		return
	}

	record, err := h.logic.GetUpload(r.Context(), &model.FileUploadQuery{ID: int32(id), UserID: userID})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"code": mint_err.ERR_CODE_DB_QUERY, "message": err.Error()})
		return
	}
	if record == nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"code": -1, "message": "文件不存在"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"code": 0, "message": "ok", "data": record})
}

// HandleListFiles 处理文件列表查询请求
func (h *UploadHandler) HandleListFiles(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authUserID(w, r)
	if !ok {
		return
	}

	list, err := h.logic.ListUploads(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"code": mint_err.ERR_CODE_DB_QUERY, "message": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"code": 0, "message": "ok", "data": list})
}

// authUserID 从 HTTP 请求中解析当前用户ID
func (h *UploadHandler) authUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	principal, err := auth.PrincipalFromRequest(r.Context(), h.tokenManager, r)
	if err != nil || principal.UserID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"code": mint_err.ERR_CODE_UNAUTHORIZED, "message": mint_err.ErrUnauthorized.Error()})
		return "", false
	}
	return principal.UserID, true
}

// writeJSON 写入 JSON 响应
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.ContextError(context.Background(), "writeJSON", zap.Error(err))
	}
}
