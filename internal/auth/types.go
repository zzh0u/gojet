package auth

// LoginReq 登录请求参数。
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResp 登录响应数据。
type LoginResp struct {
	Userid       int     `json:"userid"`        // 用户 ID
	Username     string  `json:"username"`      // 用户名
	NickName     string  `json:"nick_name"`     // 用户别名
	AccessToken  string  `json:"access_token"`  // access token
	RefreshToken string  `json:"refresh_token"` // refresh token
	ExpiresIn    float64 `json:"expires_in"`    // access token 过期时间（秒）
	TokenType    string  `json:"token_type"`    // token 类型
}
