package routes

import (
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFavoriteRoutes(
	api *gin.RouterGroup,
	favoriteController *controllers.FavoriteController,
) {

	favorite := api.Group("/favorites")

	favorite.Use(
		middleware.AuthMiddleware(),
		middleware.BuyerMiddleware(),
	)

	{
		favorite.GET("", favoriteController.GetUserFavorites)
		favorite.POST("", favoriteController.AddToFavorite)
		favorite.DELETE("/:product_id", favoriteController.RemoveFromFavorite)
	}
}
