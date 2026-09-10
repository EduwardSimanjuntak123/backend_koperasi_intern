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
		query = query.Where("stock <= ?", 3)
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

func (r *ProductRepository) FindDiscounted(filter models.ProductFilter) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).
		Where("promotion_price IS NOT NULL AND promotion_price > 0 AND promotion_price < price")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("products.created_at DESC").
		Preload("Brand").Preload("Unit").Preload("Category").Preload("Store").
		Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit)

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepository) FindBestSelling(filter models.ProductFilter, periodDays int) ([]models.ProductSales, error) {
	type salesRow struct {
		ProductID string
		TotalSold int
	}

	query := r.db.Table("products AS p").
		Select("p.id AS product_id, COALESCE(SUM(oi.qty), 0) AS total_sold").
		Joins("JOIN order_items AS oi ON oi.product_id = p.id").
		Joins("JOIN orders AS o ON o.id = oi.order_id").
		Where("o.status IN ?", []models.OrderStatus{
			models.OrderPaid, models.OrderPacking, models.OrderShipping,
			models.OrderDelivered, models.OrderCompleted,
		}).
		Where("o.created_at >= ?", time.Now().AddDate(0, 0, -periodDays))

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("(p.name ILIKE ? OR p.sku ILIKE ? OR p.barcode ILIKE ?)", search, search, search)
	}
	if filter.BrandID != nil {
		query = query.Where("p.brand_id = ?", *filter.BrandID)
	}
	if filter.CategoryID != nil {
		query = query.Where("p.category_id = ?", *filter.CategoryID)
	}
	if filter.UnitID != nil {
		query = query.Where("p.unit_id = ?", *filter.UnitID)
	}
	if filter.StoreID != nil {
		query = query.Where("p.store_id = ?", *filter.StoreID)
	}
	if filter.MinPrice != nil {
		query = query.Where("p.price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("p.price <= ?", *filter.MaxPrice)
	}
	if filter.StockStatus == "in_stock" {
		query = query.Where("p.stock > 0")
	} else if filter.StockStatus == "out_of_stock" {
		query = query.Where("p.stock <= 0")
	} else if filter.StockStatus == "low_stock" {
		query = query.Where("p.stock <= 3")
	}

	var rows []salesRow
	if err := query.Group("p.id").Order("total_sold DESC, p.created_at DESC").Limit(filter.Limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []models.ProductSales{}, nil
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ProductID)
	}
	var products []models.Product
	if err := r.db.Preload("Brand").Preload("Unit").Preload("Category").Preload("Store").Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}
	productByID := make(map[string]models.Product, len(products))
	for _, product := range products {
		productByID[product.ID] = product
	}
	result := make([]models.ProductSales, 0, len(rows))
	for _, row := range rows {
		if product, ok := productByID[row.ProductID]; ok {
			result = append(result, models.ProductSales{Product: product, TotalSold: row.TotalSold})
		}
	}
	return result, nil
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
