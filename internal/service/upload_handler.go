package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereal3x/apc/logger"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/logic"
	"go.uber.org/zap"
)

const httpMaxUploadSize = 50 << 20 // 50MB

// UploadHandler 文件上传 HTTP 处理器
type UploadHandler struct {
	logic *logic.UploadLogic
}

// NewUploadHandler 创建上传 HTTP 处理器
func NewUploadHandler(logic *logic.UploadLogic) *UploadHandler {
	return &UploadHandler{logic: logic}
}

func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "admin"
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

	result, err := h.logic.Upload(r.Context(), userID, header.Filename, file, header.Size, contentType)
	if err != nil {
		logger.ContextError(r.Context(), "UploadHandler.HandleUpload", zap.String("filename", header.Filename), zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]any{"code": mint_err.ERR_CODE_INTERNAL, "message": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"code": 0, "message": "ok", "data": result})
}

func (h *UploadHandler) HandleGetFile(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"code": mint_err.ERR_CODE_PARAM, "message": "无效的文件 ID"})
		return
	}

	record, err := h.logic.GetUpload(r.Context(), int32(id))
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

func (h *UploadHandler) HandleListFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "admin"
	}

	list, err := h.logic.ListUploads(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"code": mint_err.ERR_CODE_DB_QUERY, "message": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"code": 0, "message": "ok", "data": list})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
