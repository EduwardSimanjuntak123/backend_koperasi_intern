package routes

import (
	"backend_koperasi/internal/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes mendaftarkan seluruh endpoint User.
func RegisterUserRoutes(
	router *gin.RouterGroup,
	userController *controllers.UserController,
) {

	user := router.Group("/users")
	{
		user.GET("", userController.GetAll)
		user.GET("/:id", userController.GetByID)
		user.POST("", userController.Create)
		user.PUT("/:id", userController.Update)
		user.DELETE("/:id", userController.Delete)
	}
}
