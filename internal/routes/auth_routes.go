package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	api *gin.RouterGroup,
	authController *controllers.AuthController,
) {

	// ==========================
	// Public Routes
	// ==========================
	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// ==========================
	// Protected Routes
	// ==========================
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.AuthMiddleware())
	{
		authProtected.GET("/me", authController.Me)
		authProtected.POST("/logout", authController.Logout)
	}
}
