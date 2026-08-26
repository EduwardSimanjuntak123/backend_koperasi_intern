package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type UnitService struct {
	unitRepo *repositories.UnitRepository
}

func NewUnitService(repo *repositories.UnitRepository) *UnitService {
	return &UnitService{
		unitRepo: repo,
	}
}

// =========================
// Get All Units
// =========================
func (s *UnitService) GetAll() ([]models.Unit, error) {
	return s.unitRepo.FindAll()
}

// =========================
// Get Unit By ID
// =========================
func (s *UnitService) GetByID(id string) (*models.Unit, error) {

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("invalid unit id")
	}

	return s.unitRepo.FindByID(id)
}

// =========================
// Create Unit
// =========================
func (s *UnitService) Create(unit *models.Unit) error {

	if strings.TrimSpace(unit.Name) == "" {
		return errors.New("unit name is required")
	}

	// Generate ID
	id, err := s.unitRepo.GenerateNextID()
	if err != nil {
		return err
	}
	unit.ID = id

	return s.unitRepo.Create(unit)
}

// =========================
// Update Unit
// =========================
func (s *UnitService) Update(id string, unit *models.Unit) error {

	existing, err := s.unitRepo.FindByID(id)

	if err != nil {
		return err
	}

	if strings.TrimSpace(unit.Name) == "" {
		return errors.New("unit name is required")
	}

	existing.Name = unit.Name

	return s.unitRepo.Update(existing)
}

// =========================
// Delete Unit
// =========================
func (s *UnitService) Delete(id string) error {

	_, err := s.unitRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.unitRepo.Delete(id)
}
