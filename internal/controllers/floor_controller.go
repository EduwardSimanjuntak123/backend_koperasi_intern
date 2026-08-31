package controllers

import (
	"net/http"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type FloorController struct {
	floorService *services.FloorService
}

func NewFloorController(service *services.FloorService) *FloorController {
	return &FloorController{
		floorService: service,
	}
}

// ======================================
// Get All Floors
// ======================================

func (c *FloorController) GetAll(ctx *gin.Context) {

	floors, err := c.floorService.GetAllFloors()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Floors retrieved successfully",
		"data":    floors,
	})
}

// ======================================
// Get Floor By ID
// ======================================

func (c *FloorController) GetByID(ctx *gin.Context) {

	id := ctx.Param("id")

	floor, err := c.floorService.GetFloorByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Floor retrieved successfully",
		"data":    floor,
	})
}

// ======================================
// Get Floors By Building
// ======================================

func (c *FloorController) GetByBuildingID(ctx *gin.Context) {

	buildingID := ctx.Param("building_id")

	floors, err := c.floorService.GetFloorsByBuilding(buildingID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Floors retrieved successfully",
		"data":    floors,
	})
}

// ======================================
// Create Floor
// ======================================

func (c *FloorController) Create(ctx *gin.Context) {

	var floor models.Floor

	if err := ctx.ShouldBindJSON(&floor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if err := c.floorService.CreateFloor(&floor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Floor created successfully",
		"data":    floor,
	})
}

// ======================================
// Update Floor
// ======================================

func (c *FloorController) Update(ctx *gin.Context) {

	id := ctx.Param("id")

	var floor models.Floor

	if err := ctx.ShouldBindJSON(&floor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	floor.ID = id

	if err := c.floorService.UpdateFloor(&floor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Floor updated successfully",
		"data":    floor,
	})
}

// ======================================
// Delete Floor
// ======================================

func (c *FloorController) Delete(ctx *gin.Context) {

	id := ctx.Param("id")

	if err := c.floorService.DeleteFloor(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Floor deleted successfully",
	})
}
