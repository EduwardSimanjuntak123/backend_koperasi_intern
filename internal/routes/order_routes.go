package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(
	api *gin.RouterGroup,
	orderController *controllers.OrderController,
) {

	// ======================================
	// Buyer Routes
	// ======================================

	buyer := api.Group("/orders")

	buyer.Use(
		middleware.AuthMiddleware(),
		middleware.BuyerMiddleware(),
	)

	{
		buyer.GET("", orderController.GetMyOrders)
		buyer.POST("", orderController.CreateOrder)
		buyer.GET("/:id", orderController.GetOrderByID)
	}

	// ======================================
	// Admin / Store Routes
	// ======================================

	admin := api.Group("/admin/orders")

	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)

	{
		admin.GET("", orderController.GetAllOrders)
		admin.GET("/:id", orderController.GetOrderByID)
		admin.PUT("/:id", orderController.UpdateOrder)
		admin.DELETE("/:id", orderController.DeleteOrder)
	}
}
