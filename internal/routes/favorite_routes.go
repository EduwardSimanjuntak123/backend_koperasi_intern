package routes

import (
	"backend_koperasi/internal/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterFavoriteRoutes mendaftarkan seluruh endpoint Favorite.
func RegisterFavoriteRoutes(
	router *gin.RouterGroup,
	favoriteController *controllers.FavoriteController,
) {
	favorite := router.Group("/favorites")
	{
		favorite.GET("/user/:user_id", favoriteController.GetUserFavorites)
		favorite.POST("", favoriteController.AddToFavorite)
		favorite.DELETE("/user/:user_id/product/:product_id", favoriteController.RemoveFromFavorite)
	}
}