package response

// 常用错误消息常量
const (
	// 通用错误
	MsgInvalidParams   = "请求参数无效"
	MsgInternalError   = "服务器内部错误"
	MsgDatabaseError   = "数据库操作失败"
	MsgRecordNotFound  = "记录不存在"
	MsgOperationFailed = "操作失败"

	// 用户相关错误
	MsgUserNotFound     = "用户不存在"
	MsgUserCreateFailed = "用户创建失败"
	MsgUserUpdateFailed = "用户更新失败"
	MsgUserDeleteFailed = "用户删除失败"
	MsgInvalidUserID    = "无效的用户 ID"

	// 数据库相关错误
	MsgDBConnectionError = "数据库连接失败"
	MsgDBQueryError      = "数据查询失败"
	MsgDBInsertError     = "数据插入失败"
	MsgDBUpdateError     = "数据更新失败"
	MsgDBDeleteError     = "数据删除失败"

	// 认证相关错误
	MsgAuthFailed   = "认证失败"
	MsgUnauthorized = "未授权访问"
	MsgTokenMissing = "令牌缺失"
	MsgTokenExpired = "令牌已过期"
	MsgTokenInvalid = "无效的令牌"
)
