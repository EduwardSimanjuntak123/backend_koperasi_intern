package routes

import (
	"backend_koperasi/internal/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterStoreRoutes mendaftarkan seluruh endpoint Store.
func RegisterStoreRoutes(
	router *gin.RouterGroup,
	storeController *controllers.StoreController,
) {

	store := router.Group("/stores")
	{
		store.GET("", storeController.GetAll)
		store.GET("/:id", storeController.GetByID)
		store.POST("", storeController.Create)
		store.PUT("/:id", storeController.Update)
		store.DELETE("/:id", storeController.Delete)
	}
}
