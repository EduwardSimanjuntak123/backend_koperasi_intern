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
func (r *CategoryProductRepository) FindByID(id string) (*models.CategoryProduct, error) {

	var category models.CategoryProduct

	err := r.db.
		Preload("Products").
		Where("id = ?", id).
		First(&category).Error

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
func (r *CategoryProductRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&models.CategoryProduct{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *CategoryProductRepository) GenerateNextID() (string, error) {
	return generateNextPrefixedID(r.db, &models.CategoryProduct{}, "C")
}
