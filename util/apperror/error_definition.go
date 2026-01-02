package apperror

const (
	// 通用错误
	InvalidParams   = "请求参数无效"
	InternalError   = "服务器内部错误"
	DatabaseError   = "数据库操作失败"
	RecordNotFound  = "记录不存在"
	OperationFailed = "操作失败"

	// 用户相关错误
	UserNotFound     = "用户不存在"
	UserCreateFailed = "用户创建失败"
	UserUpdateFailed = "用户更新失败"
	UserDeleteFailed = "用户删除失败"
	InvalidUserID    = "无效的用户 ID"

	// 数据库相关错误
	DBQueryError  = "数据查询失败"
	DBInsertError = "数据插入失败"
	DBUpdateError = "数据更新失败"
	DBDeleteError = "数据删除失败"

	// 认证相关错误
	AuthFailed   = "认证失败"
	Unauthorized = "未授权访问"
	TokenMissing = "令牌缺失"
	TokenExpired = "令牌已过期"
	TokenInvalid = "无效的令牌"
)
