package controllers

import (
	"net/http"

	"backend_koperasi/internal/models"
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
// GET /api/v1/orders
// ======================================

func (c *OrderController) GetAllOrders(ctx *gin.Context) {

	orders, err := c.orderService.GetAllOrders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Orders retrieved successfully",
		"data":    orders,
	})
}

// ======================================
// GET /api/v1/orders/:id
// ======================================

func (c *OrderController) GetOrderByID(ctx *gin.Context) {

	id := ctx.Param("id")

	order, err := c.orderService.GetOrderByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order retrieved successfully",
		"data":    order,
	})
}

// ======================================
// GET /api/v1/my-orders
// ======================================

func (c *OrderController) GetMyOrders(ctx *gin.Context) {

	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	orders, err := c.orderService.GetOrdersByUser(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Orders retrieved successfully",
		"data":    orders,
	})
}

// ======================================
// POST /api/v1/orders
// ======================================

func (c *OrderController) CreateOrder(ctx *gin.Context) {

	var order models.Order

	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if err := c.orderService.CreateOrder(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Order created successfully",
		"data":    order,
	})
}

// ======================================
// PUT /api/v1/orders/:id
// ======================================

func (c *OrderController) UpdateOrder(ctx *gin.Context) {

	id := ctx.Param("id")

	var order models.Order

	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	order.ID = id

	if err := c.orderService.UpdateOrder(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order updated successfully",
	})
}

// ======================================
// DELETE /api/v1/orders/:id
// ======================================

func (c *OrderController) DeleteOrder(ctx *gin.Context) {

	id := ctx.Param("id")

	if err := c.orderService.DeleteOrder(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order deleted successfully",
	})
}
