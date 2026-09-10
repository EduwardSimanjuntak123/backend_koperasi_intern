package controllers

import (
	"backend_koperasi/internal/models"
	"backend_koperasi/internal/services"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService *services.ProductService
}

func NewProductController(service *services.ProductService) *ProductController {
	return &ProductController{
		productService: service,
	}
}

func (c *ProductController) GetAll(ctx *gin.Context) {
	filter, err := productFilterFromQuery(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	products, total, err := c.productService.GetAll(filter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Products retrieved successfully",
		"data":    products,
		"meta": gin.H{
			"page":        filter.Page,
			"limit":       filter.Limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (c *ProductController) GetDiscounted(ctx *gin.Context) {
	filter, err := productFilterFromQuery(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	products, total, err := c.productService.GetDiscounted(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk diskon berhasil diambil.",
		"data":    products,
		"meta": gin.H{
			"page": filter.Page, "limit": filter.Limit, "total": total,
			"total_pages": int(math.Ceil(float64(total) / float64(filter.Limit))),
		},
	})
}

func (c *ProductController) GetBestSelling(ctx *gin.Context) {
	filter, err := productFilterFromQuery(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	periodDays, err := queryInt(ctx, "period_days", 30)
	if err != nil || periodDays < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "parameter period_days harus berupa angka positif"})
		return
	}

	products, err := c.productService.GetBestSelling(filter, periodDays)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk terlaris berhasil diambil.",
		"data":    products,
		"meta":    gin.H{"limit": filter.Limit, "period_days": periodDays},
	})
}

func productFilterFromQuery(ctx *gin.Context) (models.ProductFilter, error) {
	page, err := queryInt(ctx, "page", 1)
	if err != nil || page < 1 {
		return models.ProductFilter{}, errors.New("parameter page harus berupa angka positif")
	}
	limit, err := queryInt(ctx, "limit", 10)
	if err != nil || limit < 1 {
		return models.ProductFilter{}, errors.New("parameter limit harus berupa angka positif")
	}

	filter := models.ProductFilter{
		Search:      ctx.Query("search"),
		Page:        page,
		Limit:       limit,
		StockStatus: strings.ToLower(ctx.Query("stock_status")),
		Expired:     strings.ToLower(ctx.Query("expired")),
		SortBy:      strings.ToLower(ctx.Query("sort_by")),
		SortOrder:   strings.ToLower(ctx.Query("sort_order")),
	}

	for name, target := range map[string]**string{
		"brand_id": &filter.BrandID, "category_id": &filter.CategoryID,
		"unit_id": &filter.UnitID, "store_id": &filter.StoreID,
	} {
		if value := strings.TrimSpace(ctx.Query(name)); value != "" {
			*target = &value
		}
	}

	if value := ctx.Query("min_price"); value != "" {
		parsed, parseErr := strconv.ParseFloat(value, 64)
		if parseErr != nil || parsed < 0 {
			return models.ProductFilter{}, errors.New("parameter min_price tidak valid")
		}
		filter.MinPrice = &parsed
	}
	if value := ctx.Query("max_price"); value != "" {
		parsed, parseErr := strconv.ParseFloat(value, 64)
		if parseErr != nil || parsed < 0 {
			return models.ProductFilter{}, errors.New("parameter max_price tidak valid")
		}
		filter.MaxPrice = &parsed
	}
	if filter.MinPrice != nil && filter.MaxPrice != nil && *filter.MinPrice > *filter.MaxPrice {
		return models.ProductFilter{}, errors.New("min_price tidak boleh lebih besar dari max_price")
	}

	if value := strings.TrimSpace(ctx.Query("inventory_movement")); value != "" {
		movement := models.InventoryMovement(strings.ToUpper(value))
		if movement != models.FastMoving && movement != models.SlowMoving {
			return models.ProductFilter{}, errors.New("inventory_movement harus FAST_MOVING atau SLOW_MOVING")
		}
		filter.InventoryMovement = &movement
	}
	if value := strings.TrimSpace(ctx.Query("badge")); value != "" {
		badge := models.ProductBadge(strings.ToUpper(value))
		if badge != models.BadgeNew && badge != models.BadgeBestSeller {
			return models.ProductFilter{}, errors.New("badge harus NEW atau BEST_SELLER")
		}
		filter.Badge = &badge
	}
	if filter.StockStatus != "" && filter.StockStatus != "in_stock" && filter.StockStatus != "out_of_stock" && filter.StockStatus != "low_stock" {
		return models.ProductFilter{}, errors.New("stock_status harus in_stock, out_of_stock, atau low_stock")
	}
	if filter.Expired != "" && filter.Expired != "expired" && filter.Expired != "not_expired" {
		return models.ProductFilter{}, errors.New("expired harus expired atau not_expired")
	}
	if filter.SortOrder != "" && filter.SortOrder != "asc" && filter.SortOrder != "desc" {
		return models.ProductFilter{}, errors.New("sort_order harus asc atau desc")
	}
	if value := ctx.Query("created_from"); value != "" {
		parsed, parseErr := time.Parse("2006-01-02", value)
		if parseErr != nil {
			return models.ProductFilter{}, errors.New("created_from harus berformat YYYY-MM-DD")
		}
		filter.CreatedFrom = &parsed
	}
	if value := ctx.Query("created_to"); value != "" {
		parsed, parseErr := time.Parse("2006-01-02", value)
		if parseErr != nil {
			return models.ProductFilter{}, errors.New("created_to harus berformat YYYY-MM-DD")
		}
		parsed = parsed.Add(24*time.Hour - time.Nanosecond)
		filter.CreatedTo = &parsed
	}
	if filter.CreatedFrom != nil && filter.CreatedTo != nil && filter.CreatedFrom.After(*filter.CreatedTo) {
		return models.ProductFilter{}, errors.New("created_from tidak boleh setelah created_to")
	}

	return filter, nil
}

func queryInt(ctx *gin.Context, name string, fallback int) (int, error) {
	value := ctx.Query(name)
	if value == "" {
		return fallback, nil
	}
	return strconv.Atoi(value)
}

func (c *ProductController) GetLowStock(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	products, total, err := c.productService.GetAll(models.ProductFilter{
		StockStatus: "low_stock",
		Page:        page,
		Limit:       limit,
		SortBy:      "stock",
		SortOrder:   "asc",
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk dengan stok menipis berhasil diambil.",
		"data":    products,
		"meta": gin.H{
			"page": page, "limit": limit, "total": total,
			"total_pages": int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

func (c *ProductController) GetByID(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	product, err := c.productService.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product retrieved successfully",
		"data":    product,
	})
}

func (c *ProductController) Create(ctx *gin.Context) {

	var product models.Product

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := c.productService.Create(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	createdProduct, err := c.productService.GetByID(product.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve created product",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Product created successfully",
		"data":    createdProduct,
	})
}

func (c *ProductController) Update(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	var product models.Product

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := c.productService.Update(id, &product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product updated successfully",
	})
}

func (c *ProductController) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	if err := c.productService.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product deleted successfully",
	})
}
