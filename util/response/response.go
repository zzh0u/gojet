package response

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"gojet/util/apperror"
)

// Response 统一响应结构体
type Response struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 消息
	Data    any    `json:"data"`    // 数据
}

// Success 返回成功响应
func Success(c *gin.Context, message string, data any) {
	if message == "" {
		message = "操作成功"
	}
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, code int, message string) {
	httpCode := http.StatusBadRequest
	switch code {
	case 400:
		httpCode = http.StatusBadRequest
	case 401:
		httpCode = http.StatusUnauthorized
	case 403:
		httpCode = http.StatusForbidden
	case 404:
		httpCode = http.StatusNotFound
	case 500:
		httpCode = http.StatusInternalServerError
	}

	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 返回400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// NotFound 返回404错误
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// InternalServerError 返回500错误
func InternalServerError(c *gin.Context, message string) {
	Error(c, 500, message)
}

// HandleError 统一处理 service 层返回的错误。
// - 如果是 *errpkg.Error，则按照其中的 Code/Message 返回对应响应。
// - 否则返回通用 500（服务器内部错误）。
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	var e *apperror.Error
	if errors.As(err, &e) {
		// 记录错误日志，包含原始错误信息（如果有）
		if e.Err != nil {
			slog.Error("应用错误", "code", e.Code, "message", e.Message, "original_error", e.Err)
		} else {
			slog.Error("应用错误", "code", e.Code, "message", e.Message)
		}

		switch e.Code {
		case 400:
			BadRequest(c, e.Message)
		case 404:
			NotFound(c, e.Message)
		case 500:
			InternalServerError(c, e.Message)
		default:
			InternalServerError(c, e.Message)
		}
		return
	}
	// 非 Error 类型，记录日志并返回通用内部错误
	slog.Error("未处理的应用错误", "error", err)
	InternalServerError(c, MsgInternalError)
}
