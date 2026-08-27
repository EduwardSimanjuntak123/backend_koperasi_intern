package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type CartService struct {
	cartRepo *repositories.CartRepository
}

func NewCartService(cartRepo *repositories.CartRepository) *CartService {
	return &CartService{
		cartRepo: cartRepo,
	}
}

// Mengambil keranjang user. Jika belum ada, buatkan keranjang baru secara otomatis.
func (s *CartService) GetCartByUserID(userID string) (*models.Cart, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("invalid user id")
	}
	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil {
		// Jika tidak ditemukan, kita buatkan keranjang kosong baru
		cartID, errID := s.cartRepo.GenerateCartID()
		if errID != nil {
			return nil, errID
		}
		newCart := models.Cart{ID: cartID, UserID: userID}
		if errCreate := s.cartRepo.CreateCart(&newCart); errCreate != nil {
			return nil, errors.New("failed to create cart")
		}
		return &newCart, nil
	}
	return cart, nil
}

// Menambahkan produk ke keranjang
func (s *CartService) AddToCart(userID string, productID string, quantity int) error {

	if strings.TrimSpace(userID) == "" {
		return errors.New("invalid user id")
	}

	if strings.TrimSpace(productID) == "" {
		return errors.New("invalid product id")
	}

	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	// 1. Pastikan keranjang user tersedia
	cart, err := s.GetCartByUserID(userID)
	if err != nil {
		return err
	}

	// 2. Cek apakah barang sudah ada di keranjang
	existingItem, err := s.cartRepo.FindCartItem(cart.ID, productID)

	if err == nil && existingItem != nil {
		// Jika sudah ada, tambahkan kuantitasnya (tidak membuat baris baru)
		newQuantity := existingItem.Quantity + quantity
		return s.cartRepo.UpdateCartItemQuantity(existingItem.ID, newQuantity)
	}

	// 3. Jika belum ada, buat item baru di keranjang
	itemID, err := s.cartRepo.GenerateCartItemID()
	if err != nil {
		return err
	}
	newItem := models.CartItem{
		ID:        itemID,
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return s.cartRepo.CreateCartItem(&newItem)
}

// Mengubah kuantitas item spesifik
func (s *CartService) UpdateItemQuantity(userID string, itemID string, quantity int) error {
	// Verifikasi kepemilikan: Pastikan item ini benar-benar ada di keranjang milik userID tersebut
	item, err := s.cartRepo.FindCartItemByID(itemID)
	if err != nil {
		return errors.New("cart item not found")
	}

	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil || cart.ID != item.CartID {
		return errors.New("unauthorized to update this cart item")
	}

	return s.cartRepo.UpdateCartItemQuantity(itemID, quantity)
}

// Menghapus item dari keranjang
func (s *CartService) RemoveItem(userID string, itemID string) error {
	// Verifikasi kepemilikan seperti di atas
	item, err := s.cartRepo.FindCartItemByID(itemID)
	if err != nil {
		return errors.New("cart item not found")
	}

	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil || cart.ID != item.CartID {
		return errors.New("unauthorized to delete this cart item")
	}

	return s.cartRepo.DeleteCartItem(itemID)
}
