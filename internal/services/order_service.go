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

type OrderService struct {
	orderRepo   *repositories.OrderRepository
	productRepo *repositories.ProductRepository
	cartRepo    *repositories.CartRepository
	paymentRepo *repositories.PaymentRepository
	courierRepo *repositories.CourierRepository
}

func NewOrderService(
	orderRepo *repositories.OrderRepository,
	productRepo *repositories.ProductRepository,
	cartRepo *repositories.CartRepository,
	paymentRepo *repositories.PaymentRepository,
	courierRepo *repositories.CourierRepository,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
		paymentRepo: paymentRepo,
		courierRepo: courierRepo,
	}
}

// ======================================
// Get All Orders (Admin)
// ======================================

func (s *OrderService) GetAll() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}

func (s *OrderService) GetByStatus(status string, page, limit int) ([]models.Order, int64, error) {
	parsedStatus := models.OrderStatus(strings.ToUpper(strings.TrimSpace(status)))
	switch parsedStatus {
	case models.OrderPending, models.OrderPaid, models.OrderPacking, models.OrderShipping, models.OrderDelivered, models.OrderCompleted, models.OrderCancelled:
	default:
		return nil, 0, errors.New("status pesanan tidak valid")
	}
	return s.orderRepo.FindByStatusPaginated(parsedStatus, page, limit)
}

type OrderDashboardSummary struct {
	PesananBaru    int64   `json:"pesanan_baru"`
	SedangDiproses int64   `json:"sedang_diproses"`
	SiapDikirim    int64   `json:"siap_dikirim"`
	SedangDiantar  int64   `json:"sedang_diantar"`
	SelesaiHariIni int64   `json:"selesai_hari_ini"`
	OmzetHariIni   float64 `json:"omzet_hari_ini"`
}

func (s *OrderService) GetDashboardSummary(now time.Time) (*OrderDashboardSummary, error) {
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)
	counts, err := s.orderRepo.CountStatuses([]models.OrderStatus{
		models.OrderPending,
		models.OrderPacking,
		models.OrderShipping,
		models.OrderDelivered,
	})
	if err != nil {
		return nil, err
	}
	completed, err := s.orderRepo.CountCompletedToday(start, end)
	if err != nil {
		return nil, err
	}
	revenue, err := s.orderRepo.SumRevenueToday(start, end)
	if err != nil {
		return nil, err
	}
	return &OrderDashboardSummary{
		PesananBaru:    counts[models.OrderPending],
		SedangDiproses: counts[models.OrderPacking],
		SiapDikirim:    counts[models.OrderShipping],
		SedangDiantar:  counts[models.OrderDelivered],
		SelesaiHariIni: completed,
		OmzetHariIni:   revenue,
	}, nil
}

func (s *OrderService) GetTodayRevenue(now time.Time) (float64, error) {
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return s.orderRepo.SumRevenueToday(start, start.Add(24*time.Hour))
}

// ======================================
// Get Order By ID
// ======================================

func (s *OrderService) GetByID(id string, userID string, isAdmin bool) (*models.Order, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("ID pesanan tidak boleh kosong")
	}

	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pesanan tidak ditemukan")
		}
		return nil, err
	}

	if !isAdmin && order.UserID != userID {
		return nil, errors.New("Anda tidak memiliki akses ke pesanan ini")
	}

	return order, nil
}

// ======================================
// Get Orders By User (Buyer)
// ======================================

func (s *OrderService) GetByUserID(userID string) ([]models.Order, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("User ID tidak valid")
	}

	return s.orderRepo.FindByUserID(userID)
}

// ======================================
// Create Order (Direct Checkout)
// ======================================

