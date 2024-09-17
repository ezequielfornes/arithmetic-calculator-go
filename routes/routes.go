package routes

import (
	"arithmetic-calculator/controllers"
	"arithmetic-calculator/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/api/v1/auth/login", controllers.Login)
	r.POST("/api/v1/auth/register", controllers.Register)
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(http.StatusNoContent)
	})
	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middlewares.JWTAuthMiddleware())
	{
		protected.POST("/operation", controllers.PerformOperation)

		protected.GET("/records", controllers.GetRecords)
		protected.GET("/user", controllers.GetUserInformation)
	}

	return r
}
