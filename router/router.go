package router

import (
	"gojet/api/v1api"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(r *gin.Engine, v1api *v1api.UserAPI) {
	// API v1 routes
	apiV1 := r.Group("/api/v1")
	{
		// User routes
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
