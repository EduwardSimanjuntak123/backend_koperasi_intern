package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"

	"gorm.io/gorm"
)

type FloorService struct {
	floorRepo *repositories.FloorRepository
}

func NewFloorService(repo *repositories.FloorRepository) *FloorService {
	return &FloorService{
		floorRepo: repo,
	}
}

// ======================================
// Get All Floors
// ======================================

func (s *FloorService) GetAllFloors() ([]models.Floor, error) {
	return s.floorRepo.FindAll()
}

// ======================================
// Get Floor By ID
// ======================================

func (s *FloorService) GetFloorByID(id string) (*models.Floor, error) {

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("floor id is required")
	}

	return s.floorRepo.FindByID(id)
}

// ======================================
// Get Floors By Building
// ======================================

func (s *FloorService) GetFloorsByBuilding(buildingID string) ([]models.Floor, error) {

	if strings.TrimSpace(buildingID) == "" {
		return nil, errors.New("building id is required")
	}

	return s.floorRepo.FindByBuildingID(buildingID)
}

// ======================================
// Create Floor
// ======================================

func (s *FloorService) CreateFloor(floor *models.Floor) error {

	if floor == nil {
		return errors.New("floor is required")
	}

	if strings.TrimSpace(floor.Name) == "" {
		return errors.New("floor name is required")
	}

	if strings.TrimSpace(floor.BuildingID) == "" {
		return errors.New("building id is required")
	}

	// Cek apakah nama lantai sudah ada pada gedung yang sama
	_, err := s.floorRepo.FindByNameAndBuilding(floor.Name, floor.BuildingID)
	if err == nil {
		return errors.New("floor already exists in this building")
	}

	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Generate ID
	id, err := s.floorRepo.GenerateNextID()
	if err != nil {
		return err
	}

	floor.ID = id

	return s.floorRepo.Create(floor)
}

// ======================================
// Update Floor
// ======================================

func (s *FloorService) UpdateFloor(floor *models.Floor) error {

	if floor == nil {
		return errors.New("floor is required")
	}

	if strings.TrimSpace(floor.ID) == "" {
		return errors.New("floor id is required")
	}

	if strings.TrimSpace(floor.Name) == "" {
		return errors.New("floor name is required")
	}

	if strings.TrimSpace(floor.BuildingID) == "" {
		return errors.New("building id is required")
	}

	return s.floorRepo.Update(floor)
}

// ======================================
// Delete Floor
// ======================================

func (s *FloorService) DeleteFloor(id string) error {

	if strings.TrimSpace(id) == "" {
		return errors.New("floor id is required")
	}

	return s.floorRepo.Delete(id)
}
