package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type BuildingService struct {
	buildingRepo *repositories.BuildingRepository
}

func NewBuildingService(repo *repositories.BuildingRepository) *BuildingService {
	return &BuildingService{
		buildingRepo: repo,
	}
}

// ======================================
// Get All Buildings
// ======================================

func (s *BuildingService) GetAllBuildings() ([]models.Building, error) {
	return s.buildingRepo.FindAll()
}

// ======================================
// Get Building By ID
// ======================================

func (s *BuildingService) GetBuildingByID(id string) (*models.Building, error) {

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("building id is required")
	}

	return s.buildingRepo.FindByID(id)
}

// ======================================
// Get Buildings By ID
// ======================================

func (s *BuildingService) GetBuildingsByID(buildingID string) ([]models.Building, error) {

	if strings.TrimSpace(buildingID) == "" {
		return nil, errors.New("building id is required")
	}

	return s.buildingRepo.FindByBuildingID(buildingID)
}

// ======================================
// Create Building
// ======================================

func (s *BuildingService) CreateBuilding(building *models.Building) error {

	if building == nil {
		return errors.New("building is required")
	}

	if strings.TrimSpace(building.Name) == "" {
		return errors.New("building name is required")
	}

	// Generate ID
	id, err := s.buildingRepo.GenerateNextID()
	if err != nil {
		return err
	}

	building.ID = id

	return s.buildingRepo.Create(building)
}

// ======================================
// Update Building
// ======================================

func (s *BuildingService) UpdateBuilding(building *models.Building) error {

	if building == nil {
		return errors.New("building is required")
	}

	if strings.TrimSpace(building.ID) == "" {
		return errors.New("building id is required")
	}

	if strings.TrimSpace(building.Name) == "" {
		return errors.New("building name is required")
	}

	return s.buildingRepo.Update(building)
}

// ======================================
// Delete Building
// ======================================

func (s *BuildingService) DeleteBuilding(id string) error {

	if strings.TrimSpace(id) == "" {
		return errors.New("building id is required")
	}

	return s.buildingRepo.Delete(id)
}
