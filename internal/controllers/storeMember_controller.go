package controllers

import (
	"backend_koperasi/internal/models"
	"backend_koperasi/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StoreMemberController struct {
	storeMemberService *services.StoreMemberService
}

func NewStoreMemberController(service *services.StoreMemberService) *StoreMemberController {
	return &StoreMemberController{
		storeMemberService: service,
	}
}

func (c *StoreMemberController) GetAll(ctx *gin.Context) {

	storeMembers, err := c.storeMemberService.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    storeMembers,
	})
}

func (c *StoreMemberController) GetByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	storeMember, err := c.storeMemberService.GetByID(uint(id))

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, storeMember)
}

func (c *StoreMemberController) Create(ctx *gin.Context) {

	var storeMember models.StoreMember

	if err := ctx.ShouldBindJSON(&storeMember); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := c.storeMemberService.Create(&storeMember); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, storeMember)
}

func (c *StoreMemberController) Update(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	var storeMember models.StoreMember

	if err := ctx.ShouldBindJSON(&storeMember); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := c.storeMemberService.Update(uint(id), &storeMember); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "updated successfully",
	})
}

func (c *StoreMemberController) Delete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	if err := c.storeMemberService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
	})
}
