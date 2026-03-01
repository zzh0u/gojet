package jwt

import (
	"gojet/util/apperror"
	"gojet/util/response"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SkipRouter 路由请求跳过的path 最后一个/匹配即可
var SkipRouter = map[string]bool{}

// TokenType token 类型常量
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

func Token(c *gin.Context) {
	path := strings.Split(c.Request.URL.Path, "/")

	lastPath := path[len(path)-1]
	if SkipRouter[lastPath] {
		c.Next()
		return
	}
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		response.Error(c, 403, apperror.TokenMissing)
		c.Abort()
		return
	}
	// Load the jwt secret from the gin config
	js, _ := c.Get("jwt-secret")
	secret := js.(string)

	// Parse the header to get the token part.
	t := strings.Replace(header, "Bearer ", "", 1)
	parseToken(t, secret, c)
}

// secretFunc validates the secret format.
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	}
}

func parseToken(tokenString string, secret string, c *gin.Context) {
	if !parseTokenWithType(tokenString, secret, c, TokenTypeAccess) {
		response.Error(c, 403, apperror.TokenInvalid)
		c.Abort()
	} else {
		c.Next()
	}
}

// parseTokenWithType 内部解析函数，添加 token 类型验证
func parseTokenWithType(tokenString string, secret string, c *gin.Context, expectedType string) bool {
	token, err := jwt.Parse(tokenString, secretFunc(secret))
	if err != nil {
		return false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 处理 token 类型验证
		switch expectedType {
		case TokenTypeAccess:
			tokenType, ok := claims["type"].(string)
			if !ok || tokenType != TokenTypeAccess {
				return false
			}
		case TokenTypeRefresh:
			tokenType, ok := claims["type"].(string)
			if !ok || tokenType != TokenTypeRefresh {
				return false
			}
		default:
			return false
		}

		userID := int(claims["id"].(float64))
		username := claims["username"].(string)
		c.Set("userid", userID)
		c.Set("username", username)
		c.Set("token", tokenString)
		return true
	}
	return false
}

type Context struct {
	ID       int
	Username string
}

// SignWithType 生成指定类型的 JWT token
func SignWithType(c Context, secret string, duration time.Duration, tokenType string) (tokenString string, err error) {
	// 创建包含用户信息、token 类型和时间戳的 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       c.ID,
		"username": c.Username,
		"type":     tokenType,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(duration).Unix(),
	})
	// 使用指定的密钥对 token 进行签名
	tokenString, err = token.SignedString([]byte(secret))
	return
}

// SignAccess 生成一个 JWT token 并返回 token 字符串
// 根据提供的上下文、用户信息、密钥和持续时间创建签名的 JWT token
func SignAccess(c Context, secret string, duration time.Duration) (tokenString string, err error) {
	return SignWithType(c, secret, duration, TokenTypeAccess)
}

// SignRefresh 生成 refresh token
func SignRefresh(c Context, secret string, duration time.Duration) (tokenString string, err error) {
	return SignWithType(c, secret, duration, TokenTypeRefresh)
}
