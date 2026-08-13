package repositories

import (
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

// Mengambil semua data favorit milik satu user, beserta detail produknya
func (r *FavoriteRepository) FindByUserID(userID uint) ([]models.Favorite, error) {
	var favorites []models.Favorite

	err := r.db.Preload("Product").Where("user_id = ?", userID).Find(&favorites).Error
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

// Mengecek apakah produk sudah ada di favorit user
func (r *FavoriteRepository) CheckIfExists(userID uint, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Favorite{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error

	return count > 0, err
}

// Menambahkan data favorit baru
func (r *FavoriteRepository) Create(favorite *models.Favorite) error {
	return r.db.Create(favorite).Error
}

// Menghapus data favorit berdasarkan user_id dan product_id
func (r *FavoriteRepository) DeleteByUserAndProduct(userID uint, productID uint) error {
	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.Favorite{}).Error
}