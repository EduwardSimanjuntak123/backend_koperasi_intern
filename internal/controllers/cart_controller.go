package controllers

import (
	"net/http"

	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	cartService *services.CartService
}

func NewCartController(cartService *services.CartService) *CartController {
	return &CartController{
		cartService: cartService,
	}
}

// =====================================
// GET /api/v1/cart
// Mengambil keranjang aktif milik user beserta isinya
// =====================================
func (cc *CartController) GetCart(c *gin.Context) {
	// Mengambil user_id dari JWT via AuthMiddleware
	userID := c.MustGet("user_id").(string)

	cart, err := cc.cartService.GetCartByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Keranjang berhasil diambil",
		"data":    cart,
	})
}

// =====================================
// POST /api/v1/cart/items
// Menambahkan produk ke keranjang
// =====================================
func (cc *CartController) AddToCart(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	// Mendefinisikan struktur request body
	var req struct {
		ProductID string `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Permintaan tidak valid: product_id dan quantity (minimal 1) wajib diisi",
		})
		return
	}

	err := cc.cartService.AddToCart(userID, req.ProductID, req.Quantity)
	if err != nil {
		statusCode := services.ResolveCartErrorCode(err)
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Produk berhasil ditambahkan ke keranjang",
	})
}

// =====================================
// PUT /api/v1/cart/items/:item_id
// Mengubah jumlah (quantity) barang di keranjang
// =====================================
func (cc *CartController) UpdateCartItemQuantity(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	itemIDStr := c.Param("item_id")
	if itemIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID item keranjang tidak valid",
		})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Permintaan tidak valid: quantity wajib diisi dan minimal bernilai 1",
		})
		return
	}

	err := cc.cartService.UpdateItemQuantity(userID, itemIDStr, req.Quantity)
	if err != nil {
		statusCode := services.ResolveCartErrorCode(err)
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jumlah item keranjang berhasil diperbarui",
	})
}

// =====================================
// DELETE /api/v1/cart/items/:item_id
// Menghapus satu jenis barang dari keranjang
// =====================================
func (cc *CartController) RemoveFromCart(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	itemIDStr := c.Param("item_id")
	if itemIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID item keranjang tidak valid",
		})
		return
	}

	err := cc.cartService.RemoveItem(userID, itemIDStr)
	if err != nil {
		statusCode := services.ResolveCartErrorCode(err)
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item berhasil dihapus dari keranjang",
	})
}

// =====================================
// DELETE /api/v1/cart
// Mengosongkan seluruh isi keranjang (clear cart)
// =====================================
func (cc *CartController) ClearCart(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	err := cc.cartService.ClearCart(userID)
	if err != nil {
		statusCode := services.ResolveCartErrorCode(err)
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Keranjang berhasil dikosongkan",
	})
}
