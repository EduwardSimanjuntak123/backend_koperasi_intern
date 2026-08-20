package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUnitRoutes(
	router *gin.RouterGroup,
	UnitController *controllers.UnitController,
) {

	// Public
	router.GET("/units", UnitController.GetAll)
	router.GET("/units/:id", UnitController.GetByID)

	// Admin Only
	admin := router.Group("/units")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		admin.POST("", UnitController.Create)
		admin.PUT("/:id", UnitController.Update)
		admin.DELETE("/:id", UnitController.Delete)
	}
}
