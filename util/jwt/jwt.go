package jwt

import (
	"gojet/util/response"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SkipRouter 路由请求跳过的path 最后一个/匹配即可
var SkipRouter = map[string]bool{}

func Token(c *gin.Context) {
	path := strings.Split(c.Request.URL.Path, "/")

	lastPath := path[len(path)-1]
	if SkipRouter[lastPath] {
		c.Next()
		return
	}
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		response.Error(c, 403, response.MsgTokenMissing)
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
		// Make sure the `alg` is what we except.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	}
}

func parseToken(tokenString string, secret string, c *gin.Context) {
	// Parse the token.
	token, err := jwt.Parse(tokenString, secretFunc(secret))

	// Parse error.
	if err != nil {
		response.Error(c, 403, response.MsgTokenInvalid)
		c.Abort()
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["id"].(float64))
		username := claims["username"].(string)
		c.Set("userid", userID)
		c.Set("username", username)
		c.Set("token", tokenString)
		c.Next()
	} else {
		// token 过期了
		response.Error(c, 403, response.MsgTokenExpired)
		c.Abort()
	}
}

type Context struct {
	ID       int
	Username string
}

// Sign 生成一个JWT token并返回token字符串
// 根据提供的上下文、用户信息、密钥和持续时间创建签名的JWT token
func Sign(c Context, secret string, duration time.Duration) (tokenString string, err error) {
	// 创建包含用户信息和时间戳的JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       c.ID,
		"username": c.Username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(duration).Unix(),
	})
	// 使用指定的密钥对token进行签名
	tokenString, err = token.SignedString([]byte(secret))

	return
}
