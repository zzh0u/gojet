package middleware

import (
	"gojet/cache"
	"gojet/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InjectDependencies 将启动时初始化的依赖注入到 Gin 上下文中。
func InjectDependencies(db *gorm.DB, redisClient *cache.RedisClient, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("jwt-secret", cfg.JWT.Secret)

		if sqlDB, err := db.DB(); err == nil {
			c.Set("db", sqlDB)
		}
		if redisClient != nil {
			c.Set("redis", redisClient)
		}

		c.Set("config", cfg)
		c.Next()
	}
}
