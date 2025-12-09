package router

import (
	"gojet/api/v1api"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有应用路由
func SetupRoutes(r *gin.Engine, v1api *v1api.UserAPI) {
	apiV1 := r.Group("/v1")
	{
		users := apiV1.Group("/users")
		{
			users.POST("/insert", v1api.InsertInitialData)
			users.DELETE("/:id", v1api.DeleteUser)
			users.GET("/:id", v1api.GetUserByID)
			users.GET("", v1api.GetAllUsers)
			users.POST("", v1api.CreateUser)
		}
	}
}
