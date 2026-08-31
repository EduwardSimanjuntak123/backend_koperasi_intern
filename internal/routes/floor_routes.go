package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFloorRoutes(
	router *gin.RouterGroup,
	floorController *controllers.FloorController,
) {

	// Public
	router.GET("/floors", floorController.GetAll)
	router.GET("/floors/:id", floorController.GetByID)
	router.GET("/floors/buildings/:building_id", floorController.GetByBuildingID)

	// Admin Only
	admin := router.Group("/floors")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		admin.POST("", floorController.Create)
		admin.PUT("/:id", floorController.Update)
		admin.DELETE("/:id", floorController.Delete)
	}
}
