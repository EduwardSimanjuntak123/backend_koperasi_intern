package controllers

import (
	"net/http"
	"strconv"

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
	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user id parameter",
		})
		return
	}

	favorites, err := c.favoriteService.GetUserFavorites(uint(userID))
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
		UserID    uint `json:"user_id" binding:"required"`
		ProductID uint `json:"product_id" binding:"required"`
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

	userID, err1 := strconv.Atoi(userIDStr)
	productID, err2 := strconv.Atoi(productIDStr)

	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user id or product id parameter",
		})
		return
	}

	if err := c.favoriteService.RemoveFromFavorite(uint(userID), uint(productID)); err != nil {
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
