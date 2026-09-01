package controllers

import (
	"net/http"

	"backend_koperasi/internal/requests"
	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentService *services.PaymentService
}

func NewPaymentController(service *services.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: service,
	}
}

// ======================================
// GET /api/v1/payments (Admin)
// ======================================

func (c *PaymentController) GetAll(ctx *gin.Context) {
	payments, err := c.paymentService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pembayaran berhasil diambil.",
		"data":    payments,
	})
}

// ======================================
// GET /api/v1/payments/:id
// ======================================

func (c *PaymentController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	payment, err := c.paymentService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pembayaran berhasil diambil.",
		"data":    payment,
	})
}

// ======================================
// GET /api/v1/payments/order/:order_id
// ======================================

func (c *PaymentController) GetByOrderID(ctx *gin.Context) {
	orderID := ctx.Param("order_id")

	payment, err := c.paymentService.GetByOrderID(orderID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pembayaran berhasil diambil.",
		"data":    payment,
	})
}

// ======================================
// POST /api/v1/payments (Create Payment)
// ======================================

func (c *PaymentController) Create(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var req requests.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid: " + err.Error(),
		})
		return
	}

	payment, err := c.paymentService.Create(req, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Entri pembayaran berhasil dibuat.",
		"data":    payment,
	})
}

// ======================================
// POST /api/v1/payments/pay (Confirm Payment)
// ======================================

func (c *PaymentController) Pay(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var req requests.PayOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid: " + err.Error(),
		})
		return
	}

	payment, err := c.paymentService.Pay(req, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pembayaran pesanan berhasil dikonfirmasi.",
		"data":    payment,
	})
}

// ======================================
// PUT /api/v1/payments/:id (Admin)
// ======================================

func (c *PaymentController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req requests.UpdatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid: " + err.Error(),
		})
		return
	}

	payment, err := c.paymentService.Update(id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pembayaran berhasil diperbarui.",
		"data":    payment,
	})
}

// ======================================
// DELETE /api/v1/payments/:id (Admin)
// ======================================

func (c *PaymentController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.paymentService.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pembayaran berhasil dihapus.",
	})
}
