package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type CourierRepository struct {
	db *gorm.DB
}

func NewCourierRepository(db *gorm.DB) *CourierRepository {
	return &CourierRepository{
		db: db,
	}
}

// Mengambil seluruh bangunan
func (r *CourierRepository) FindAll() ([]models.Courier, error) {

	var couriers []models.Courier

	err := r.db.
		Find(&couriers).Error

	if err != nil {
		return nil, err
	}

	return couriers, nil
}

// Mengambil bangunan berdasarkan ID
func (r *CourierRepository) FindByID(id string) (*models.Courier, error) {

	var courier models.Courier

	err := r.db.
		Where("id = ?", id).
		First(&courier).Error

	if err != nil {
		return nil, err
	}

	return &courier, nil
}

// Mengambil bangunan berdasarkan code
func (r *CourierRepository) FindByCode(code string) (*models.Courier, error) {

	var courier models.Courier

	err := r.db.
		Where("code = ?", code).
		First(&courier).Error

	if err != nil {
		return nil, err
	}

	return &courier, nil
}

// Menambahkan bangunan
func (r *CourierRepository) Create(courier *models.Courier) error {
	return r.db.Create(courier).Error
}

// Mengubah bangunan
func (r *CourierRepository) Update(courier *models.Courier) error {
	return r.db.Save(courier).Error
}

// Menghapus bangunan
func (r *CourierRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&models.Courier{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *CourierRepository) GenerateNextID() (string, error) {
	return generateNextPrefixedID(r.db, &models.Courier{}, "Co")
}
