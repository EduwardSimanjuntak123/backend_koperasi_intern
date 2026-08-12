package routes

import (
	"backend_koperasi/internal/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterProductRoutes mendaftarkan seluruh endpoint Product.
func RegisterProductRoutes(
	router *gin.RouterGroup,
	productController *controllers.ProductController,
) {

	product := router.Group("/products")
	{
		product.GET("", productController.GetAll)
		product.GET("/:id", productController.GetByID)
		product.POST("", productController.Create)
		product.PUT("/:id", productController.Update)
		product.DELETE("/:id", productController.Delete)
	}
}
