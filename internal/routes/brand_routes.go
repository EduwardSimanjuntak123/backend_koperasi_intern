package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterBrandRoutes(
	router *gin.RouterGroup,
	brandController *controllers.BrandController,
) {

	// Public
	router.GET("/brands", brandController.GetAll)
	router.GET("/brands/:id", brandController.GetByID)

	// Admin Only
	admin := router.Group("/brands")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		admin.POST("", brandController.Create)
		admin.PUT("/:id", brandController.Update)
		admin.DELETE("/:id", brandController.Delete)
	}
}
