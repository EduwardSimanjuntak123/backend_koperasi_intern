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

func NewFavoriteController(favoriteService *services.FavoriteService) *FavoriteController {
	return &FavoriteController{
		favoriteService: favoriteService,
	}
}

// ======================================================
// GET /favorites
// Mengambil semua produk favorit milik user
// ======================================================
func (fc *FavoriteController) GetMyFavorites(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userID := userIDValue.(uint)

	favorites, err := fc.favoriteService.GetFavorites(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Favorites retrieved successfully",
		"data":    favorites,
	})
}

// ======================================================
// POST /favorites/:product_id
// Menambahkan produk ke favorit
// ======================================================
func (fc *FavoriteController) AddFavorite(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userID := userIDValue.(uint)

	productID64, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid product ID",
		})
		return
	}

	err = fc.favoriteService.Add(userID, uint(productID64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Product added to favorites",
	})
}

// ======================================================
// DELETE /favorites/:product_id
// Menghapus produk dari favorit
// ======================================================
func (fc *FavoriteController) RemoveFavorite(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userID := userIDValue.(uint)

	productID64, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid product ID",
		})
		return
	}

	err = fc.favoriteService.Remove(userID, uint(productID64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product removed from favorites",
	})
}
