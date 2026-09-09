package controllers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"backend_koperasi/internal/requests"
	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService *services.OrderService
}

func NewOrderController(service *services.OrderService) *OrderController {
	return &OrderController{
		orderService: service,
	}
}

// ======================================
// POST /api/v1/orders (Create Order)
// ======================================

func (c *OrderController) Create(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Silakan login terlebih dahulu.",
		})
		return
	}

	var req requests.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid: " + err.Error(),
		})
		return
	}

	order, err := c.orderService.Create(userID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pesanan berhasil dibuat.",
		"data":    order,
	})
}

// ======================================
// POST /api/v1/orders/checkout-cart
// ======================================

func (c *OrderController) CheckoutCart(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Silakan login terlebih dahulu.",
		})
		return
	}

	var req requests.CheckoutCartRequest
	_ = ctx.ShouldBindJSON(&req)

	order, err := c.orderService.CheckoutCart(userID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pesanan berhasil dibuat dari keranjang belanja.",
		"data":    order,
	})
}

// ======================================
// GET /api/v1/orders (Buyer's Orders)
// ======================================

func (c *OrderController) GetMyOrders(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Silakan login terlebih dahulu.",
		})
		return
	}

	orders, err := c.orderService.GetByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar pesanan Anda berhasil diambil.",
		"data":    orders,
	})
}

// ======================================
// GET /api/v1/orders/:id (Detail Order)
// ======================================

func (c *OrderController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")
	roleID, _ := ctx.Get("role_id")

	isAdmin := false
	if rID, ok := roleID.(uint); ok && rID == 1 {
		isAdmin = true
	}

	order, err := c.orderService.GetByID(id, userID, isAdmin)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail pesanan berhasil diambil.",
		"data":    order,
	})
}

// ======================================
// POST /api/v1/orders/:id/cancel (Cancel)
// ======================================

func (c *OrderController) Cancel(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")
	roleID, _ := ctx.Get("role_id")

	isAdmin := false
	if rID, ok := roleID.(uint); ok && rID == 1 {
		isAdmin = true
	}

	var req requests.CancelOrderRequest
	_ = ctx.ShouldBindJSON(&req)

	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}

	if err := c.orderService.Cancel(id, userID, reason, isAdmin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pesanan berhasil dibatalkan dan stok telah dikembalikan.",
	})
}

// ======================================
// GET /api/v1/orders/all (Admin: All Orders)
// ======================================

func (c *OrderController) GetAll(ctx *gin.Context) {
	orders, err := c.orderService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Seluruh data pesanan berhasil diambil.",
		"data":    orders,
	})
}

func (c *OrderController) GetByStatus(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	orders, total, err := c.orderService.GetByStatus(ctx.Query("status"), page, limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar pesanan berdasarkan status berhasil diambil.",
		"data":    orders,
		"meta": gin.H{
			"page": page, "limit": limit, "total": total,
			"total_pages": int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

func (c *OrderController) GetDashboardSummary(ctx *gin.Context) {
	summary, err := c.orderService.GetDashboardSummary(time.Now())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Ringkasan pesanan berhasil diambil.", "data": summary})
}

func (c *OrderController) GetTodayRevenue(ctx *gin.Context) {
	revenue, err := c.orderService.GetTodayRevenue(time.Now())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Omzet hari ini berhasil diambil.", "data": gin.H{"omzet_hari_ini": revenue}})
}

// ======================================
// PUT /api/v1/orders/:id/status (Admin)
// ======================================

func (c *OrderController) UpdateStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	adminID := ctx.GetString("user_id")

	var req requests.UpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid: " + err.Error(),
		})
		return
	}

	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	if err := c.orderService.UpdateStatus(id, req.Status, notes, adminID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status pesanan berhasil diperbarui.",
	})
}

// ======================================
// PUT /api/v1/orders/:id/courier (Admin)
// ======================================

func (c *OrderController) AssignCourier(ctx *gin.Context) {
	id := ctx.Param("id")
	adminID := ctx.GetString("user_id")

	var req requests.AssignCourierRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid: " + err.Error(),
		})
		return
	}

	if err := c.orderService.AssignCourier(id, req.CourierID, adminID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kurir berhasil ditugaskan untuk pesanan.",
	})
}

// ======================================
// DELETE /api/v1/orders/:id (Admin)
// ======================================

func (c *OrderController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.orderService.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pesanan beserta relasinya berhasil dihapus.",
	})
}
