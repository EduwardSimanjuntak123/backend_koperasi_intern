package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
	"backend_koperasi/internal/requests"

	"gorm.io/gorm"
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
	orderRepo   *repositories.OrderRepository
}

func NewPaymentService(
	paymentRepo *repositories.PaymentRepository,
	orderRepo *repositories.OrderRepository,
) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

// ======================================
// Get All Payments (Admin)
// ======================================

func (s *PaymentService) GetAll() ([]models.Payment, error) {
	return s.paymentRepo.FindAll()
}

// ======================================
// Get Payment By ID
// ======================================

func (s *PaymentService) GetByID(id string) (*models.Payment, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("ID pembayaran tidak boleh kosong")
	}

	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("data pembayaran tidak ditemukan")
		}
		return nil, err
	}

	return payment, nil
}

// ======================================
// Get Payment By Order ID
// ======================================

func (s *PaymentService) GetByOrderID(orderID string) (*models.Payment, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.New("ID pesanan tidak boleh kosong")
	}

	payment, err := s.paymentRepo.FindByOrderID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("data pembayaran untuk pesanan ini tidak ditemukan")
		}
		return nil, err
	}

	return payment, nil
}

// ======================================
// Create / Register Payment
// ======================================

func (s *PaymentService) Create(req requests.CreatePaymentRequest, userID string) (*models.Payment, error) {
	if strings.TrimSpace(req.OrderID) == "" {
		return nil, errors.New("ID pesanan wajib diisi")
	}

	if req.Amount <= 0 {
		return nil, errors.New("Jumlah pembayaran harus lebih dari 0")
	}

	db := s.paymentRepo.GetDB()
	var payment models.Payment

	err := db.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Where("id = ?", req.OrderID).First(&order).Error; err != nil {
			return errors.New("pesanan tidak ditemukan")
		}

		if order.Status == models.OrderPaid || order.Status == models.OrderCompleted {
			return errors.New("pesanan ini sudah dibayar/selesai")
		}

		if order.Status == models.OrderCancelled {
			return errors.New("tidak dapat melakukan pembayaran untuk pesanan yang telah dibatalkan")
		}

		// Cek apakah payment sudah ada
		var existingPayment models.Payment
		errFind := tx.Where("order_id = ?", req.OrderID).First(&existingPayment).Error

		if errFind == nil {
			// Update payment yang sudah ada
			existingPayment.Method = req.Method
			existingPayment.Amount = req.Amount
			existingPayment.ReferenceNo = req.ReferenceNo
			existingPayment.PaymentProof = req.PaymentProof

			if err := tx.Save(&existingPayment).Error; err != nil {
				return err
			}
			payment = existingPayment
			return nil
		}

		// Jika belum ada, buat baru
		payID, err := s.paymentRepo.GenerateNextID()
		if err != nil {
			return err
		}

		payment = models.Payment{
			ID:           payID,
			OrderID:      req.OrderID,
			Method:       req.Method,
			Status:       models.PaymentPending,
			Amount:       req.Amount,
			ReferenceNo:  req.ReferenceNo,
			PaymentProof: req.PaymentProof,
		}

		return tx.Create(&payment).Error
	})

	if err != nil {
		return nil, err
	}

	return s.paymentRepo.FindByID(payment.ID)
}

// ======================================
// Pay / Confirm Payment
// ======================================

