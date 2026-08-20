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
	cart := api.Group("/cart")

	cart.Use(
		middleware.AuthMiddleware(),
		middleware.BuyerMiddleware(),
	)

	{
		cart.GET("", cartController.GetCart)
		cart.POST("/items", cartController.AddToCart)
		cart.PUT("/items/:item_id", cartController.UpdateCartItemQuantity)
		cart.DELETE("/items/:item_id", cartController.RemoveFromCart)
	}
}