func (s *OrderService) Create(userID string, req requests.CreateOrderRequest) (*models.Order, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("User wajib login untuk membuat pesanan")
	}

	if len(req.Items) == 0 {
		return nil, errors.New("Daftar barang pesanan (items) tidak boleh kosong")
	}

	db := s.orderRepo.GetDB()
	var createdOrder models.Order

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Generate Order ID & Invoice Number
		orderID, err := s.orderRepo.GenerateNextID()
		if err != nil {
			return fmt.Errorf("gagal membuat ID pesanan: %w", err)
		}

		invoiceNo, err := s.orderRepo.GenerateInvoiceNumber()
		if err != nil {
			return fmt.Errorf("gagal membuat nomor invoice: %w", err)
		}

		// 2. Validasi Kurir jika dipilih
		if req.CourierID != nil && strings.TrimSpace(*req.CourierID) != "" {
			var courier models.Courier
			if err := tx.Where("id = ?", *req.CourierID).First(&courier).Error; err != nil {
				return errors.New("Kurir yang dipilih tidak ditemukan")
			}
		}

		var (
			orderItems    []models.OrderItem
			totalSubtotal float64
		)

		// 3. Proses setiap item: cek stok, kurangi stok, buat OrderItem
		for _, itemReq := range req.Items {
			if itemReq.Qty <= 0 {
				return errors.New("Jumlah barang (qty) harus lebih dari 0")
			}

			var product models.Product
			if err := tx.Where("id = ?", itemReq.ProductID).First(&product).Error; err != nil {
				return fmt.Errorf("produk dengan ID '%s' tidak ditemukan", itemReq.ProductID)
			}

			if product.Stock < itemReq.Qty {
				return fmt.Errorf("stok produk '%s' tidak mencukupi (tersedia: %d, diminta: %d)", product.Name, product.Stock, itemReq.Qty)
			}

			// Tentukan harga
			price := product.Price
			if product.PromotionPrice != nil && *product.PromotionPrice > 0 {
				price = *product.PromotionPrice
			}

			discount := 0.0
			subtotalItem := (price - discount) * float64(itemReq.Qty)
			totalSubtotal += subtotalItem

			// Generate Item ID
			oitID, err := s.orderRepo.GenerateOrderItemID()
			if err != nil {
				return fmt.Errorf("gagal membuat ID item pesanan: %w", err)
			}

			notes := ""
			if itemReq.Notes != nil {
				notes = *itemReq.Notes
			}

			orderItems = append(orderItems, models.OrderItem{
				ID:        oitID,
				OrderID:   orderID,
				ProductID: product.ID,
				Qty:       itemReq.Qty,
				Price:     price,
				Discount:  discount,
				Subtotal:  subtotalItem,
				Notes:     notes,
			})
		}

		// 4. Kalkulasi Grand Total
		grandTotal := totalSubtotal - req.DiscountAmount + req.ShippingCost
		if grandTotal < 0 {
			grandTotal = 0
		}

		// 5. Buat entri Order
		order := models.Order{
			ID:             orderID,
			UserID:         userID,
			CourierID:      req.CourierID,
			InvoiceNumber:  invoiceNo,
			Status:         models.OrderPending,
			Notes:          req.Notes,
			Subtotal:       totalSubtotal,
			DiscountAmount: req.DiscountAmount,
			ShippingCost:   req.ShippingCost,
			GrandTotal:     grandTotal,
			Items:          orderItems,
		}

		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("gagal menyimpan pesanan: %w", err)
		}

		// 6. Buat initial OrderStatusHistory
		oshID, err := s.orderRepo.GenerateStatusHistoryID()
		if err != nil {
			return fmt.Errorf("gagal membuat ID riwayat status: %w", err)
		}

		history := models.OrderStatusHistory{
			ID:          oshID,
			OrderID:     orderID,
			Status:      string(models.OrderPending),
			Description: "Pesanan berhasil dibuat, menunggu pembayaran.",
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("gagal mencatat riwayat status: %w", err)
		}

		// 7. Buat initial Payment jika payment method ditentukan
		if req.PaymentMethod != nil && *req.PaymentMethod != "" {
			payID, err := s.paymentRepo.GenerateNextID()
			if err != nil {
				return fmt.Errorf("gagal membuat ID pembayaran: %w", err)
			}

			payment := models.Payment{
				ID:      payID,
				OrderID: orderID,
				Method:  *req.PaymentMethod,
				Status:  models.PaymentPending,
				Amount:  grandTotal,
			}

			if err := tx.Create(&payment).Error; err != nil {
				return fmt.Errorf("gagal membuat entri pembayaran: %w", err)
			}
		}

		createdOrder = order
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.orderRepo.FindByID(createdOrder.ID)
}

// ======================================
// Checkout Cart (Checkout Keranjang Belanja)
// ======================================

func (s *OrderService) CheckoutCart(userID string, req requests.CheckoutCartRequest) (*models.Order, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("User wajib login untuk checkout keranjang")
	}

	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil || cart == nil || len(cart.CartItems) == 0 {
		return nil, errors.New("Keranjang belanja kosong atau tidak ditemukan")
	}

	var items []requests.OrderItemRequest
	for _, ci := range cart.CartItems {
		items = append(items, requests.OrderItemRequest{
			ProductID: ci.ProductID,
			Qty:       ci.Quantity,
		})
	}

	createReq := requests.CreateOrderRequest{
		Items:          items,
		CourierID:      req.CourierID,
		PaymentMethod:  req.PaymentMethod,
		ShippingCost:   req.ShippingCost,
		DiscountAmount: req.DiscountAmount,
		Notes:          req.Notes,
	}

	order, err := s.Create(userID, createReq)
	if err != nil {
		return nil, err
	}

	// Bersihkan keranjang belanja setelah checkout sukses
	_ = s.cartRepo.ClearCartItems(cart.ID)

	return order, nil
}

// ======================================
// Update Status Order (Admin / Sistem)
// ======================================

