package api

import (
	"net/http"
	"otp-service/pkg/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, authHandler *AuthHandler, userHandler *UserHandler, jwtSecret string) {
	// health check of the service
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// api version
	v1 := router.Group("/api/v1")

	// Authentication routes
	auth := v1.Group("/auth")
	{
		auth.POST("/send-otp", authHandler.SendOTP)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
	}
	// User routes (authentication required)
	users := v1.Group("/users")
	users.Use(middleware.JWTAuth(jwtSecret))
	{
		users.GET("/:id", userHandler.GetUser)
		users.GET("", userHandler.GetUsers)
	}
}
