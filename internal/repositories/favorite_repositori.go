package repositories

import (
	"backend_koperasi/internal/models"
	"fmt"

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
func (r *FavoriteRepository) FindByUserID(userID string) ([]models.Favorite, error) {
	var favorites []models.Favorite

	err := r.db.Preload("Product").Where("user_id = ?", userID).Find(&favorites).Error
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

// Mengecek apakah produk sudah ada di favorit user
func (r *FavoriteRepository) CheckIfExists(userID string, productID string) (bool, error) {
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
func (r *FavoriteRepository) DeleteByUserAndProduct(userID string, favoriteID string) error {
	fmt.Println("userID:", userID)
	fmt.Println("favoriteID:", favoriteID)
	result := r.db.
		Where("user_id = ? AND id = ?", userID, favoriteID).
		Delete(&models.Favorite{})

	fmt.Println("RowsAffected:", result.RowsAffected)

	return result.Error
}
func (r *FavoriteRepository) GenerateNextID() (string, error) {
	return generateNextPrefixedID(r.db, &models.Favorite{}, "FAV")
}
