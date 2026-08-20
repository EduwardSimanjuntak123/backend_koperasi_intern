package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
	"backend_koperasi/internal/utils"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: repo,
	}
}

func (s *ProductService) GetAll(search string) ([]models.Product, error) {
	return s.productRepo.FindAll(search)
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {

	if id == 0 {
		return nil, errors.New("invalid product id")
	}

	return s.productRepo.FindByID(id)
}
func IsValidBadge(badge *models.ProductBadge) bool {
	if badge == nil {
		return true // badge boleh kosong
	}

	switch *badge {
	case models.BadgeNew,
		models.BadgeBestSeller:
		return true
	default:
		return false
	}
}

func IsValidInventoryMovement(movement *models.InventoryMovement) bool {
	if movement == nil {
		return true // movement boleh kosong
	}
	switch *movement {
	case models.FastMoving,
		models.SlowMoving:
		return true
	default:
		return false
	}
}

func (s *ProductService) Create(product *models.Product) error {

	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name is required")
	}
	// product.Slug = utils.GenerateSlug(product.Name)
	slug := utils.GenerateSlug(product.Name)

	_, err := s.productRepo.FindBySlug(slug)
	if err == nil {
		slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())
	}

	product.Slug = slug
	if product.Price <= 0 {
		return errors.New("price must be greater than zero")
	}

	if product.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	if product.BrandID == nil {
		return errors.New("brand id is required")
	}
	if product.UnitID == nil {
		return errors.New("unit id is required")
	}
	if product.CategoryID == nil {
		return errors.New("category id is required")
	}
	if product.Badge != nil {
		switch *product.Badge {
		case models.BadgeNew, models.BadgeBestSeller:
		default:
			return errors.New("badge tidak sesuai dengan enum yang tersedia")
		}
	}
	if product.InventoryMovement != nil {
		switch *product.InventoryMovement {
		case models.FastMoving, models.SlowMoving:
		default:
			return errors.New("inventory movement tidak sesuai dengan enum yang tersedia")
		}
	}

	return s.productRepo.Create(product)
}

func (s *ProductService) Update(id uint, product *models.Product) error {

	existing, err := s.productRepo.FindByID(id)

	if err != nil {
		return err
	}

	existing.Name = product.Name
	existing.Slug = product.Slug
	existing.Image = product.Image
	existing.Price = product.Price
	existing.Stock = product.Stock
	existing.Badge = product.Badge

	return s.productRepo.Update(existing)
}

func (s *ProductService) Delete(id uint) error {

	_, err := s.productRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.productRepo.Delete(id)
}
