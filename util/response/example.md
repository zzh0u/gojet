# 响应工具包使用示例

## 基本响应格式

所有 API 响应现在都使用统一的响应格式：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": ...
}
```

## 使用示例

### 1. 返回单个对象

```go
// 获取用户信息
user, err := service.GetUserByID(id)
if err != nil {
    response.NotFound(c, "用户不存在")
    return
}
response.Success(c, "", user)
```

### 2. 返回数组数据

```go
// 获取所有用户
users, err := service.GetAllUsers()
if err != nil {
    response.InternalServerError(c, "查询失败")
    return
}
response.Success(c, "", users)
```

### 3. 创建资源

```go
// 创建用户
newUser, err := service.CreateUser(name)
if err != nil {
    response.InternalServerError(c, "创建失败")
    return
}
response.Success(c, "创建成功", newUser)
```

### 4. 更新资源

```go
// 更新用户
updatedUser, err := service.UpdateUser(id, name)
if err != nil {
    response.InternalServerError(c, "更新失败")
    return
}
response.Success(c, "更新成功", updatedUser)
```

### 5. 删除资源

```go
// 删除用户
err := service.DeleteUser(id)
if err != nil {
    response.InternalServerError(c, "删除失败")
    return
}
response.Success(c, "删除成功", nil)
```

### 6. 错误响应

```go
// 参数错误
response.BadRequest(c, "参数无效")

// 未授权
response.Unauthorized(c, "请先登录")

// 无权限
response.Forbidden(c, "权限不足")

// 资源不存在
response.NotFound(c, "用户不存在")

// 服务器错误
response.InternalServerError(c, "数据库连接失败")

// 自定义错误码
response.Error(c, 503, "服务不可用")
```

## 错误消息常量

使用预定义的错误消息常量，保持错误信息一致性：

```go
// 通用错误
response.MsgInvalidParams     // "请求参数无效"
response.MsgInternalError     // "服务器内部错误"
response.MsgDatabaseError     // "数据库操作失败"
response.MsgRecordNotFound    // "记录不存在"

// 用户相关错误
response.MsgUserNotFound      // "用户不存在"
response.MsgUserCreateFailed  // "用户创建失败"
response.MsgUserUpdateFailed  // "用户更新失败"
response.MsgUserDeleteFailed  // "用户删除失败"
response.MsgInvalidUserID     // "无效的用户 ID"
```

## 最佳实践

1. **统一使用响应工具函数**：不要再直接使用 `c.JSON()`
2. **使用预定义错误消息**：保持错误信息一致性
3. **选择合适的 HTTP 状态码**：
   - 200 OK: 查询、更新成功
   - 201 Created: 创建成功
   - 400 Bad Request: 参数错误
   - 401 Unauthorized: 未授权
   - 403 Forbidden: 无权限
   - 404 Not Found: 资源不存在
   - 500 Internal Server Error: 服务器错误
4. **删除操作使用 Success**：返回标准删除成功响应
