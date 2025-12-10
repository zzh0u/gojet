package v1api

import (
	"gojet/models"
	"gojet/service"
	"gojet/util/response"

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

// DeleteUser 根据 ID 删除用户
func DeleteUser(c *gin.Context) {
	var idParam IDParam
	if err := c.ShouldBindUri(&idParam); err != nil {
		response.BadRequest(c, response.MsgInvalidUserID)
		return
	}

	if err := service.DeleteUser(uint(idParam.ID)); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "删除成功", nil)
}

// GetUserByID 根据 ID 获取用户信息 - 返回单个用户详情
func GetUserByID(c *gin.Context) {
	var idParam IDParam
	if err := c.ShouldBindUri(&idParam); err != nil {
		response.BadRequest(c, response.MsgInvalidUserID)
		return
	}

	user, err := service.GetUserByID(uint(idParam.ID))
	if err != nil {
		// 使用 HandleError 统一处理，支持 400/404/500 等错误码
		response.HandleError(c, err)
		return
	}
	response.Success(c, "", user)
}

// GetAllUsers 获取所有用户列表 - 返回用户数组
func GetAllUsers(c *gin.Context) {
	users, err := service.GetAllUsers()
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "", users)
}

// CreateUser 创建新用户 - 从请求体获取用户信息
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, response.MsgInvalidParams)
		return
	}

	hashedPassword, err := models.HashPassword(user.Password)
	if err != nil {
		response.Error(c, 500, "密码加密失败")
		return
	}
	user.Password = hashedPassword

	newUser, err := service.CreateUser(&user)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "创建成功", newUser)
}

// UpdateUserRequest 更新用户请求结构体
type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	var idParam IDParam
	if err := c.ShouldBindUri(&idParam); err != nil {
		response.BadRequest(c, response.MsgInvalidUserID)
		return
	}

	var updateReq UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		response.BadRequest(c, response.MsgInvalidParams)
		return
	}

	updatedUser, err := service.UpdateUser(uint(idParam.ID), updateReq.Name)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "更新成功", updatedUser)
}
