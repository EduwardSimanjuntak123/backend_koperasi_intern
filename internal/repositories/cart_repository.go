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
func (r *CartRepository) FindCartByUserID(userID string) (*models.Cart, error) {
	var cart models.Cart

	// Preload "CartItems" dan relasi bersarang "CartItems.Product"
	// agar saat data keranjang ditarik, detail produknya ikut terbawa secara otomatis.
	err := r.db.Preload("CartItems").Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("keranjang tidak ditemukan")
		}
		return nil, err
	}

	return &cart, nil
}

// Membuat keranjang baru (dipanggil otomatis jika user belum pernah memiliki keranjang)
func (r *CartRepository) CreateCart(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

// =====================================
// Operasi untuk model CartItem (Isi Keranjang)
// =====================================

// Mencari item spesifik di dalam keranjang berdasarkan cart_id dan product_id
// (digunakan untuk mengecek apakah produk sudah ada sebelum menambahkan)
func (r *CartRepository) FindCartItem(cartID string, productID string) (*models.CartItem, error) {
	var item models.CartItem

	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk belum ada di keranjang")
		}
		return nil, err
	}

	return &item, nil
}

// Mencari item keranjang berdasarkan ID item-nya sendiri (untuk proses Update/Delete)
func (r *CartRepository) FindCartItemByID(itemID string) (*models.CartItem, error) {
	var item models.CartItem

	// Gunakan .Where() agar aman untuk tipe ID string (menghindari ambiguitas primary key)
	err := r.db.Where("id = ?", itemID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item keranjang tidak ditemukan")
		}
		return nil, err
	}

	return &item, nil
}

// Menambahkan item baru ke dalam keranjang
func (r *CartRepository) CreateCartItem(item *models.CartItem) error {
	return r.db.Create(item).Error
}

// Mengubah jumlah (quantity) item yang sudah ada di keranjang
func (r *CartRepository) UpdateCartItemQuantity(itemID string, quantity int) error {
	return r.db.Model(&models.CartItem{}).Where("id = ?", itemID).Update("quantity", quantity).Error
}

// Menghapus satu item dari keranjang
func (r *CartRepository) DeleteCartItem(itemID string) error {
	return r.db.Where("id = ?", itemID).Delete(&models.CartItem{}).Error
}

// Menghapus semua item dalam keranjang (clear cart)
func (r *CartRepository) ClearCartItems(cartID string) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

func (r *CartRepository) GenerateCartID() (string, error) {
	return generateNextPrefixedID(r.db, &models.Cart{}, "CRT")
}

func (r *CartRepository) GenerateCartItemID() (string, error) {
	return generateNextPrefixedID(r.db, &models.CartItem{}, "CI")
}
