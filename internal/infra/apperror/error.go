package apperror

import "fmt"

// Error 是应用层统一错误类型，包含业务码和用户可读信息
type Error struct {
	Code    int    // 业务错误码（按需定义，例如 400/404/500 等）
	Message string // 返回给客户端的友好消息
	Err     error  // 原始错误（可为 nil）
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 使 errors.Is / As 能够访问底层错误
func (e *Error) Unwrap() error { return e.Err }

// New 创建一个新的 AppError
func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

// Wrap 包装底层 error 为 AppError（保留原始错误）
func Wrap(err error, code int, message string) *Error {
	return &Error{Code: code, Message: message, Err: err}
}
