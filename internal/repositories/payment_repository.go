package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) GetDB() *gorm.DB {
	return r.db
}

// ======================================
// Mengambil seluruh pembayaran
// ======================================

func (r *PaymentRepository) FindAll() ([]models.Payment, error) {
	var payments []models.Payment

	err := r.db.
		Preload("Order").
		Preload("Order.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, username, email, no_hp, role_id")
		}).
		Preload("Order.Items").
		Preload("Order.Items.Product").
		Order("created_at DESC").
		Find(&payments).Error

	if err != nil {
		return nil, err
	}

	return payments, nil
}

// ======================================
// Mengambil pembayaran berdasarkan ID
// ======================================

func (r *PaymentRepository) FindByID(id string) (*models.Payment, error) {
	var payment models.Payment

	err := r.db.
		Where("id = ?", id).
		Preload("Order").
		Preload("Order.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, username, email, no_hp, role_id")
		}).
		Preload("Order.Items").
		Preload("Order.Items.Product").
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// ======================================
// Mengambil pembayaran berdasarkan Order ID
// ======================================

func (r *PaymentRepository) FindByOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment

	err := r.db.
		Where("order_id = ?", orderID).
		Preload("Order").
		Preload("Order.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, username, email, no_hp, role_id")
		}).
		Preload("Order.Items").
		Preload("Order.Items.Product").
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// ======================================
// Menambahkan pembayaran
// ======================================

func (r *PaymentRepository) Create(tx *gorm.DB, payment *models.Payment) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(payment).Error
}

// ======================================
// Mengubah pembayaran
// ======================================

func (r *PaymentRepository) Update(tx *gorm.DB, payment *models.Payment) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Save(payment).Error
}

// ======================================
// Menghapus pembayaran
// ======================================

func (r *PaymentRepository) Delete(id string) error {
	result := r.db.
		Where("id = ?", id).
		Delete(&models.Payment{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// ======================================
// Generate ID Payment
// ======================================

func (r *PaymentRepository) GenerateNextID() (string, error) {
	return generateNextPrefixedID(r.db, &models.Payment{}, "PAY")
}
