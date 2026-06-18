package v1api

import (
	"gojet/service"
	"gojet/utils/apperror"
	"gojet/utils/response"

	"github.com/gin-gonic/gin"
)

// IDParam 用于绑定路径参数中的ID
type IDParam struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// InsertInitialData 插入初始学生数据
func InsertInitialData(c *gin.Context) {
	// 调用服务层创建初始数据
	if err := service.CreateInitialData(); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "数据插入成功", nil)
}

// DeleteUser
// @Summary 	根据 ID 删除用户
// @Description 根据 ID 删除系统用户
// @Id 			DeleteUser
// @Tags 		auth
// @Param 		id 		path 		int true "用户ID"
// @Success		200		{object}	response.Response{data=nil}	"删除成功"
// @Failure 	400 	{object} 	response.Response "请求参数无效"
// @Failure 	401 	{object} 	response.Response "认证失败"
// @Failure 	404 	{object} 	response.Response "用户不存在"
// @Failure 	500 	{object} 	response.Response "服务器内部错误"
// @Router 		/v1/user/{id} [delete]
func DeleteUser(c *gin.Context) {
	var idParam IDParam
	if err := c.ShouldBindUri(&idParam); err != nil {
		response.BadRequest(c, apperror.InvalidUserID)
		return
	}

	if err := service.DeleteUser(idParam.ID); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "删除成功", nil)
}

// GetUserByID
// @Summary 	根据 ID 获取用户信息
// @Description 根据 ID 获取系统用户详情
// @Id 			GetUserByID
// @Tags 		auth
// @Param 		id 		path 		int true "用户ID"
// @Success		200		{object}	response.Response{data=service.UserResponse}	"用户详情"
// @Failure 	400 	{object} 	response.Response "请求参数无效"
// @Failure 	401 	{object} 	response.Response "认证失败"
// @Failure 	404 	{object} 	response.Response "用户不存在"
// @Failure 	500 	{object} 	response.Response "服务器内部错误"
// @Router 		/v1/user/{id} [get]
func GetUserByID(c *gin.Context) {
	var idParam IDParam
	if err := c.ShouldBindUri(&idParam); err != nil {
		response.BadRequest(c, apperror.InvalidUserID)
		return
	}

	user, err := service.GetUserByID(idParam.ID)
	if err != nil {
		// 使用 HandleError 统一处理，支持 400/404/500 等错误码
		response.HandleError(c, err)
		return
	}
	response.Success(c, "", user)
}

// GetAllUsers
// @Summary 	获取所有用户列表
// @Description 获取系统中所有用户的详细信息
// @Id 			GetAllUsers
// @Tags 		auth
// @Success		200		{object}	response.Response{data=[]service.UserResponse}	"用户列表"
// @Failure 	401 	{object} 	response.Response "认证失败"
// @Failure 	500 	{object} 	response.Response "服务器内部错误"
// @Router 		/v1/users [get]
func GetAllUsers(c *gin.Context) {
	users, err := service.GetAllUsers()
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "", users)
}

// CreateUser
// @Summary 	创建新用户
// @Description 创建一个新的系统用户，从请求体获取用户信息
// @Id 			CreateUser
// @Tags 		auth
// @Param 		user 	body 		service.CreateUserRequest true "用户信息"
// @Success		200		{object}	response.Response{data=service.UserResponse}	"创建成功"
// @Failure 	400 	{object} 	response.Response "请求参数无效"
// @Failure 	401 	{object} 	response.Response "认证失败"
// @Failure 	500 	{object} 	response.Response "服务器内部错误"
// @Router 		/v1/user [post]
func CreateUser(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperror.InvalidParams)
		return
	}

	newUser, err := service.CreateUser(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "创建成功", newUser)
}

// UpdateUser
// @Summary 	更新用户信息
// @Description 根据 ID 更新系统用户的姓名
// @Id 			UpdateUser
// @Tags 		auth
// @Param 		id 		path 		int true "用户ID"
// @Param 		user 	body 		service.UpdateUserRequest true "更新用户信息"
// @Success		200		{object}	response.Response{data=service.UserResponse}	"更新成功"
// @Failure 	400 	{object} 	response.Response "请求参数无效"
// @Failure 	401 	{object} 	response.Response "认证失败"
// @Failure 	404 	{object} 	response.Response "用户不存在"
// @Failure 	500 	{object} 	response.Response "服务器内部错误"
// @Router 		/v1/user/{id} [put]
func UpdateUser(c *gin.Context) {
	var idParam IDParam
	if err := c.ShouldBindUri(&idParam); err != nil {
		response.BadRequest(c, apperror.InvalidUserID)
		return
	}

	var updateReq service.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		response.BadRequest(c, apperror.InvalidParams)
		return
	}

	updatedUser, err := service.UpdateUser(idParam.ID, updateReq)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "更新成功", updatedUser)
}
