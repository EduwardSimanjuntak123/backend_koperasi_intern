package controllers

import (
	"net/http"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type CourierController struct {
	courierService *services.CourierService
}

func NewCourierController(service *services.CourierService) *CourierController {
	return &CourierController{
		courierService: service,
	}
}

// =====================================
// GET /api/couriers
// =====================================
func (c *CourierController) GetAll(ctx *gin.Context) {

	couriers, err := c.courierService.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Couriers retrieved successfully",
		"data":    couriers,
	})
}

// =====================================
// GET /api/couriers/:id
// =====================================
func (c *CourierController) GetByID(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid courier id",
		})
		return
	}

	courier, err := c.courierService.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Courier retrieved successfully",
		"data":    courier,
	})
}

// =====================================
// POST /api/couriers
// =====================================
func (c *CourierController) Create(ctx *gin.Context) {

	var courier models.Courier

	if err := ctx.ShouldBindJSON(&courier); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data yang dikirim tidak valid.",
		})
		return
	}

	if err := c.courierService.Create(&courier); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Courier created successfully",
		"data":    courier,
	})
}

// =====================================
// PUT /api/couriers/:id
// =====================================
func (c *CourierController) Update(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid courier id",
		})
		return
	}

	var courier models.Courier

	if err := ctx.ShouldBindJSON(&courier); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data yang dikirim tidak valid.",
		})
		return
	}

	if err := c.courierService.Update(id, &courier); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Courier updated successfully",
	})
}

// =====================================
// DELETE /api/couriers/:id
// =====================================
func (c *CourierController) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid courier id",
		})
		return
	}

	if err := c.courierService.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Courier deleted successfully",
	})
}
