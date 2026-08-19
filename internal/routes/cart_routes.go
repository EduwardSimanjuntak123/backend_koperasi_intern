package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(
	api *gin.RouterGroup,
	cartController *controllers.CartController,
) {
	// Membuat grup route "/cart"
	cartRoutes := api.Group("/cart")
	
	// Menerapkan AuthMiddleware agar hanya user yang sudah login yang bisa mengakses
	cartRoutes.Use(middleware.AuthMiddleware())
	{
		cartRoutes.GET("", cartController.GetCart)
		cartRoutes.POST("/items", cartController.AddToCart)
		cartRoutes.PUT("/items/:item_id", cartController.UpdateCartItemQuantity)
		cartRoutes.DELETE("/items/:item_id", cartController.RemoveFromCart)
	}
}