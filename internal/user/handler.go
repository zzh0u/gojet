package user

import (
	"gojet/internal/infra/apperror"
	"gojet/internal/infra/httputil"

	"github.com/gin-gonic/gin"
)

// Handler 用户 HTTP 处理器。
type Handler struct {
	svc *Service
}

// NewHandler 创建用户处理器。
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes 注册用户路由。
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	users := router.Group("/v1/user")
	users.POST("/insert", h.InsertInitialData)
	users.POST("", h.CreateUser)
	users.GET("/:id", h.GetUserByID)
	users.GET("", h.GetAllUsers)
	users.PUT("/:id", h.UpdateUser)
	users.DELETE("/:id", h.DeleteUser)
}

// IDParam 用于绑定路径参数中的 ID。
type IDParam struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// InsertInitialData 插入初始用户数据。
func (h *Handler) InsertInitialData(c *gin.Context) {
	if err := h.svc.CreateInitialData(); err != nil {
		httputil.HandleError(c, err)
		return
	}
	httputil.Success(c, "数据插入成功", nil)
}

// DeleteUser
// @Summary 	根据 ID 删除用户
// @Description 根据 ID 删除系统用户
// @Id 			DeleteUser
// @Tags 		auth
// @Param 		id 		path 		int true "用户ID"
// @Success		200		{object}	httputil.Response{data=nil}	"删除成功"
// @Failure 	400 	{object} 	httputil.Response "请求参数无效"
// @Failure 	401 	{object} 	httputil.Response "认证失败"
// @Failure 	404 	{object} 	httputil.Response "用户不存在"
// @Failure 	500 	{object} 	httputil.Response "服务器内部错误"
// @Router 		/v1/user/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	var idParam IDParam
	if err := c.ShouldBindUri(&idParam); err != nil {
		httputil.BadRequest(c, apperror.InvalidUserID)
		return
	}

	if err := h.svc.DeleteUser(idParam.ID); err != nil {
		httputil.HandleError(c, err)
		return
	}
	httputil.Success(c, "删除成功", nil)
}

// GetUserByID
// @Summary 	根据 ID 获取用户信息
// @Description 根据 ID 获取系统用户详情
// @Id 			GetUserByID
// @Tags 		auth
// @Param 		id 		path 		int true "用户ID"
// @Success		200		{object}	httputil.Response{data=UserResponse}	"用户详情"
// @Failure 	400 	{object} 	httputil.Response "请求参数无效"
// @Failure 	401 	{object} 	httputil.Response "认证失败"
// @Failure 	404 	{object} 	httputil.Response "用户不存在"
// @Failure 	500 	{object} 	httputil.Response "服务器内部错误"
// @Router 		/v1/user/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	var idParam IDParam
	if err := c.ShouldBindUri(&idParam); err != nil {
		httputil.BadRequest(c, apperror.InvalidUserID)
		return
	}

	user, err := h.svc.GetUserByID(idParam.ID)
	if err != nil {
		httputil.HandleError(c, err)
		return
	}
	httputil.Success(c, "", user)
}

// GetAllUsers
// @Summary 	获取所有用户列表
// @Description 获取系统中所有用户的详细信息
// @Id 			GetAllUsers
// @Tags 		auth
// @Success		200		{object}	httputil.Response{data=[]UserResponse}	"用户列表"
// @Failure 	401 	{object} 	httputil.Response "认证失败"
// @Failure 	500 	{object} 	httputil.Response "服务器内部错误"
// @Router 		/v1/users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.svc.GetAllUsers()
	if err != nil {
		httputil.HandleError(c, err)
		return
	}
	httputil.Success(c, "", users)
}

// CreateUser
// @Summary 	创建新用户
// @Description 创建一个新的系统用户，从请求体获取用户信息
// @Id 			CreateUser
// @Tags 		auth
// @Param 		user 	body 		CreateUserRequest true "用户信息"
// @Success		200		{object}	httputil.Response{data=UserResponse}	"创建成功"
// @Failure 	400 	{object} 	httputil.Response "请求参数无效"
// @Failure 	401 	{object} 	httputil.Response "认证失败"
// @Failure 	500 	{object} 	httputil.Response "服务器内部错误"
// @Router 		/v1/user [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, apperror.InvalidParams)
		return
	}

	newUser, err := h.svc.CreateUser(req)
	if err != nil {
		httputil.HandleError(c, err)
		return
	}
	httputil.Success(c, "创建成功", newUser)
}

// UpdateUser
// @Summary 	更新用户信息
// @Description 根据 ID 更新系统用户的姓名
// @Id 			UpdateUser
// @Tags 		auth
// @Param 		id 		path 		int true "用户ID"
// @Param 		user 	body 		UpdateUserRequest true "更新用户信息"
// @Success		200		{object}	httputil.Response{data=UserResponse}	"更新成功"
// @Failure 	400 	{object} 	httputil.Response "请求参数无效"
// @Failure 	401 	{object} 	httputil.Response "认证失败"
// @Failure 	404 	{object} 	httputil.Response "用户不存在"
// @Failure 	500 	{object} 	httputil.Response "服务器内部错误"
// @Router 		/v1/user/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	var idParam IDParam
	if err := c.ShouldBindUri(&idParam); err != nil {
		httputil.BadRequest(c, apperror.InvalidUserID)
		return
	}

	var updateReq UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		httputil.BadRequest(c, apperror.InvalidParams)
		return
	}

	updatedUser, err := h.svc.UpdateUser(idParam.ID, updateReq)
	if err != nil {
		httputil.HandleError(c, err)
		return
	}
	httputil.Success(c, "更新成功", updatedUser)
}
