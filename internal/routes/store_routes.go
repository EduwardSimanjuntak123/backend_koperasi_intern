package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterStoreRoutes mendaftarkan seluruh endpoint Store.
func RegisterStoreRoutes(
	router *gin.RouterGroup,
	storeController *controllers.StoreController,
) {
	router.GET("/stores", storeController.GetAll)
	router.GET("/stores/:id", storeController.GetByID)
	admin := router.Group("/stores")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{

		admin.POST("", storeController.Create)
		admin.PUT("/:id", storeController.Update)
		admin.DELETE("/:id", storeController.Delete)
	}
}
