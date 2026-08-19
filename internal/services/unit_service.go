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
func (s *UnitService) GetByID(id uint) (*models.Unit, error) {

	if id == 0 {
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

	return s.unitRepo.Create(unit)
}

// =========================
// Update Unit
// =========================
func (s *UnitService) Update(id uint, unit *models.Unit) error {

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
func (s *UnitService) Delete(id uint) error {

	_, err := s.unitRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.unitRepo.Delete(id)
}
