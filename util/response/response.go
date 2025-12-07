package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// Unauthorized 返回401错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

// Forbidden 返回403错误
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, message)
}

// NotFound 返回404错误
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// InternalServerError 返回500错误
func InternalServerError(c *gin.Context, message string) {
	Error(c, 500, message)
}
