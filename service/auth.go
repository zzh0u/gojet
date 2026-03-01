package service

import (
	"gojet/config"
	"gojet/util/apperror"
	"gojet/util/jwt"
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
	Userid       int     `json:"userid"`        // 用户 ID
	Username     string  `json:"username"`      // 用户名
	NickName     string  `json:"nick_name"`     // 用户别名
	AccessToken  string  `json:"access_token"`  // access token
	RefreshToken string  `json:"refresh_token"` // refresh token
	ExpiresIn    float64 `json:"expires_in"`    // access token 过期时间（秒）
	TokenType    string  `json:"token_type"`    // token 类型
}

// Login 执行登录逻辑
func (req *LoginReq) Login(ctx *gin.Context) (*LoginResp, error) {
	user, err := userRepo.GetUserByUserName(req.Username)
	if err != nil {
		return nil, apperror.Wrap(err, 404, apperror.UserNotFound)
	}

	// 验证密码
	if !user.CompareSimple(req.Password) {
		return nil, apperror.New(401, apperror.AuthFailed)
	}

	// 从 Gin 上下文获取 JWT 密钥
	secret, exists := ctx.Get("jwt-secret")
	if !exists {
		return nil, apperror.New(500, "JWT secret 未配置")
	}
	jwtSecret := secret.(string)

	// 生成 access token
	accessDuration := time.Duration(cfg.JWT.AccessExpireHours) * time.Hour
	accessToken, err := jwt.SignAccess(jwt.Context{ID: user.ID, Username: user.Username}, jwtSecret, accessDuration)
	if err != nil {
		return nil, apperror.Wrap(err, 500, "生成 Access Token 失败")
	}

	// 生成 refresh token
	refreshDuration := time.Duration(cfg.JWT.RefreshExpireDays) * 24 * time.Hour
	refreshToken, err := jwt.SignRefresh(jwt.Context{ID: user.ID, Username: user.Username}, jwtSecret, refreshDuration)
	if err != nil {
		return nil, apperror.Wrap(err, 500, "生成 Refresh Token 失败")
	}

	resp := &LoginResp{
		Userid:       user.ID,
		Username:     user.Username,
		NickName:     user.NickName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    accessDuration.Seconds(),
	}
	return resp, nil
}
