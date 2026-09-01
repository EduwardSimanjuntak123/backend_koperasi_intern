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

func (r *OrderRepository) GetDB() *gorm.DB {
	return r.db
}

// ======================================
// Get All Orders (Admin)
// ======================================

func (r *OrderRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order

	err := r.db.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, username, email, no_hp, role_id")
		}).
		Preload("Courier").
		Preload("Payment").
		Preload("Items").
		Preload("Items.Product").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Order("created_at DESC").
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
		Where("id = ?", id).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, username, email, no_hp, role_id")
		}).
		Preload("Courier").
		Preload("Payment").
		Preload("Items").
		Preload("Items.Product").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		First(&order).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

// ======================================
// Get Orders By User (Buyer)
// ======================================

func (r *OrderRepository) FindByUserID(userID string) ([]models.Order, error) {
	var orders []models.Order

	err := r.db.
		Where("user_id = ?", userID).
		Preload("Courier").
		Preload("Payment").
		Preload("Items").
		Preload("Items.Product").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// ======================================
// Get Orders By Status
// ======================================

func (r *OrderRepository) FindByStatus(status models.OrderStatus) ([]models.Order, error) {
	var orders []models.Order

	err := r.db.
		Where("status = ?", status).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, username, email, no_hp, role_id")
		}).
		Preload("Courier").
		Preload("Payment").
		Preload("Items").
		Preload("Items.Product").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// ======================================
// Create Order (Supports DB Transaction)
// ======================================

func (r *OrderRepository) Create(tx *gorm.DB, order *models.Order) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(order).Error
}

// ======================================
// Update Order (Supports DB Transaction)
// ======================================

func (r *OrderRepository) Update(tx *gorm.DB, order *models.Order) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Save(order).Error
}

// ======================================
// Delete Order
// ======================================

func (r *OrderRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Hapus status history
		if err := tx.Where("order_id = ?", id).Delete(&models.OrderStatusHistory{}).Error; err != nil {
			return err
		}

		// Hapus payment
		if err := tx.Where("order_id = ?", id).Delete(&models.Payment{}).Error; err != nil {
			return err
		}

		// Hapus order items
		if err := tx.Where("order_id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
			return err
		}

		// Hapus order
		result := tx.Where("id = ?", id).Delete(&models.Order{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

// ======================================
// Update Status
// ======================================

func (r *OrderRepository) UpdateStatus(tx *gorm.DB, id string, status models.OrderStatus) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.
		Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// ======================================
// Assign Courier
// ======================================

func (r *OrderRepository) AssignCourier(tx *gorm.DB, orderID, courierID string) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("courier_id", courierID).Error
}

// ======================================
// Create Order Status History
// ======================================

func (r *OrderRepository) CreateStatusHistory(tx *gorm.DB, history *models.OrderStatusHistory) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(history).Error
}

// ======================================
// ID & Invoice Generators
// ======================================

func (r *OrderRepository) GenerateNextID() (string, error) {
	return generateNextPrefixedID(r.db, &models.Order{}, "ORD")
}

func (r *OrderRepository) GenerateInvoiceNumber() (string, error) {
	return generateInvoiceNumber(r.db)
}

func (r *OrderRepository) GenerateOrderItemID() (string, error) {
	return generateNextPrefixedID(r.db, &models.OrderItem{}, "OIT")
}

func (r *OrderRepository) GenerateStatusHistoryID() (string, error) {
	return generateNextPrefixedID(r.db, &models.OrderStatusHistory{}, "OSH")
}
