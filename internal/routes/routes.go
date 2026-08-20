package routes

import (
	"backend_koperasi/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	productController *controllers.ProductController,
	categoryController *controllers.CategoryProductController,
	userController *controllers.UserController,
	rolesController *controllers.RolesController,
	storeController *controllers.StoreController,
	storeMemberController *controllers.StoreMemberController,
	favoriteController *controllers.FavoriteController,
	authController *controllers.AuthController,
	cartController *controllers.CartController,
	unitController *controllers.UnitController,
	brandController *controllers.BrandController,
) {

	api := router.Group("/api/v1")

	RegisterProductRoutes(api, productController)
	RegisterCategoryRoutes(api, categoryController)
	RegisterUserRoutes(api, userController)
	RegisterRoleRoutes(api, rolesController)
	RegisterStoreRoutes(api, storeController)
	RegisterStoreMemberRoutes(api, storeMemberController)
	RegisterFavoriteRoutes(api, favoriteController)
	RegisterAuthRoutes(api, authController)
	RegisterCartRoutes(api, cartController)
	RegisterBrandRoutes(api, brandController)
	RegisterUnitRoutes(api, unitController)

}
