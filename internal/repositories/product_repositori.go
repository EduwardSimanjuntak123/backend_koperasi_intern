package repositories

import (
	"backend_koperasi/internal/models"

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
func (r *ProductRepository) FindAll(search string) ([]models.Product, error) {

	var products []models.Product

	query := r.db.
		Preload("Category").
		Preload("Brand").
		Preload("Unit").
		Preload("Store")

	if search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%")
	}

	err := query.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {

	var product models.Product

	err := r.db.
		Preload("Category").
		Preload("Brand").
		Preload("Unit").
		Preload("Store").
		First(&product, id).Error

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

func (r *ProductRepository) FindByCategoryID(categoryID uint) ([]models.Product, error) {

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

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