func (s *PaymentService) Pay(req requests.PayOrderRequest, verifiedBy string) (*models.Payment, error) {
	if strings.TrimSpace(req.OrderID) == "" {
		return nil, errors.New("ID pesanan wajib diisi")
	}

	db := s.paymentRepo.GetDB()
	var updatedPayment models.Payment

	err := db.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Where("id = ?", req.OrderID).First(&order).Error; err != nil {
			return errors.New("pesanan tidak ditemukan")
		}

		if order.Status == models.OrderCancelled {
			return errors.New("pesanan sudah dibatalkan, tidak dapat dibayar")
		}

		now := time.Now()

		var payment models.Payment
		errFind := tx.Where("order_id = ?", req.OrderID).First(&payment).Error

		if errFind != nil {
			// Jika belum ada, buat entri payment
			payID, err := s.paymentRepo.GenerateNextID()
			if err != nil {
				return err
			}

			method := models.PaymentCash
			if req.Method != nil && *req.Method != "" {
				method = *req.Method
			}

			payment = models.Payment{
				ID:           payID,
				OrderID:      order.ID,
				Method:       method,
				Status:       models.PaymentPaid,
				Amount:       order.GrandTotal,
				PaidAt:       &now,
				ReferenceNo:  req.ReferenceNo,
				PaymentProof: req.PaymentProof,
			}

			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
		} else {
			// Update entri payment yang sudah ada
			if req.Method != nil && *req.Method != "" {
				payment.Method = *req.Method
			}
			if req.ReferenceNo != nil {
				payment.ReferenceNo = req.ReferenceNo
			}
			if req.PaymentProof != nil {
				payment.PaymentProof = req.PaymentProof
			}
			payment.Status = models.PaymentPaid
			payment.PaidAt = &now

			if err := tx.Save(&payment).Error; err != nil {
				return err
			}
		}

		// Update status Order menjadi PAID
		order.Status = models.OrderPaid
		order.PaidAt = &now
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		// Catat ke riwayat status Order
		oshID, err := s.orderRepo.GenerateStatusHistoryID()
		if err != nil {
			return err
		}

		desc := fmt.Sprintf("Pembayaran berhasil dikonfirmasi via %s", payment.Method)
		if payment.ReferenceNo != nil && *payment.ReferenceNo != "" {
			desc += fmt.Sprintf(" (Ref: %s)", *payment.ReferenceNo)
		}

		history := models.OrderStatusHistory{
			ID:          oshID,
			OrderID:     order.ID,
			Status:      string(models.OrderPaid),
			Description: desc,
			CreatedBy:   verifiedBy,
			CreatedAt:   time.Now(),
		}

		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		updatedPayment = payment
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.paymentRepo.FindByID(updatedPayment.ID)
}

// ======================================
// Update Payment (Admin)
// ======================================

func (s *PaymentService) Update(id string, req requests.UpdatePaymentRequest) (*models.Payment, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("ID pembayaran tidak valid")
	}

	db := s.paymentRepo.GetDB()
	var updatedPayment models.Payment

	err := db.Transaction(func(tx *gorm.DB) error {
		var oldPayment models.Payment
		if err := tx.Where("id = ?", id).First(&oldPayment).Error; err != nil {
			return errors.New("data pembayaran tidak ditemukan")
		}

		if req.Method != nil {
			oldPayment.Method = *req.Method
		}
		if req.Amount != nil {
			oldPayment.Amount = *req.Amount
		}
		if req.ReferenceNo != nil {
			oldPayment.ReferenceNo = req.ReferenceNo
		}
		if req.PaymentProof != nil {
			oldPayment.PaymentProof = req.PaymentProof
		}

		if req.Status != nil {
			oldPayment.Status = *req.Status

			if *req.Status == models.PaymentPaid && oldPayment.PaidAt == nil {
				now := time.Now()
				oldPayment.PaidAt = &now

				// Sinkronkan status order
				_ = tx.Model(&models.Order{}).
					Where("id = ?", oldPayment.OrderID).
					Updates(map[string]interface{}{
						"status":  models.OrderPaid,
						"paid_at": &now,
					}).Error
			}
		}

		if err := tx.Save(&oldPayment).Error; err != nil {
			return err
		}

		updatedPayment = oldPayment
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.paymentRepo.FindByID(updatedPayment.ID)
}

// ======================================
// Delete Payment (Admin)
// ======================================

func (s *PaymentService) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("ID pembayaran tidak valid")
	}

	return s.paymentRepo.Delete(id)
}
