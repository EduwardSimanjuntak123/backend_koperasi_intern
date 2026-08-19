package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type BrandService struct {
	brandRepo *repositories.BrandRepository
}

func NewBrandService(repo *repositories.BrandRepository) *BrandService {
	return &BrandService{
		brandRepo: repo,
	}
}

// =========================
// Get All Brands
// =========================
func (s *BrandService) GetAll() ([]models.Brand, error) {
	return s.brandRepo.FindAll()
}

// =========================
// Get Brand By ID
// =========================
func (s *BrandService) GetByID(id uint) (*models.Brand, error) {

	if id == 0 {
		return nil, errors.New("invalid brand id")
	}

	return s.brandRepo.FindByID(id)
}

// =========================
// Create Brand
// =========================
func (s *BrandService) Create(brand *models.Brand) error {

	if strings.TrimSpace(brand.Name) == "" {
		return errors.New("brand name is required")
	}

	return s.brandRepo.Create(brand)
}

// =========================
// Update Brand
// =========================
func (s *BrandService) Update(id uint, brand *models.Brand) error {

	existing, err := s.brandRepo.FindByID(id)

	if err != nil {
		return err
	}

	if strings.TrimSpace(brand.Name) == "" {
		return errors.New("brand name is required")
	}

	existing.Name = brand.Name

	return s.brandRepo.Update(existing)
}

// =========================
// Delete Brand
// =========================
func (s *BrandService) Delete(id uint) error {

	_, err := s.brandRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.brandRepo.Delete(id)
}
