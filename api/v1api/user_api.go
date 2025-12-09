package v1api

import (
	"strconv"

	"gojet/api"
	"gojet/models"
	"gojet/util/response"

	"github.com/gin-gonic/gin"
)

// UserAPI 用户 API 处理器
type UserAPI struct {
	userService api.User
}

// NewUserAPI 创建用户 API 实例
func NewUserAPI(userService api.User) *UserAPI {
	return &UserAPI{userService: userService}
}

// InsertInitialData 插入初始学生数据
func (api *UserAPI) InsertInitialData(c *gin.Context) {
	// 调用服务层创建初始数据
	if err := api.userService.CreateInitialData(); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "数据插入成功", nil)
}

// DeleteUser 根据 ID 删除用户
func (api *UserAPI) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, response.MsgInvalidUserID)
		return
	}

	if err := api.userService.DeleteUser(uint(id)); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "删除成功", nil)
}

// GetUserByID 根据 ID 获取用户信息 - 返回单个用户详情
func (api *UserAPI) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, response.MsgInvalidUserID)
		return
	}

	user, err := api.userService.GetUserByID(uint(id))
	if err != nil {
		// 使用 HandleError 统一处理，支持 400/404/500 等错误码
		response.HandleError(c, err)
		return
	}
	response.Success(c, "", user)
}

// GetAllUsers 获取所有用户列表 - 返回用户数组
func (api *UserAPI) GetAllUsers(c *gin.Context) {
	users, err := api.userService.GetAllUsers()
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "", users)
}

// CreateUser 创建新用户 - 从请求体获取用户信息
func (api *UserAPI) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, response.MsgInvalidParams)
		return
	}

	// if err := user.Validate(); err != nil {
	// 	fieldErrors := models.FormatValidationError(err)
	// 	response.BadRequestWithData(c, "用户信息验证失败", fieldErrors)
	// 	return
	// }

	newUser, err := api.userService.CreateUser(user.Name)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "创建成功", newUser)
}

// UpdateUser 更新用户信息
func (api *UserAPI) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, response.MsgInvalidUserID)
		return
	}

	var updateReq struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		response.BadRequest(c, response.MsgInvalidParams)
		return
	}

	// // 验证 Name 字段
	// tempUser := &models.User{Name: updateReq.Name}
	// if err := tempUser.Validate(); err != nil {
	// 	fieldErrors := models.FormatValidationError(err)
	// 	response.BadRequestWithData(c, "用户信息验证失败", fieldErrors)
	// 	return
	// }

	updatedUser, err := api.userService.UpdateUser(uint(id), updateReq.Name)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "更新成功", updatedUser)
}
