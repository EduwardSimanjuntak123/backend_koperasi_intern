package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(
	router *gin.RouterGroup,
	categoryController *controllers.CategoryProductController,
) {

	// Public
	router.GET("/categories", categoryController.GetAll)
	router.GET("/categories/:id", categoryController.GetByID)

	// Admin Only
	admin := router.Group("/categories")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		admin.POST("", categoryController.Create)
		admin.PUT("/:id", categoryController.Update)
		admin.DELETE("/:id", categoryController.Delete)
	}
}
