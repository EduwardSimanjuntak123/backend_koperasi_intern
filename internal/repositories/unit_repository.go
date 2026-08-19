package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type UnitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) *UnitRepository {
	return &UnitRepository{
		db: db,
	}
}

// Mengambil seluruh kategori
func (r *UnitRepository) FindAll() ([]models.Unit, error) {

	var units []models.Unit

	err := r.db.
		Find(&units).Error

	if err != nil {
		return nil, err
	}

	return units, nil
}

// Mengambil kategori berdasarkan ID
func (r *UnitRepository) FindByID(id uint) (*models.Unit, error) {

	var unit models.Unit

	err := r.db.
		Find(&unit, id).Error

	if err != nil {
		return nil, err
	}

	return &unit, nil
}

// Mengambil kategori berdasarkan slug
func (r *UnitRepository) FindBySlug(slug string) (*models.Unit, error) {

	var unit models.Unit

	err := r.db.
		Where("slug = ?", slug).
		First(&unit).Error

	if err != nil {
		return nil, err
	}

	return &unit, nil
}

// Menambahkan kategori
func (r *UnitRepository) Create(unit *models.Unit) error {
	return r.db.Create(unit).Error
}

// Mengubah kategori
func (r *UnitRepository) Update(unit *models.Unit) error {
	return r.db.Save(unit).Error
}

// Menghapus kategori
func (r *UnitRepository) Delete(id uint) error {
	return r.db.Delete(&models.Unit{}, id).Error
}
