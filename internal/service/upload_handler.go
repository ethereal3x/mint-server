package service

import (
	"net/http"
	"strconv"

	"github.com/ethereal3x/apc/logger"
	"github.com/ethereal3x/mint-server/internal/dto"
	mint_err "github.com/ethereal3x/mint-server/internal/errs"
	"github.com/ethereal3x/mint-server/internal/middleware"
	"github.com/ethereal3x/mint-server/internal/model"
	"github.com/ethereal3x/mint-server/internal/util"
	"go.uber.org/zap"
)

const httpMaxUploadSize = 50 << 20 // 50MB

// UploadHandler 文件上传 HTTP 处理器
type UploadHandler struct {
	logic        UploadServiceLogic
	tokenManager *util.TokenManager
}

// NewUploadHandler 创建上传 HTTP 处理器
func NewUploadHandler(uploadLogic UploadServiceLogic, tokenManager *util.TokenManager) *UploadHandler {
	return &UploadHandler{logic: uploadLogic, tokenManager: tokenManager}
}

// HandleUpload 处理文件上传请求
func (handler *UploadHandler) HandleUpload(w http.ResponseWriter, request *http.Request) {
	userID, ok := handler.requireUserID(w, request)
	if !ok {
		return
	}

	if err := request.ParseMultipartForm(httpMaxUploadSize); err != nil {
		writeHTTPError(w, http.StatusBadRequest, mint_err.ErrParam)
		return
	}

	file, header, err := request.FormFile("file")
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest, mint_err.ErrParam)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	result, err := handler.logic.Upload(request.Context(), &dto.UploadRequest{
		UserID:      userID,
		FileName:    header.Filename,
		Reader:      file,
		Size:        header.Size,
		ContentType: contentType,
	})
	if err != nil {
		logger.ContextError(request.Context(), "UploadHandler.HandleUpload", zap.String("filename", header.Filename), zap.Error(err))
		writeHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	writeHTTPSuccess(w, result)
}

// HandleGetFile 处理文件详情查询请求
func (handler *UploadHandler) HandleGetFile(w http.ResponseWriter, request *http.Request) {
	userID, ok := handler.requireUserID(w, request)
	if !ok {
		return
	}
	idStr := request.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest, mint_err.ErrParam)
		return
	}

	record, err := handler.logic.GetUpload(request.Context(), &model.FileUploadQuery{ID: int32(id), UserID: userID})
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	if record == nil {
		writeHTTPJSON(w, http.StatusNotFound, HTTPResponse{Code: -1, Message: "文件不存在"})
		return
	}

	writeHTTPSuccess(w, record)
}

// HandleListFiles 处理文件列表查询请求
func (handler *UploadHandler) HandleListFiles(w http.ResponseWriter, request *http.Request) {
	userID, ok := handler.requireUserID(w, request)
	if !ok {
		return
	}

	list, err := handler.logic.ListUploads(request.Context(), userID)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	writeHTTPSuccess(w, list)
}

// requireUserID 从 HTTP 请求中解析当前用户 ID
func (handler *UploadHandler) requireUserID(w http.ResponseWriter, request *http.Request) (string, bool) {
	userID, err := middleware.RequireUserIDFromRequest(request, handler.tokenManager)
	if err != nil {
		writeHTTPError(w, http.StatusUnauthorized, err)
		return "", false
	}
	return userID, true
}
