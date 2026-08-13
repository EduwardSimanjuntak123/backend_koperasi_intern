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

// ======================================================
// Mengambil seluruh favorit milik user
// ======================================================

func (r *FavoriteRepository) FindByUserID(userID uint) ([]models.Favorite, error) {
	var favorites []models.Favorite

	err := r.db.
		Preload("Product").
		Where("user_id = ?", userID).
		Find(&favorites).Error

	return favorites, err
}

// ======================================================
// Mengecek apakah produk sudah difavoritkan
// ======================================================

func (r *FavoriteRepository) CheckIfExists(userID, productID uint) (bool, error) {
	var count int64

	err := r.db.Model(&models.Favorite{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error

	return count > 0, err
}

// ======================================================
// Menambahkan produk ke favorit
// ======================================================

func (r *FavoriteRepository) AddProduct(userID, productID uint) error {

	exists, err := r.CheckIfExists(userID, productID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("product already in favorites")
	}

	favorite := models.Favorite{
		UserID:    userID,
		ProductID: productID,
	}

	return r.db.Create(&favorite).Error
}

// ======================================================
// Menghapus produk dari favorit
// ======================================================

func (r *FavoriteRepository) RemoveProduct(userID, productID uint) error {

	result := r.db.
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.Favorite{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("favorite not found")
	}

	return nil
}
