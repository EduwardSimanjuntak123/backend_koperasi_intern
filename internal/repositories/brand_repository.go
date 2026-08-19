package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type BrandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) *BrandRepository {
	return &BrandRepository{
		db: db,
	}
}

// Mengambil seluruh brand
func (r *BrandRepository) FindAll() ([]models.Brand, error) {

	var brands []models.Brand

	err := r.db.
		Find(&brands).Error

	if err != nil {
		return nil, err
	}

	return brands, nil
}

// Mengambil brand berdasarkan ID
func (r *BrandRepository) FindByID(id uint) (*models.Brand, error) {

	var brand models.Brand

	err := r.db.
		Find(&brand, id).Error

	if err != nil {
		return nil, err
	}

	return &brand, nil
}

// Mengambil brand berdasarkan slug
func (r *BrandRepository) FindBySlug(slug string) (*models.Brand, error) {

	var brand models.Brand

	err := r.db.
		Where("slug = ?", slug).
		First(&brand).Error

	if err != nil {
		return nil, err
	}

	return &brand, nil
}

// Menambahkan brand
func (r *BrandRepository) Create(brand *models.Brand) error {
	return r.db.Create(brand).Error
}

// Mengubah brand
func (r *BrandRepository) Update(brand *models.Brand) error {
	return r.db.Save(brand).Error
}

// Menghapus brand
func (r *BrandRepository) Delete(id uint) error {
	return r.db.Delete(&models.Brand{}, id).Error
}
