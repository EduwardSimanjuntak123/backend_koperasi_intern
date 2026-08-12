package controllers

import (
	"net/http"
	"strconv"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type RolesController struct {
	rolesService *services.RolesService
}

func NewRolesController(service *services.RolesService) *RolesController {
	return &RolesController{
		rolesService: service,
	}
}

// =====================================
// GET /api/roles
// =====================================
func (c *RolesController) GetAll(ctx *gin.Context) {

	roles, err := c.rolesService.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Roles retrieved successfully",
		"data":    roles,
	})
}

// =====================================
// GET /api/roles/:id
// =====================================
func (c *RolesController) GetByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid role id",
		})
		return
	}

	role, err := c.rolesService.GetByID(uint(id))

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role retrieved successfully",
		"data":    role,
	})
}

// =====================================
// POST /api/roles
// =====================================
func (c *RolesController) Create(ctx *gin.Context) {

	var role models.Roles

	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := c.rolesService.Create(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Role created successfully",
		"data":    role,
	})
}

// =====================================
// PUT /api/roles/:id
// =====================================
func (c *RolesController) Update(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid role id",
		})
		return
	}

	var role models.Roles

	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := c.rolesService.Update(uint(id), &role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role updated successfully",
	})
}

// =====================================
// DELETE /api/roles/:id
// =====================================
func (c *RolesController) Delete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid role id",
		})
		return
	}

	if err := c.rolesService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role deleted successfully",
	})
}
