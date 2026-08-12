package routes

import (
	"backend_koperasi/internal/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoleRoutes mendaftarkan seluruh endpoint Role.
func RegisterRoleRoutes(
	router *gin.RouterGroup,
	roleController *controllers.RolesController,
) {

	role := router.Group("/roles")
	{
		role.GET("", roleController.GetAll)
		role.GET("/:id", roleController.GetByID)
		role.POST("", roleController.Create)
		role.PUT("/:id", roleController.Update)
		role.DELETE("/:id", roleController.Delete)
	}
}