func (s *OrderService) UpdateStatus(orderID string, status models.OrderStatus, notes string, updatedBy string) error {
	if strings.TrimSpace(orderID) == "" {
		return errors.New("ID pesanan tidak valid")
	}

	db := s.orderRepo.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Where("id = ?", orderID).Preload("Items").First(&order).Error; err != nil {
			return errors.New("pesanan tidak ditemukan")
		}

		now := time.Now()
		order.Status = status

		switch status {
		case models.OrderPaid:
			order.PaidAt = &now
			// Update payment juga jika ada
			_ = tx.Model(&models.Payment{}).
				Where("order_id = ?", orderID).
				Updates(map[string]interface{}{
					"status":  models.PaymentPaid,
					"paid_at": &now,
				}).Error

		case models.OrderDelivered:
			order.DeliveredAt = &now

		case models.OrderCancelled:
			order.CancelledAt = &now
			// Order pending belum pernah mengurangi stok.
			if order.Status != models.OrderPending {
				for _, item := range order.Items {
					if err := tx.Model(&models.Product{}).
						Where("id = ?", item.ProductID).
						Update("stock", gorm.Expr("stock + ?", item.Qty)).Error; err != nil {
						return fmt.Errorf("gagal mengembalikan stok produk: %w", err)
					}
				}
			}
			// Update payment menjadi REFUNDED / FAILED jika ada
			_ = tx.Model(&models.Payment{}).
				Where("order_id = ?", orderID).
				Update("status", models.PaymentRefunded).Error
		}

		if err := tx.Save(&order).Error; err != nil {
			return fmt.Errorf("gagal memperbarui pesanan: %w", err)
		}

		// Catat riwayat status
		oshID, err := s.orderRepo.GenerateStatusHistoryID()
		if err != nil {
			return err
		}

		desc := fmt.Sprintf("Status pesanan diperbarui menjadi %s", status)
		if notes != "" {
			desc += ": " + notes
		}

		history := models.OrderStatusHistory{
			ID:          oshID,
			OrderID:     orderID,
			Status:      string(status),
			Description: desc,
			CreatedBy:   updatedBy,
			CreatedAt:   time.Now(),
		}

		return tx.Create(&history).Error
	})
}

// ======================================
// Assign Courier (Admin)
// ======================================

func (s *OrderService) AssignCourier(orderID, courierID, assignedBy string) error {
	if strings.TrimSpace(orderID) == "" {
		return errors.New("ID pesanan tidak valid")
	}

	if strings.TrimSpace(courierID) == "" {
		return errors.New("Kurir harus dipilih")
	}

	db := s.orderRepo.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Where("id = ?", orderID).First(&order).Error; err != nil {
			return errors.New("pesanan tidak ditemukan")
		}

		var courier models.Courier
		if err := tx.Where("id = ?", courierID).First(&courier).Error; err != nil {
			return errors.New("kurir tidak ditemukan")
		}

		order.CourierID = &courierID
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		// Catat riwayat penugasan kurir
		oshID, err := s.orderRepo.GenerateStatusHistoryID()
		if err != nil {
			return err
		}

		desc := fmt.Sprintf("Kurir '%s' (%s) ditugaskan untuk pengiriman", courier.Name, courier.PhoneNumber)
		history := models.OrderStatusHistory{
			ID:          oshID,
			OrderID:     orderID,
			Status:      string(order.Status),
			Description: desc,
			CreatedBy:   assignedBy,
			CreatedAt:   time.Now(),
		}

		return tx.Create(&history).Error
	})
}

// ======================================
// Cancel Order (Buyer / Admin)
// ======================================

func (s *OrderService) Cancel(orderID string, userID string, reason string, isAdmin bool) error {
	if strings.TrimSpace(orderID) == "" {
		return errors.New("ID pesanan tidak valid")
	}

	db := s.orderRepo.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Where("id = ?", orderID).Preload("Items").First(&order).Error; err != nil {
			return errors.New("pesanan tidak ditemukan")
		}

		if !isAdmin {
			if order.UserID != userID {
				return errors.New("Anda tidak memiliki izin untuk membatalkan pesanan ini")
			}
			if order.Status != models.OrderPending {
				return errors.New("Hanya pesanan dengan status PENDING yang dapat dibatalkan oleh pembeli")
			}
		}

		if order.Status == models.OrderCancelled {
			return errors.New("Pesanan sudah dibatalkan sebelumnya")
		}

		now := time.Now()
		order.Status = models.OrderCancelled
		order.CancelledAt = &now

		// Stok hanya dikembalikan jika sebelumnya sudah dibayar.
		if order.Status != models.OrderPending {
			for _, item := range order.Items {
				if err := tx.Model(&models.Product{}).
					Where("id = ?", item.ProductID).
					Update("stock", gorm.Expr("stock + ?", item.Qty)).Error; err != nil {
					return fmt.Errorf("gagal mengembalikan stok produk: %w", err)
				}
			}
		}

		// Update payment jika ada
		_ = tx.Model(&models.Payment{}).
			Where("order_id = ?", orderID).
			Update("status", models.PaymentRefunded).Error

		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		// Catat ke riwayat
		oshID, err := s.orderRepo.GenerateStatusHistoryID()
		if err != nil {
			return err
		}

		desc := "Pesanan dibatalkan"
		if reason != "" {
			desc += ". Alasan: " + reason
		}

		history := models.OrderStatusHistory{
			ID:          oshID,
			OrderID:     orderID,
			Status:      string(models.OrderCancelled),
			Description: desc,
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
		}

		return tx.Create(&history).Error
	})
}

// ======================================
// Delete Order (Admin)
// ======================================

func (s *OrderService) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("ID pesanan tidak valid")
	}

	return s.orderRepo.Delete(id)
}
