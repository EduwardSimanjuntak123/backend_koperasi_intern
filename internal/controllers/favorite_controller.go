package controllers

import (
	"net/http"

	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type FavoriteController struct {
	favoriteService *services.FavoriteService
}

func NewFavoriteController(service *services.FavoriteService) *FavoriteController {
	return &FavoriteController{
		favoriteService: service,
	}
}

// GET /api/v1/favorites/user/:user_id
func (c *FavoriteController) GetUserFavorites(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user id parameter",
		})
		return
	}

	favorites, err := c.favoriteService.GetUserFavorites(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Favorites retrieved successfully",
		"data":    favorites,
	})
}

// POST /api/v1/favorites
func (c *FavoriteController) AddToFavorite(ctx *gin.Context) {
	var req struct {
		UserID    string `json:"user_id" binding:"required"`
		ProductID string `json:"product_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if err := c.favoriteService.AddToFavorite(req.UserID, req.ProductID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Product added to favorites successfully",
	})
}

// DELETE /api/v1/favorites/user/:user_id/product/:product_id
func (c *FavoriteController) RemoveFromFavorite(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	productIDStr := ctx.Param("product_id")

	if userIDStr == "" || productIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user id or product id parameter",
		})
		return
	}

	if err := c.favoriteService.RemoveFromFavorite(userIDStr, productIDStr); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product removed from favorites successfully",
	})
}
