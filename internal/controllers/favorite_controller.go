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

// GET /api/v1/favorites
func (c *FavoriteController) GetUserFavorites(ctx *gin.Context) {

	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	favorites, err := c.favoriteService.GetUserFavorites(userID)
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
		ProductID string `json:"product_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	if err := c.favoriteService.AddToFavorite(userID, req.ProductID); err != nil {
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

func (c *FavoriteController) RemoveFromFavorite(ctx *gin.Context) {

	// Ambil user_id yang telah disimpan oleh AuthMiddleware
	// AuthMiddleware membaca JWT dari HttpOnly Cookie,
	// memverifikasi token, lalu menyimpan user_id ke gin.Context.
	userID := ctx.GetString("user_id")

	// Jika user_id tidak ada, berarti user belum login
	// atau token tidak valid.
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	// Ambil product_id dari parameter URL.
	// Contoh:
	// DELETE /api/v1/favorites/123
	// favoriteID = "123"
	favoriteID := ctx.Param("product_id")

	// Hapus produk dari daftar favorite milik user yang sedang login.
	// userID berasal dari JWT, bukan dari request,
	// sehingga user tidak dapat menghapus favorite milik user lain.
	if err := c.favoriteService.RemoveFromFavorite(userID, favoriteID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Berhasil menghapus favorite.
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product removed from favorites successfully",
	})
}
