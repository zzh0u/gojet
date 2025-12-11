package router

import (
	"gojet/api/v1api"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有应用路由
func SetupRoutes(r *gin.Engine) {
	apiV1 := r.Group("/v1")
	{
		health := apiV1.Group("/health")
		{
			health.GET("", v1api.HealthCheck)
		}

		users := apiV1.Group("/user")
		{
			users.POST("/insert", v1api.InsertInitialData)
			users.POST("", v1api.CreateUser)
			users.GET("/:id", v1api.GetUserByID)
			users.GET("", v1api.GetAllUsers)
			users.PUT("/:id", v1api.UpdateUser)
			users.DELETE("/:id", v1api.DeleteUser)
		}
		auth := apiV1.Group("")
		{
			auth.POST("/login", v1api.Login)
			auth.POST("/register", v1api.Register)
		}
	}
}
