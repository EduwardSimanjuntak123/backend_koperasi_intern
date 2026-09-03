package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type BuildingRepository struct {
	db *gorm.DB
}

func NewBuildingRepository(db *gorm.DB) *BuildingRepository {
	return &BuildingRepository{
		db: db,
	}
}

// ======================================
// Mengambil seluruh lantai
// ======================================

func (r *BuildingRepository) FindAll() ([]models.Building, error) {

	var buildings []models.Building

	err := r.db.
		Find(&buildings).Error

	if err != nil {
		return nil, err
	}

	return buildings, nil
}

// ======================================
// Mengambil lantai berdasarkan ID
// ======================================

func (r *BuildingRepository) FindByID(id string) (*models.Building, error) {

	var building models.Building

	err := r.db.
		Where("id = ?", id).
		Preload("Floors").
		First(&building).Error

	if err != nil {
		return nil, err
	}

	return &building, nil
}

// ======================================
// Mengambil seluruh lantai berdasarkan gedung
// ======================================

func (r *BuildingRepository) FindByBuildingID(buildingID string) ([]models.Building, error) {

	var buildings []models.Building

	err := r.db.
		Where("building_id = ?", buildingID).
		Order("name ASC").
		Find(&buildings).Error

	if err != nil {
		return nil, err
	}

	return buildings, nil
}

// ======================================
// Mencari lantai berdasarkan nama & gedung
// Digunakan saat CREATE
// ======================================

func (r *BuildingRepository) FindByNameAndBuilding(name, buildingID string) (*models.Building, error) {

	var building models.Building

	err := r.db.
		Where("LOWER(name) = LOWER(?) AND building_id = ?", name, buildingID).
		First(&building).Error

	if err != nil {
		return nil, err
	}

	return &building, nil
}

// ======================================
// Mencari lantai berdasarkan nama & gedung
// kecuali ID tertentu (untuk UPDATE)
// ======================================

func (r *BuildingRepository) FindByNameAndBuildingExceptID(
	name,
	buildingID,
	id string,
) (*models.Building, error) {

	var building models.Building

	err := r.db.
		Where(
			"LOWER(name) = LOWER(?) AND building_id = ? AND id <> ?",
			name,
			buildingID,
			id,
		).
		First(&building).Error

	if err != nil {
		return nil, err
	}

	return &building, nil
}

// ======================================
// Menambahkan lantai
// ======================================

func (r *BuildingRepository) Create(building *models.Building) error {
	return r.db.Create(building).Error
}

// ======================================
// Mengubah lantai
// ======================================

func (r *BuildingRepository) Update(building *models.Building) error {
	return r.db.Save(building).Error
}

// ======================================
// Menghapus lantai
// ======================================

func (r *BuildingRepository) Delete(id string) error {

	result := r.db.
		Where("id = ?", id).
		Delete(&models.Building{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// ======================================
// Generate ID gedung
// ======================================

func (r *BuildingRepository) GenerateNextID() (string, error) {
	return generateNextPrefixedID(r.db, &models.Building{}, "GD")
}
