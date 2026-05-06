package middleware

import (
	"strings"

	"gojet/utils/apperror"
	"gojet/utils/jwt"
	"gojet/utils/response"

	"github.com/gin-gonic/gin"
)

// JWT 验证访问令牌，并将用户信息写入请求上下文。
func JWT(secret string, skipPaths ...string) gin.HandlerFunc {
	skipLastPath := make(map[string]struct{}, len(skipPaths))
	for _, path := range skipPaths {
		skipLastPath[path] = struct{}{}
	}

	return func(c *gin.Context) {
		lastPath := lastSegment(c.Request.URL.Path)
		if _, ok := skipLastPath[lastPath]; ok {
			c.Next()
			return
		}

		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if header == "" {
			response.Error(c, 403, apperror.TokenMissing)
			c.Abort()
			return
		}

		tokenString, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || strings.TrimSpace(tokenString) == "" {
			response.Error(c, 403, apperror.TokenInvalid)
			c.Abort()
			return
		}

		claims, err := jwt.ParseAccessToken(strings.TrimSpace(tokenString), secret)
		if err != nil {
			response.Error(c, 403, apperror.TokenInvalid)
			c.Abort()
			return
		}

		c.Set("userid", claims.ID)
		c.Set("username", claims.Username)
		c.Set("token", strings.TrimSpace(tokenString))
		c.Next()
	}
}

func lastSegment(path string) string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return ""
	}

	parts := strings.Split(trimmed, "/")
	return parts[len(parts)-1]
}
