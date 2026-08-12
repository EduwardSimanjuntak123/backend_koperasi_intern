package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type CategoryProductRepository struct {
	db *gorm.DB
}

func NewCategoryProductRepository(db *gorm.DB) *CategoryProductRepository {
	return &CategoryProductRepository{
		db: db,
	}
}

// Mengambil seluruh kategori
func (r *CategoryProductRepository) FindAll() ([]models.CategoryProduct, error) {

	var categories []models.CategoryProduct

	err := r.db.
		Preload("Products").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// Mengambil kategori berdasarkan ID
func (r *CategoryProductRepository) FindByID(id uint) (*models.CategoryProduct, error) {

	var category models.CategoryProduct

	err := r.db.
		Preload("Products").
		First(&category, id).Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// Mengambil kategori berdasarkan slug
func (r *CategoryProductRepository) FindBySlug(slug string) (*models.CategoryProduct, error) {

	var category models.CategoryProduct

	err := r.db.
		Where("slug = ?", slug).
		First(&category).Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// Menambahkan kategori
func (r *CategoryProductRepository) Create(category *models.CategoryProduct) error {
	return r.db.Create(category).Error
}

// Mengubah kategori
func (r *CategoryProductRepository) Update(category *models.CategoryProduct) error {
	return r.db.Save(category).Error
}

// Menghapus kategori
func (r *CategoryProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.CategoryProduct{}, id).Error
}
