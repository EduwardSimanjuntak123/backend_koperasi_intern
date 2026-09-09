package repositories

import (
	"backend_koperasi/internal/models"
	"time"

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
			return db.Select("id, name, username, email, no_hp, role_id, building_id, floor_id")
		}).
		Preload("User.Building").
		Preload("User.Floor").
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
			return db.Select("id, name, username, email, no_hp, role_id, building_id, floor_id")
		}).
		Preload("User.Building").
		Preload("User.Floor").
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
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, username, email, no_hp, role_id, building_id, floor_id")
		}).
		Preload("User.Building").
		Preload("User.Floor").
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
			return db.Select("id, name, username, email, no_hp, role_id, building_id, floor_id")
		}).
		Preload("User.Building").
		Preload("User.Floor").
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

func (r *OrderRepository) FindByStatusPaginated(status models.OrderStatus, page, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	query := r.db.Model(&models.Order{}).Where("status = ?", status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	err := query.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, username, email, no_hp, role_id, building_id, floor_id")
		}).
		Preload("User.Building").
		Preload("User.Floor").
		Preload("Courier").
		Preload("Payment").
		Preload("Items").
		Preload("Items.Product").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *OrderRepository) CountStatuses(statuses []models.OrderStatus) (map[models.OrderStatus]int64, error) {
	counts := make(map[models.OrderStatus]int64, len(statuses))
	for _, status := range statuses {
		var count int64
		if err := r.db.Model(&models.Order{}).Where("status = ?", status).Count(&count).Error; err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, nil
}

func (r *OrderRepository) CountCompletedToday(start, end time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).
		Where("status = ? AND updated_at >= ? AND updated_at < ?", models.OrderCompleted, start, end).
		Count(&count).Error
	return count, err
}

func (r *OrderRepository) SumRevenueToday(start, end time.Time) (float64, error) {
	var revenue float64
	err := r.db.Model(&models.Order{}).
		Where("paid_at >= ? AND paid_at < ? AND status NOT IN ?", start, end, []models.OrderStatus{models.OrderCancelled}).
		Select("COALESCE(SUM(grand_total), 0)").Scan(&revenue).Error
	return revenue, err
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
