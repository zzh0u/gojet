package v1api

import (
	"gojet/service"
	"gojet/utils/apperror"
	"gojet/utils/response"

	"github.com/gin-gonic/gin"
)

// Login
// @Summary 	用户登录
// @Description 系统用户登录
// @Id 			Login
// @Tags 		auth
// @Param 		m 		body 		service.LoginReq true "账号密码信息"
// @Success		200		{object}	response.Response{data=service.LoginResp}	"登录后token信息"
// @Failure 	400 	{object} 	response.Response "请求参数无效"
// @Failure 	401 	{object} 	response.Response "认证失败"
// @Failure 	404 	{object} 	response.Response "用户不存在"
// @Failure 	500 	{object} 	response.Response "服务器内部错误"
// @Router /v1/login [post]
func Login(ctx *gin.Context) {
	var req service.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, apperror.InvalidParams)
		return
	}

	resp, err := req.Login(ctx)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, "登录成功", resp)
}
