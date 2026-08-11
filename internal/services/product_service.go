package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type ProductService struct {
	repository *repositories.ProductRepository
}

func NewProductService(
	repository *repositories.ProductRepository,
) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) GetProducts() ([]models.Product, error) {
	return s.repository.FindAll()
}

func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	if id == 0 {
		return nil, errors.New("invalid product ID")
	}

	return s.repository.FindByID(id)
}

func (s *ProductService) CreateProduct(
	product *models.Product,
) error {

	// Validasi nama
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name is required")
	}

	// Validasi slug
	if strings.TrimSpace(product.Slug) == "" {
		return errors.New("product slug is required")
	}

	// Validasi harga
	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}

	// Validasi stock
	if product.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	return s.repository.Create(product)
}

func (s *ProductService) UpdateProduct(
	id uint,
	product *models.Product,
) error {

	if id == 0 {
		return errors.New("invalid product ID")
	}

	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name is required")
	}

	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}

	if product.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	existingProduct, err := s.repository.FindByID(id)

	if err != nil {
		return err
	}

	existingProduct.Name = product.Name
	existingProduct.Slug = product.Slug
	existingProduct.Image = product.Image
	existingProduct.Price = product.Price
	existingProduct.OriginalPrice = product.OriginalPrice
	existingProduct.Stock = product.Stock
	existingProduct.Volume = product.Volume
	existingProduct.WeightGram = product.WeightGram
	existingProduct.LengthCm = product.LengthCm
	existingProduct.WidthCm = product.WidthCm
	existingProduct.HeightCm = product.HeightCm
	existingProduct.Badge = product.Badge

	return s.repository.Update(existingProduct)
}

func (s *ProductService) DeleteProduct(id uint) error {

	if id == 0 {
		return errors.New("invalid product ID")
	}

	_, err := s.repository.FindByID(id)

	if err != nil {
		return err
	}

	return s.repository.Delete(id)
}
