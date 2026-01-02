package v1api

import (
	"gojet/models"
	"gojet/service"
	"gojet/util/apperror"
	"gojet/util/response"

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

// Register 用户注册
// @Summary 	用户注册
// @Description 注册新用户
// @Id 			Register
// @Tags 		auth
// @Param 		user 	body 		models.User true "用户信息"
// @Success		200		{object}	response.Response{data=models.User}	"注册成功的用户信息"
// @Failure 	400 	{object} 	response.Response "请求参数无效"
// @Failure 	500 	{object} 	response.Response "服务器内部错误"
// @Router /v1/register [post]
func Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		response.BadRequest(ctx, apperror.InvalidParams)
		return
	}

	// 对密码进行哈希处理
	if user.Password != "" {
		hashedPassword, err := models.HashPassword(user.Password)
		if err != nil {
			response.Error(ctx, 500, "密码加密失败")
			return
		}
		user.Password = hashedPassword
	}

	// 创建用户
	newUser, err := service.CreateUser(&user)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, "注册成功", newUser)
}
