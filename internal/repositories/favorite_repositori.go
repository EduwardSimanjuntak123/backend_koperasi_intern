package repositories

import (
	"errors"
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{
		db: db,
	}
}

// ==============================================================
// Mengambil wadah favorit milik user beserta daftar produknya
// ==============================================================
func (r *FavoriteRepository) FindByUserID(userID uint) (*models.Favorite, error) {
	var favorite models.Favorite

	// Menggunakan Preload("Products") untuk mengambil isi produk sekaligus
	err := r.db.Preload("Products").Where("user_id = ?", userID).First(&favorite).Error
	if err != nil {
		return nil, err
	}

	return &favorite, nil
}

// ==============================================================
// Mengecek apakah sebuah produk sudah ada di favorit user
// ==============================================================
func (r *FavoriteRepository) CheckIfExists(userID uint, productID uint) (bool, error) {
	var favorite models.Favorite

	// Cari favorit user dulu
	err := r.db.Where("user_id = ?", userID).First(&favorite).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // Jika wadah favorit belum ada, berarti produk juga belum ada
		}
		return false, err
	}

	// Cari apakah di dalam wadah favorit tersebut ada productID yang dimaksud
	var products []models.Product
	err = r.db.Model(&favorite).Where("id = ?", productID).Association("Products").Find(&products)

	return len(products) > 0, err
}

// ==============================================================
// Menambahkan produk ke favorit user (relasi Many-to-Many)
// ==============================================================
func (r *FavoriteRepository) AddProduct(userID uint, productID uint) error {
	var favorite models.Favorite

	// FirstOrCreate: Cari wadah favoritnya, jika belum ada, buatkan baru untuk user ini
	err := r.db.FirstOrCreate(&favorite, models.Favorite{UserID: userID}).Error
	if err != nil {
		return err
	}

	// Masukkan productID ke dalam tabel perantara (pivot table)
	product := models.Product{ID: productID}
	return r.db.Model(&favorite).Association("Products").Append(&product)
}

// ==============================================================
// Menghapus produk dari favorit user
// ==============================================================
func (r *FavoriteRepository) RemoveProduct(userID uint, productID uint) error {
	var favorite models.Favorite

	err := r.db.Where("user_id = ?", userID).First(&favorite).Error
	if err != nil {
		return err // Error atau tidak ditemukan
	}

	// Hapus productID dari tabel perantara relasi
	product := models.Product{ID: productID}
	return r.db.Model(&favorite).Association("Products").Delete(&product)
}