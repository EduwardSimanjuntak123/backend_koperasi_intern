package services

import (
	"errors"
	"net/http"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

// Sentinel errors untuk membedakan jenis kesalahan pada cart
var (
	ErrCartUserIDKosong     = errors.New("user ID tidak boleh kosong")
	ErrCartProductIDKosong  = errors.New("product ID tidak boleh kosong")
	ErrCartQuantityTidakValid = errors.New("jumlah barang harus lebih dari 0")
	ErrCartItemTidakDitemukan = errors.New("item keranjang tidak ditemukan")
	ErrCartTidakDitemukan   = errors.New("keranjang tidak ditemukan")
	ErrCartItemTidakDimiliki = errors.New("Anda tidak memiliki akses untuk mengubah item keranjang ini")
	ErrCartItemTidakBisaDihapus = errors.New("Anda tidak memiliki akses untuk menghapus item keranjang ini")
	ErrCartGagalDibuat      = errors.New("gagal membuat keranjang baru, silakan coba lagi")
	ErrCartGagalDikosongkan = errors.New("gagal mengosongkan keranjang, silakan coba lagi")
)

// ResolveCartErrorCode memetakan error ke HTTP status code yang sesuai
func ResolveCartErrorCode(err error) int {
	switch {
	case errors.Is(err, ErrCartUserIDKosong),
		errors.Is(err, ErrCartProductIDKosong),
		errors.Is(err, ErrCartQuantityTidakValid):
		return http.StatusBadRequest
	case errors.Is(err, ErrCartItemTidakDitemukan),
		errors.Is(err, ErrCartTidakDitemukan):
		return http.StatusNotFound
	case errors.Is(err, ErrCartItemTidakDimiliki),
		errors.Is(err, ErrCartItemTidakBisaDihapus):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

type CartService struct {
	cartRepo *repositories.CartRepository
}

func NewCartService(cartRepo *repositories.CartRepository) *CartService {
	return &CartService{
		cartRepo: cartRepo,
	}
}

// Mengambil keranjang user. Jika belum ada, buatkan keranjang baru secara otomatis (lazy creation).
func (s *CartService) GetCartByUserID(userID string) (*models.Cart, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrCartUserIDKosong
	}

	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil {
		// Jika belum ada, buat keranjang kosong baru secara otomatis
		cartID, errID := s.cartRepo.GenerateCartID()
		if errID != nil {
			return nil, ErrCartGagalDibuat
		}
		newCart := models.Cart{ID: cartID, UserID: userID}
		if errCreate := s.cartRepo.CreateCart(&newCart); errCreate != nil {
			return nil, ErrCartGagalDibuat
		}
		return &newCart, nil
	}
	return cart, nil
}

// Menambahkan produk ke keranjang.
// Jika produk sudah ada, kuantitasnya akan ditambahkan (bukan membuat baris baru).
// Jika keranjang belum ada, akan dibuat otomatis.
func (s *CartService) AddToCart(userID string, productID string, quantity int) error {
	if strings.TrimSpace(userID) == "" {
		return ErrCartUserIDKosong
	}
	if strings.TrimSpace(productID) == "" {
		return ErrCartProductIDKosong
	}
	if quantity <= 0 {
		return ErrCartQuantityTidakValid
	}

	// 1. Pastikan keranjang user tersedia (buat otomatis jika belum ada)
	cart, err := s.GetCartByUserID(userID)
	if err != nil {
		return err
	}

	// 2. Cek apakah produk sudah ada di keranjang
	existingItem, err := s.cartRepo.FindCartItem(cart.ID, productID)
	if err == nil && existingItem != nil {
		// Jika sudah ada, tambahkan kuantitasnya (tidak membuat baris baru)
		newQuantity := existingItem.Quantity + quantity
		return s.cartRepo.UpdateCartItemQuantity(existingItem.ID, newQuantity)
	}

	// 3. Jika belum ada, buat item baru di keranjang
	itemID, err := s.cartRepo.GenerateCartItemID()
	if err != nil {
		return errors.New("gagal membuat ID item keranjang, silakan coba lagi")
	}
	newItem := models.CartItem{
		ID:        itemID,
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return s.cartRepo.CreateCartItem(&newItem)
}

// Mengubah kuantitas item spesifik di keranjang
func (s *CartService) UpdateItemQuantity(userID string, itemID string, quantity int) error {
	// Verifikasi item tersebut benar-benar ada
	item, err := s.cartRepo.FindCartItemByID(itemID)
	if err != nil {
		return ErrCartItemTidakDitemukan
	}

	// Verifikasi kepemilikan: pastikan item ini ada di keranjang milik user yang bersangkutan
	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil {
		return ErrCartTidakDitemukan
	}
	if cart.ID != item.CartID {
		return ErrCartItemTidakDimiliki
	}

	return s.cartRepo.UpdateCartItemQuantity(itemID, quantity)
}

// Menghapus satu item dari keranjang
func (s *CartService) RemoveItem(userID string, itemID string) error {
	// Verifikasi item tersebut benar-benar ada
	item, err := s.cartRepo.FindCartItemByID(itemID)
	if err != nil {
		return ErrCartItemTidakDitemukan
	}

	// Verifikasi kepemilikan
	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil {
		return ErrCartTidakDitemukan
	}
	if cart.ID != item.CartID {
		return ErrCartItemTidakBisaDihapus
	}

	return s.cartRepo.DeleteCartItem(itemID)
}

// Mengosongkan seluruh isi keranjang milik user (clear cart)
func (s *CartService) ClearCart(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrCartUserIDKosong
	}

	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil {
		return ErrCartTidakDitemukan
	}

	if errClear := s.cartRepo.ClearCartItems(cart.ID); errClear != nil {
		return ErrCartGagalDikosongkan
	}
	return nil
}
