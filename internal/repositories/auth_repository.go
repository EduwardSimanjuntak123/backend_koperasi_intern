package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) FindByEmail(email string) (*models.User, error) {

	var user models.User

	err := r.db.
		Preload("Role").
		Preload("Building").
		Preload("Floor").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	return r.db.
		Preload("Role").
		Preload("Building").
		Preload("Floor").
		Where("id = ?", user.ID).
		First(user).Error
}

func (r *AuthRepository) FindByID(id string) (*models.User, error) {

	var user models.User

	err := r.db.
		Preload("Role").
		Preload("Building").
		Preload("Floor").
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetMaxUserNumberByPrefix(prefix string) (int, error) {
	var result struct {
		MaxNumber int
	}

	err := r.db.
		Model(&models.User{}).
		Select("COALESCE(MAX(CAST(SPLIT_PART(id, '-', 2) AS INTEGER)), 0) AS max_number").
		Where("id LIKE ?", prefix+"-%").
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.MaxNumber, nil
}
