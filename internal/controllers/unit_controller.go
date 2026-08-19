package controllers

import (
	"net/http"
	"strconv"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type UnitController struct {
	unitService *services.UnitService
}

func NewUnitController(service *services.UnitService) *UnitController {
	return &UnitController{
		unitService: service,
	}
}

// =====================================
// GET /api/units
// =====================================
func (c *UnitController) GetAll(ctx *gin.Context) {

	units, err := c.unitService.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Units retrieved successfully",
		"data":    units,
	})
}

// =====================================
// GET /api/units/:id
// =====================================
func (c *UnitController) GetByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid unit id",
		})
		return
	}

	unit, err := c.unitService.GetByID(uint(id))

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Unit retrieved successfully",
		"data":    unit,
	})
}

// =====================================
// POST /api/units
// =====================================
func (c *UnitController) Create(ctx *gin.Context) {

	var unit models.Unit

	if err := ctx.ShouldBindJSON(&unit); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := c.unitService.Create(&unit); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Unit created successfully",
		"data":    unit,
	})
}

// =====================================
// PUT /api/units/:id
// =====================================
func (c *UnitController) Update(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid unit id",
		})
		return
	}

	var unit models.Unit

	if err := ctx.ShouldBindJSON(&unit); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := c.unitService.Update(uint(id), &unit); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Unit updated successfully",
	})
}

// =====================================
// DELETE /api/units/:id
// =====================================
func (c *UnitController) Delete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid unit id",
		})
		return
	}

	if err := c.unitService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Unit deleted successfully",
	})
}
