package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type OrderService struct {
	orderRepo *repositories.OrderRepository
}

func NewOrderService(repo *repositories.OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: repo,
	}
}

// ======================================
// Get All Orders
// ======================================

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}

// ======================================
// Get Order By ID
// ======================================

func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("invalid order id")
	}

	return s.orderRepo.FindByID(id)
}

// ======================================
// Get Order By Invoice Number
// ======================================

func (s *OrderService) GetOrderByInvoice(invoice string) (*models.Order, error) {

	if strings.TrimSpace(invoice) == "" {
		return nil, errors.New("invoice number is required")
	}

	return s.orderRepo.FindByInvoiceNumber(invoice)
}

// ======================================
// Get Orders By User
// ======================================

func (s *OrderService) GetOrdersByUser(userID string) ([]models.Order, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("invalid user id")
	}

	return s.orderRepo.FindByUserID(userID)
}

// ======================================
// Create Order
// ======================================

func (s *OrderService) CreateOrder(order *models.Order) error {

	if order == nil {
		return errors.New("order is required")
	}

	if strings.TrimSpace(order.UserID) == "" {
		return errors.New("user id is required")
	}

	if strings.TrimSpace(order.InvoiceNumber) == "" {
		return errors.New("invoice number is required")
	}

	return s.orderRepo.Create(order)
}

// ======================================
// Update Order
// ======================================

func (s *OrderService) UpdateOrder(order *models.Order) error {

	if order == nil {
		return errors.New("order is required")
	}

	if strings.TrimSpace(order.ID) == "" {
		return errors.New("order id is required")
	}

	return s.orderRepo.Update(order)
}

// ======================================
// Delete Order
// ======================================

func (s *OrderService) DeleteOrder(id string) error {

	if strings.TrimSpace(id) == "" {
		return errors.New("invalid order id")
	}

	return s.orderRepo.Delete(id)
}
