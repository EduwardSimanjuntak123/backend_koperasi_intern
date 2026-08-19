package controllers

import (
	"net/http"
	"strconv"

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
// GET /cart
// Mengambil keranjang aktif milik user beserta isinya
// =====================================
func (cc *CartController) GetCart(c *gin.Context) {
	// Mengambil user_id dari JWT via AuthMiddleware
	userID := c.MustGet("user_id").(uint)

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
		"message": "Cart retrieved successfully",
		"data":    cart,
	})
}

// =====================================
// POST /cart/items
// Menambahkan produk ke keranjang
// =====================================
func (cc *CartController) AddToCart(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	// Mendefinisikan struktur request body
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: product_id and valid quantity are required",
		})
		return
	}

	err := cc.cartService.AddToCart(userID, req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Product added to cart successfully",
	})
}

// =====================================
// PUT /cart/items/:item_id
// Mengubah jumlah (quantity) barang di keranjang
// =====================================
func (cc *CartController) UpdateCartItemQuantity(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	itemIDStr := c.Param("item_id")
	itemID, err := strconv.Atoi(itemIDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid cart item ID",
		})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: valid quantity is required",
		})
		return
	}

	err = cc.cartService.UpdateItemQuantity(userID, uint(itemID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cart item quantity updated successfully",
	})
}

// =====================================
// DELETE /cart/items/:item_id
// Menghapus satu jenis barang dari keranjang
// =====================================
func (cc *CartController) RemoveFromCart(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	itemIDStr := c.Param("item_id")
	itemID, err := strconv.Atoi(itemIDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid cart item ID",
		})
		return
	}

	err = cc.cartService.RemoveItem(userID, uint(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item removed from cart successfully",
	})
}