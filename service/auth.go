package service

import (
	"gojet/config"
	"gojet/util/apperror"
	"gojet/util/jwt"
	"gojet/util/response"
	"time"

	"github.com/gin-gonic/gin"
)

var cfg *config.Config

// InitAuth 初始化认证服务
func InitAuth(config *config.Config) {
	cfg = config
}

// LoginReq 登录请求参数
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResp 登录响应数据
type LoginResp struct {
	Userid      int     `json:"userid"`       // 用户ID
	Username    string  `json:"username"`     // 用户名称
	NickName    string  `json:"nick_name"`    // 用户别名
	AccessToken string  `json:"access_token"` // accessToken
	ExpiresIn   float64 `json:"expires_in"`   // 过期时间
	TokenType   string  `json:"token_type"`   // token类型
}

// Login 执行登录逻辑
func (req *LoginReq) Login(ctx *gin.Context) (*LoginResp, error) {
	user, err := userRepo.GetUserByUserName(req.Username)
	if err != nil {
		return nil, apperror.Wrap(err, 404, response.MsgUserNotFound)
	}

	// 验证密码
	if !user.CompareSimple(req.Password) {
		return nil, apperror.New(401, response.MsgAuthFailed)
	}

	// 设置token过期时间
	var duration = time.Duration(cfg.JWT.ExpireHours) * time.Hour

	// 生成JWT token
	secret, exists := ctx.Get("jwt-secret")
	if !exists {
		return nil, apperror.New(500, "JWT secret 未配置")
	}

	token, err := jwt.Sign(jwt.Context{ID: user.ID, Username: user.Username}, secret.(string), duration)
	if err != nil {
		return nil, apperror.Wrap(err, 500, "生成Token失败")
	}

	resp := &LoginResp{
		Userid:      user.ID,
		Username:    user.Username,
		NickName:    user.NickName,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   duration.Seconds(),
	}
	return resp, nil
}
