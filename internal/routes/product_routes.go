package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterProductRoutes mendaftarkan seluruh endpoint Product.
func RegisterProductRoutes(
	router *gin.RouterGroup,
	productController *controllers.ProductController,
) {
	router.GET("/products", productController.GetAll)
	router.GET("/products/:id", productController.GetByID)
	admin := router.Group("/products")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		admin.POST("", productController.Create)
		admin.PUT("/:id", productController.Update)
		admin.DELETE("/:id", productController.Delete)
	}
}
