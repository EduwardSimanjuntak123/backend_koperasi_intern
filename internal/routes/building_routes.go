package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterBuildingRoutes(
	router *gin.RouterGroup,
	buildingController *controllers.BuildingController,
) {

	// Public
	router.GET("/buildings", buildingController.GetAll)
	router.GET("/buildings/:id", buildingController.GetByID)

	// Admin Only
	admin := router.Group("/buildings")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		admin.POST("", buildingController.Create)
		admin.PUT("/:id", buildingController.Update)
		admin.DELETE("/:id", buildingController.Delete)
	}
}
