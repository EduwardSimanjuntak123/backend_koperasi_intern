package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCourierRoutes(
	router *gin.RouterGroup,
	courierController *controllers.CourierController,
) {

	// Public
	router.GET("/couriers", courierController.GetAll)
	router.GET("/couriers/:id", courierController.GetByID)

	// Admin Only
	admin := router.Group("/couriers")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		admin.POST("", courierController.Create)
		admin.PUT("/:id", courierController.Update)
		admin.DELETE("/:id", courierController.Delete)
	}
}
