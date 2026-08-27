package repositories

import (
	"backend_koperasi/internal/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// ini untuk mengambil semua produk dari database, dengan opsi pencarian berdasarkan nama produk.
// Fungsi ini menggunakan GORM untuk melakukan query ke database dan mengembalikan daftar produk yang ditemukan atau error jika terjadi kesalahan.
func (r *ProductRepository) FindAll(
	filter models.ProductFilter,
) ([]models.Product, int64, error) {

	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	// Search
	if filter.Search != "" {
		search := "%" + filter.Search + "%"

		query = query.Where(
			"name ILIKE ? OR sku ILIKE ? OR barcode ILIKE ?",
			search,
			search,
			search,
		)
	}

	// Relation filter
	if filter.BrandID != nil {
		query = query.Where("brand_id = ?", *filter.BrandID)
	}

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.UnitID != nil {
		query = query.Where("unit_id = ?", *filter.UnitID)
	}

	if filter.StoreID != nil {
		query = query.Where("store_id = ?", *filter.StoreID)
	}

	// Price
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}

	// Stock
	if filter.MinStock != nil {
		query = query.Where("stock >= ?", *filter.MinStock)
	}

	if filter.MaxStock != nil {
		query = query.Where("stock <= ?", *filter.MaxStock)
	}

	// Inventory movement
	if filter.InventoryMovement != nil {
		query = query.Where(
			"inventory_movement = ?",
			*filter.InventoryMovement,
		)
	}

	// Badge
	if filter.Badge != nil {
		query = query.Where("badge = ?", *filter.Badge)
	}

	// Stock status
	switch filter.StockStatus {
	case "out_of_stock":
		query = query.Where("stock <= ?", 0)

	case "in_stock":
		query = query.Where("stock > ?", 0)

	case "low_stock":
		query = query.Where(
			"min_stock IS NOT NULL AND stock <= min_stock",
		)
	}

	// Expired
	switch filter.Expired {
	case "expired":
		query = query.Where(
			"expired_date IS NOT NULL AND expired_date < ?",
			time.Now(),
		)

	case "not_expired":
		query = query.Where(
			"expired_date IS NULL OR expired_date >= ?",
			time.Now(),
		)
	}

	// Created date
	if filter.CreatedFrom != nil {
		query = query.Where(
			"created_at >= ?",
			*filter.CreatedFrom,
		)
	}

	if filter.CreatedTo != nil {
		query = query.Where(
			"created_at <= ?",
			*filter.CreatedTo,
		)
	}

	// Total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	allowedSort := map[string]string{
		"name":           "name",
		"price":          "price",
		"stock":          "stock",
		"created_at":     "created_at",
		"updated_at":     "updated_at",
		"purchase_price": "purchase_price",
	}

	sortColumn, ok := allowedSort[filter.SortBy]
	if !ok {
		sortColumn = "created_at"
	}

	sortOrder := "DESC"

	if strings.ToLower(filter.SortOrder) == "asc" {
		sortOrder = "ASC"
	}

	query = query.Order(sortColumn + " " + sortOrder)

	// Pagination
	if filter.Page < 1 {
		filter.Page = 1
	}

	if filter.Limit < 1 {
		filter.Limit = 20
	}

	offset := (filter.Page - 1) * filter.Limit

	query = query.
		Preload("Brand").
		Preload("Unit").
		Preload("Category").
		Preload("Store").
		Offset(offset).
		Limit(filter.Limit)

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepository) FindByID(id string) (*models.Product, error) {

	var product models.Product

	err := r.db.
		Preload("Category").
		Preload("Brand").
		Preload("Unit").
		Preload("Store").
		Where("id = ?", id).
		First(&product).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) FindBySlug(slug string) (*models.Product, error) {

	var product models.Product

	err := r.db.
		Where("slug = ?", slug).
		First(&product).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) FindByCategoryID(categoryID string) ([]models.Product, error) {

	var products []models.Product
	err := r.db.
		Where("category_id = ?", categoryID).
		Preload("Category").
		Find(&products).Error

	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id string) error {

	result := r.db.Where("id = ?", id).Delete(&models.Product{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ProductRepository) GenerateNextID() (string, error) {
	return generateNextPrefixedID(r.db, &models.Product{}, "P")
}
