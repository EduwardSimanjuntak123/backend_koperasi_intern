package routes

import (
	"backend_koperasi/internal/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterCategoryRoutes mendaftarkan seluruh endpoint Category.
func RegisterCategoryRoutes(
	router *gin.RouterGroup,
	categoryController *controllers.CategoryProductController,
) {

	category := router.Group("/categories")
	{
		category.GET("", categoryController.GetAll)
		category.GET("/:id", categoryController.GetByID)
		category.POST("", categoryController.Create)
		category.PUT("/:id", categoryController.Update)
		category.DELETE("/:id", categoryController.Delete)
	}
}
