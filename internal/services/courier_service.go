package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type CourierService struct {
	courierRepo *repositories.CourierRepository
}

func NewCourierService(repo *repositories.CourierRepository) *CourierService {
	return &CourierService{
		courierRepo: repo,
	}
}

// GetAll Couriers
func (s *CourierService) GetAll() ([]models.Courier, error) {
	return s.courierRepo.FindAll()
}

// GetByID Courier
func (s *CourierService) GetByID(id string) (*models.Courier, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("invalid courier id")
	}

	return s.courierRepo.FindByID(id)
}

// Create Courier
func (s *CourierService) Create(courier *models.Courier) error {
	if strings.TrimSpace(courier.Name) == "" {
		return errors.New("courier name is required")
	}
	if courier.IsActive == nil {
		return errors.New("courier active status is required")
	}
	if strings.TrimSpace(courier.PhoneNumber) == "" {
		return errors.New("courier phone number is required")
	}

	id, err := s.courierRepo.GenerateNextID()
	if err != nil {
		return err
	}
	courier.ID = id

	return s.courierRepo.Create(courier)
}

// Update Courier
func (s *CourierService) Update(id string, courier *models.Courier) error {
	existing, err := s.courierRepo.FindByID(id)
	if err != nil {
		return err
	}

	if strings.TrimSpace(courier.Name) == "" {
		return errors.New("courier name is required")
	}
	if strings.TrimSpace(courier.PhoneNumber) == "" {
		return errors.New("courier phone number is required")
	}

	existing.Name = courier.Name
	existing.PhoneNumber = courier.PhoneNumber
	if courier.IsActive != nil {
		existing.IsActive = courier.IsActive
	}

	return s.courierRepo.Update(existing)
}

// Delete Courier
func (s *CourierService) Delete(id string) error {
	_, err := s.courierRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.courierRepo.Delete(id)
}
