package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

// ======================================
// Get All Orders
// ======================================

func (r *OrderRepository) FindAll() ([]models.Order, error) {

	var orders []models.Order

	err := r.db.
		Preload("User").
		Preload("Courier").
		Preload("Payment").
		Preload("Items").
		Preload("StatusHistory").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// ======================================
// Get Order By ID
// ======================================

func (r *OrderRepository) FindByID(id string) (*models.Order, error) {

	var order models.Order

	err := r.db.
		Preload("User").
		Preload("Courier").
		Preload("Payment").
		Preload("Items").
		Preload("StatusHistory").
		First(&order, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

// ======================================
// Get Order By Invoice Number
// ======================================

func (r *OrderRepository) FindByInvoiceNumber(invoice string) (*models.Order, error) {

	var order models.Order

	err := r.db.
		Preload("User").
		Preload("Courier").
		Preload("Payment").
		Preload("Items").
		Preload("StatusHistory").
		Where("invoice_number = ?", invoice).
		First(&order).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

// ======================================
// Get Orders By User ID
// ======================================

func (r *OrderRepository) FindByUserID(userID string) ([]models.Order, error) {

	var orders []models.Order

	err := r.db.
		Where("user_id = ?", userID).
		Preload("Courier").
		Preload("Payment").
		Preload("Items").
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// ======================================
// Create Order
// ======================================

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

// ======================================
// Update Order
// ======================================

func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

// ======================================
// Delete Order
// ======================================

func (r *OrderRepository) Delete(id string) error {
	return r.db.Delete(&models.Order{}, "id = ?", id).Error
}
