package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterStoreMemberRoutes mendaftarkan seluruh endpoint StoreMember.
func RegisterStoreMemberRoutes(
	router *gin.RouterGroup,
	storeMemberController *controllers.StoreMemberController,
) {
	router.GET("/store-members", storeMemberController.GetAll)
	router.GET("/store-members/:id", storeMemberController.GetByID)
	admin := router.Group("/store-members")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{

		admin.POST("", storeMemberController.Create)
		admin.PUT("/:id", storeMemberController.Update)
		admin.DELETE("/:id", storeMemberController.Delete)
	}
}
