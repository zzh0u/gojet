package auth

import (
	"gojet/internal/infra/apperror"
	"gojet/internal/infra/httputil"

	"github.com/gin-gonic/gin"
)

// Handler 认证 HTTP 处理器。
type Handler struct {
	svc *Service
}

// NewHandler 创建认证处理器。
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes 注册认证路由。
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	auth := router.Group("/v1/auth")
	auth.POST("/login", h.Login)
}

// Login
// @Summary 	用户登录
// @Description 系统用户登录
// @Id 			Login
// @Tags 		auth
// @Param 		m 		body 		LoginReq true "账号密码信息"
// @Success		200		{object}	httputil.Response{data=LoginResp}	"登录后token信息"
// @Failure 	400 	{object} 	httputil.Response "请求参数无效"
// @Failure 	401 	{object} 	httputil.Response "认证失败"
// @Failure 	404 	{object} 	httputil.Response "用户不存在"
// @Failure 	500 	{object} 	httputil.Response "服务器内部错误"
// @Router /v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, apperror.InvalidParams)
		return
	}

	resp, err := h.svc.Login(req)
	if err != nil {
		httputil.HandleError(c, err)
		return
	}

	httputil.Success(c, "登录成功", resp)
}
