package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(
	router *gin.RouterGroup,
	orderController *controllers.OrderController,
) {
	// ===========================
	// BUYER
	// ===========================
	buyer := router.Group("/orders")
	buyer.Use(
		middleware.AuthMiddleware(),
		middleware.BuyerMiddleware(),
	)
	{
		buyer.GET("", orderController.GetMyOrders)
		buyer.GET("/:id", orderController.GetByID)
		buyer.POST("", orderController.Create)
		buyer.POST("/checkout-cart", orderController.CheckoutCart)
		buyer.POST("/:id/cancel", orderController.Cancel)
		buyer.DELETE("/:id", orderController.Cancel)
	}

	// ===========================
	// ADMIN / TOKO
	// ===========================
	admin := router.Group("/orders")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		admin.GET("/all", orderController.GetAll)
		admin.GET("/admin/by-status", orderController.GetByStatus)
		admin.GET("/admin/:id", orderController.GetAdminByID)
		admin.GET("/dashboard/summary", orderController.GetDashboardSummary)
		admin.GET("/dashboard/revenue/today", orderController.GetTodayRevenue)
		admin.PUT("/:id/status", orderController.UpdateStatus)
		admin.PATCH("/:id/status", orderController.UpdateStatus)
		admin.PUT("/:id/courier", orderController.AssignCourier)
		admin.PATCH("/:id/courier", orderController.AssignCourier)
		admin.DELETE("/admin/:id", orderController.Delete)
	}
}
