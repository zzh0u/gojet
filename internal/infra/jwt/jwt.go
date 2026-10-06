package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType token 类型常量
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// Context 是生成 token 时携带的用户信息。
type Context struct {
	ID       int
	Username string
}

// Claims 是项目使用的 JWT 声明结构。
type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

// SignWithType 生成指定类型的 JWT token。
func SignWithType(c Context, secret string, duration time.Duration, tokenType string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID:       c.ID,
		Username: c.Username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	})

	return token.SignedString([]byte(secret))
}

// SignAccess 生成 access token。
func SignAccess(c Context, secret string, duration time.Duration) (string, error) {
	return SignWithType(c, secret, duration, TokenTypeAccess)
}

// SignRefresh 生成 refresh token。
func SignRefresh(c Context, secret string, duration time.Duration) (string, error) {
	return SignWithType(c, secret, duration, TokenTypeRefresh)
}

// ParseAccessToken 解析 access token。
func ParseAccessToken(tokenString string, secret string) (*Claims, error) {
	return ParseWithType(tokenString, secret, TokenTypeAccess)
}

// ParseRefreshToken 解析 refresh token。
func ParseRefreshToken(tokenString string, secret string) (*Claims, error) {
	return ParseWithType(tokenString, secret, TokenTypeRefresh)
}

// ParseWithType 解析 token 并校验 token 类型。
func ParseWithType(tokenString string, secret string, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, secretFunc(secret))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if claims.Type != expectedType {
		return nil, errors.New("unexpected token type")
	}

	return claims, nil
}

func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	}
}
