package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: repo,
	}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.productRepo.FindAll()
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {

	if id == 0 {
		return nil, errors.New("invalid product id")
	}

	return s.productRepo.FindByID(id)
}

func (s *ProductService) Create(product *models.Product) error {

	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name is required")
	}

	if strings.TrimSpace(product.Slug) == "" {
		return errors.New("product slug is required")
	}

	if product.Price <= 0 {
		return errors.New("price must be greater than zero")
	}

	if product.Stock < 0 {
		return errors.New("stock cannot be negative")
	}
	if product.OriginalPrice != nil && *product.OriginalPrice < 0 {
		return errors.New("original price must be greater than or equal to 0")
	}
	if product.WeightGram != nil && *product.WeightGram < 0 {
		return errors.New("weight cannot be negative")
	}
	if product.LengthCm != nil && *product.LengthCm < 0 {
		return errors.New("length cannot be negative")
	}
	if product.WidthCm != nil && *product.WidthCm < 0 {
		return errors.New("width cannot be negative")
	}
	if product.HeightCm != nil && *product.HeightCm < 0 {
		return errors.New("height cannot be negative")
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
	existing.OriginalPrice = product.OriginalPrice
	existing.Stock = product.Stock
	existing.Volume = product.Volume
	existing.WeightGram = product.WeightGram
	existing.LengthCm = product.LengthCm
	existing.WidthCm = product.WidthCm
	existing.HeightCm = product.HeightCm
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
