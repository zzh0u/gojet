package v1api

import (
	"net/http"
	"strconv"

	"gojet/api"
	"gojet/models"

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

// InsertInitialData 插入初始学生数据 - 用于初始化测试数据
func (api *UserAPI) InsertInitialData(c *gin.Context) {
	// 调用服务层创建初始数据
	if err := api.userService.CreateInitialData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "数据插入成功"})
}

// DeleteUser 根据 ID 删除用户 - 从 URL 参数获取用户 ID
func (api *UserAPI) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"})
		return
	}

	if err := api.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "学生删除成功"})
}

// GetUserByID 根据 ID 获取用户信息 - 返回单个用户详情
func (api *UserAPI) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"})
		return
	}

	user, err := api.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetAllUsers 获取所有用户列表 - 返回用户数组
func (api *UserAPI) GetAllUsers(c *gin.Context) {
	users, err := api.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser 创建新用户 - 从请求体获取用户信息
func (api *UserAPI) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := api.userService.CreateUser(user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}
