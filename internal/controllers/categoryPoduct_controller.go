package controllers

import (
	"net/http"
	"strconv"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type CategoryProductController struct {
	categoryProductService *services.CategoryProductService
}

func NewCategoryProductController(service *services.CategoryProductService) *CategoryProductController {
	return &CategoryProductController{
		categoryProductService: service,
	}
}

// =====================================
// GET /api/categories
// =====================================
func (c *CategoryProductController) GetAll(ctx *gin.Context) {

	categories, err := c.categoryProductService.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Categories retrieved successfully",
		"data":    categories,
	})
}

// =====================================
// GET /api/categories/:id
// =====================================
func (c *CategoryProductController) GetByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid category id",
		})
		return
	}

	category, err := c.categoryProductService.GetByID(uint(id))

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category retrieved successfully",
		"data":    category,
	})
}

// =====================================
// POST /api/categories
// =====================================
func (c *CategoryProductController) Create(ctx *gin.Context) {

	var category models.CategoryProduct

	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := c.categoryProductService.Create(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Category created successfully",
		"data":    category,
	})
}

// =====================================
// PUT /api/categories/:id
// =====================================
func (c *CategoryProductController) Update(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid category id",
		})
		return
	}

	var category models.CategoryProduct

	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := c.categoryProductService.Update(uint(id), &category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category updated successfully",
	})
}

// =====================================
// DELETE /api/categories/:id
// =====================================
func (c *CategoryProductController) Delete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid category id",
		})
		return
	}

	if err := c.categoryProductService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category deleted successfully",
	})
}
