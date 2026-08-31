package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type FloorRepository struct {
	db *gorm.DB
}

func NewFloorRepository(db *gorm.DB) *FloorRepository {
	return &FloorRepository{
		db: db,
	}
}

// ======================================
// Mengambil seluruh lantai
// ======================================

func (r *FloorRepository) FindAll() ([]models.Floor, error) {

	var floors []models.Floor

	err := r.db.
		Preload("Building").
		Find(&floors).Error

	if err != nil {
		return nil, err
	}

	return floors, nil
}

// ======================================
// Mengambil lantai berdasarkan ID
// ======================================

func (r *FloorRepository) FindByID(id string) (*models.Floor, error) {

	var floor models.Floor

	err := r.db.
		Where("id = ?", id).
		Preload("Building").
		First(&floor).Error

	if err != nil {
		return nil, err
	}

	return &floor, nil
}

// ======================================
// Mengambil seluruh lantai berdasarkan gedung
// ======================================

func (r *FloorRepository) FindByBuildingID(buildingID string) ([]models.Floor, error) {

	var floors []models.Floor

	err := r.db.
		Where("building_id = ?", buildingID).
		Order("name ASC").
		Find(&floors).Error

	if err != nil {
		return nil, err
	}

	return floors, nil
}

// ======================================
// Mencari lantai berdasarkan nama & gedung
// Digunakan saat CREATE
// ======================================

func (r *FloorRepository) FindByNameAndBuilding(name, buildingID string) (*models.Floor, error) {

	var floor models.Floor

	err := r.db.
		Where("LOWER(name) = LOWER(?) AND building_id = ?", name, buildingID).
		First(&floor).Error

	if err != nil {
		return nil, err
	}

	return &floor, nil
}

// ======================================
// Mencari lantai berdasarkan nama & gedung
// kecuali ID tertentu (untuk UPDATE)
// ======================================

func (r *FloorRepository) FindByNameAndBuildingExceptID(
	name,
	buildingID,
	id string,
) (*models.Floor, error) {

	var floor models.Floor

	err := r.db.
		Where(
			"LOWER(name) = LOWER(?) AND building_id = ? AND id <> ?",
			name,
			buildingID,
			id,
		).
		First(&floor).Error

	if err != nil {
		return nil, err
	}

	return &floor, nil
}

// ======================================
// Menambahkan lantai
// ======================================

func (r *FloorRepository) Create(floor *models.Floor) error {
	return r.db.Create(floor).Error
}

// ======================================
// Mengubah lantai
// ======================================

func (r *FloorRepository) Update(floor *models.Floor) error {
	return r.db.Save(floor).Error
}

// ======================================
// Menghapus lantai
// ======================================

func (r *FloorRepository) Delete(id string) error {

	result := r.db.
		Where("id = ?", id).
		Delete(&models.Floor{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// ======================================
// Generate ID lantai
// ======================================

func (r *FloorRepository) GenerateNextID() (string, error) {
	return generateNextPrefixedID(r.db, &models.Floor{}, "fl")
}
