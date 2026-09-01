package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(
	router *gin.RouterGroup,
	paymentController *controllers.PaymentController,
) {
	// ==========================
	// Authenticated (Buyer / Admin)
	// ==========================
	authenticated := router.Group("/payments")
	authenticated.Use(
		middleware.AuthMiddleware(),
	)
	{
		authenticated.POST("", paymentController.Create)
		authenticated.POST("/pay", paymentController.Pay)
		authenticated.GET("/:id", paymentController.GetByID)
		authenticated.GET("/order/:order_id", paymentController.GetByOrderID)
	}

	// ==========================
	// Admin
	// ==========================
	admin := router.Group("/payments")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		admin.GET("", paymentController.GetAll)
		admin.PUT("/:id", paymentController.Update)
		admin.PATCH("/:id", paymentController.Update)
		admin.DELETE("/:id", paymentController.Delete)
	}
}
