package repositories

import (
	"errors"

	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

// =====================================
// Operasi untuk model Cart (Keranjang Induk)
// =====================================

// Mencari keranjang berdasarkan UserID beserta seluruh isinya
func (r *CartRepository) FindCartByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart

	// Preload "CartItems" dan relasi bersarang "CartItems.Product" 
	// agar saat data keranjang ditarik, detail produknya ikut terbawa secara otomatis.
	err := r.db.Preload("CartItems").Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart not found")
		}
		return nil, err
	}

	return &cart, nil
}

// Membuat keranjang baru (biasanya dipanggil jika user baru pertama kali akses keranjang)
func (r *CartRepository) CreateCart(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

// =====================================
// Operasi untuk model CartItem (Isi Keranjang)
// =====================================

// Mencari spesifik satu item di dalam keranjang (untuk mengecek apakah barang sudah ada)
func (r *CartRepository) FindCartItem(cartID uint, productID uint) (*models.CartItem, error) {
	var item models.CartItem

	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart item not found")
		}
		return nil, err
	}

	return &item, nil
}

// Mencari item keranjang berdasarkan ID item-nya sendiri (untuk proses Update/Delete)
func (r *CartRepository) FindCartItemByID(itemID uint) (*models.CartItem, error) {
	var item models.CartItem

	err := r.db.First(&item, itemID).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Menambahkan item baru ke dalam keranjang
func (r *CartRepository) CreateCartItem(item *models.CartItem) error {
	return r.db.Create(item).Error
}

// Mengubah jumlah (quantity) item yang sudah ada di keranjang
func (r *CartRepository) UpdateCartItemQuantity(itemID uint, quantity int) error {
	return r.db.Model(&models.CartItem{}).Where("id = ?", itemID).Update("quantity", quantity).Error
}

// Menghapus item dari keranjang
func (r *CartRepository) DeleteCartItem(itemID uint) error {
	return r.db.Delete(&models.CartItem{}, itemID).Error
}