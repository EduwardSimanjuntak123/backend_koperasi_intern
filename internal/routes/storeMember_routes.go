package routes

import (
	"backend_koperasi/internal/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterStoreMemberRoutes mendaftarkan seluruh endpoint StoreMember.
func RegisterStoreMemberRoutes(
	router *gin.RouterGroup,
	storeMemberController *controllers.StoreMemberController,
) {

	storeMember := router.Group("/store-members")
	{
		storeMember.GET("", storeMemberController.GetAll)
		storeMember.GET("/:id", storeMemberController.GetByID)
		storeMember.POST("", storeMemberController.Create)
		storeMember.PUT("/:id", storeMemberController.Update)
		storeMember.DELETE("/:id", storeMemberController.Delete)
	}
}
