package services

import (
	"errors"

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
func (s *CartService) GetCartByUserID(userID uint) (*models.Cart, error) {
	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil {
		// Jika tidak ditemukan, kita buatkan keranjang kosong baru
		newCart := models.Cart{UserID: userID}
		if errCreate := s.cartRepo.CreateCart(&newCart); errCreate != nil {
			return nil, errors.New("failed to create cart")
		}
		return &newCart, nil
	}
	return cart, nil
}

// Menambahkan produk ke keranjang
func (s *CartService) AddToCart(userID uint, productID uint, quantity int) error {
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
	newItem := models.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return s.cartRepo.CreateCartItem(&newItem)
}

// Mengubah kuantitas item spesifik
func (s *CartService) UpdateItemQuantity(userID uint, itemID uint, quantity int) error {
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
func (s *CartService) RemoveItem(userID uint, itemID uint) error {
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