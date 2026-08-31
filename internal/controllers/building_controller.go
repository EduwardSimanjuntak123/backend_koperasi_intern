package controllers

import (
	"net/http"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type BuildingController struct {
	buildingService *services.BuildingService
}

func NewBuildingController(service *services.BuildingService) *BuildingController {
	return &BuildingController{
		buildingService: service,
	}
}

// ======================================
// Get All Buildings
// ======================================

func (c *BuildingController) GetAll(ctx *gin.Context) {

	buildings, err := c.buildingService.GetAllBuildings()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Buildings retrieved successfully",
		"data":    buildings,
	})
}

// ======================================
// Get Building By ID
// ======================================

func (c *BuildingController) GetByID(ctx *gin.Context) {

	id := ctx.Param("id")

	building, err := c.buildingService.GetBuildingByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Building retrieved successfully",
		"data":    building,
	})
}

// ======================================
// Get Floors By Building
// ======================================

// ======================================
// Create Floor
// ======================================

func (c *BuildingController) Create(ctx *gin.Context) {

	var building models.Building

	if err := ctx.ShouldBindJSON(&building); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if err := c.buildingService.CreateBuilding(&building); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Building created successfully",
		"data":    building,
	})
}

// ======================================
// Update Building
// ======================================

func (c *BuildingController) Update(ctx *gin.Context) {

	id := ctx.Param("id")

	var building models.Building

	if err := ctx.ShouldBindJSON(&building); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	building.ID = id

	if err := c.buildingService.UpdateBuilding(&building); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Building updated successfully",
		"data":    building,
	})
}

// ======================================
// Delete Building
// ======================================

func (c *BuildingController) Delete(ctx *gin.Context) {

	id := ctx.Param("id")

	if err := c.buildingService.DeleteBuilding(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Building deleted successfully",
	})
}
