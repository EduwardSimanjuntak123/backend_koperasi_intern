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
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) FindByID(id uint) (*models.User, error) {

	var user models.User

	err := r.db.Preload("Role").First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
