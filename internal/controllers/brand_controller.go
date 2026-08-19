package controllers

import (
	"net/http"
	"strconv"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type BrandController struct {
	brandService *services.BrandService
}

func NewBrandController(service *services.BrandService) *BrandController {
	return &BrandController{
		brandService: service,
	}
}

// =====================================
// GET /api/brands
// =====================================
func (c *BrandController) GetAll(ctx *gin.Context) {

	brands, err := c.brandService.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Brands retrieved successfully",
		"data":    brands,
	})
}

// =====================================
// GET /api/brands/:id
// =====================================
func (c *BrandController) GetByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid brand id",
		})
		return
	}

	brand, err := c.brandService.GetByID(uint(id))

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Brand retrieved successfully",
		"data":    brand,
	})
}

// =====================================
// POST /api/brands
// =====================================
func (c *BrandController) Create(ctx *gin.Context) {

	var brand models.Brand

	if err := ctx.ShouldBindJSON(&brand); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := c.brandService.Create(&brand); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Brand created successfully",
		"data":    brand,
	})
}

// =====================================
// PUT /api/brands/:id
// =====================================
func (c *BrandController) Update(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid brand id",
		})
		return
	}

	var brand models.Brand

	if err := ctx.ShouldBindJSON(&brand); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := c.brandService.Update(uint(id), &brand); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Brand updated successfully",
	})
}

// =====================================
// DELETE /api/brands/:id
// =====================================
func (c *BrandController) Delete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid brand id",
		})
		return
	}

	if err := c.brandService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Brand	 deleted successfully",
	})
}
