package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"

	"gorm.io/gorm"
)

type CategoryProductService struct {
	categoryProductRepo *repositories.CategoryProductRepository
}

func NewCategoryProductService(repo *repositories.CategoryProductRepository) *CategoryProductService {
	return &CategoryProductService{
		categoryProductRepo: repo,
	}
}

// =========================
// Get All Categories
// =========================
func (s *CategoryProductService) GetAll() ([]models.CategoryProduct, error) {
	return s.categoryProductRepo.FindAll()
}

// =========================
// Get Category By ID
// =========================
func (s *CategoryProductService) GetByID(id uint) (*models.CategoryProduct, error) {

	if id == 0 {
		return nil, errors.New("invalid category id")
	}

	return s.categoryProductRepo.FindByID(id)
}

// =========================
// Create Category
// =========================
func (s *CategoryProductService) Create(category *models.CategoryProduct) error {

	if strings.TrimSpace(category.Name) == "" {
		return errors.New("category name is required")
	}

	if strings.TrimSpace(category.Slug) == "" {
		return errors.New("category slug is required")
	}

	// Cek slug sudah ada atau belum
	_, err := s.categoryProductRepo.FindBySlug(category.Slug)

	if err == nil {
		return errors.New("category slug already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.categoryProductRepo.Create(category)
}

// =========================
// Update Category
// =========================
func (s *CategoryProductService) Update(id uint, category *models.CategoryProduct) error {

	existing, err := s.categoryProductRepo.FindByID(id)

	if err != nil {
		return err
	}

	if strings.TrimSpace(category.Name) == "" {
		return errors.New("category name is required")
	}

	if strings.TrimSpace(category.Slug) == "" {
		return errors.New("category slug is required")
	}

	// Cek slug jika berubah
	if existing.Slug != category.Slug {

		slug, err := s.categoryProductRepo.FindBySlug(category.Slug)

		if err == nil && slug.ID != existing.ID {
			return errors.New("category slug already exists")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	existing.Name = category.Name
	existing.Slug = category.Slug

	return s.categoryProductRepo.Update(existing)
}

// =========================
// Delete Category
// =========================
func (s *CategoryProductService) Delete(id uint) error {

	_, err := s.categoryProductRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.categoryProductRepo.Delete(id)
}
