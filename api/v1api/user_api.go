package v1api

import (
	"net/http"
	"strconv"

	"gojet/api"
	"gojet/models"

	"github.com/gin-gonic/gin"
)

// UserAPI handles HTTP requests for user operations
type UserAPI struct {
	userService api.UserService
}

// NewUserAPI creates a new UserAPI
func NewUserAPI(userService api.UserService) *UserAPI {
	return &UserAPI{userService: userService}
}

// InsertInitialData inserts initial student data
func (api *UserAPI) InsertInitialData(c *gin.Context) {
	if err := api.userService.CreateInitialData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "数据插入成功"})
}

// DeleteUser deletes a user by ID
func (api *UserAPI) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := api.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "学生删除成功"})
}

// GetUserByID gets a user by ID
func (api *UserAPI) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := api.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetAllUsers gets all users
func (api *UserAPI) GetAllUsers(c *gin.Context) {
	users, err := api.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser creates a new user
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
