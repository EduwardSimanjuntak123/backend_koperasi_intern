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
	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	return r.db.Preload("Role").
		First(user, user.ID).
		Error
}

func (r *AuthRepository) FindByID(id string) (*models.User, error) {

	var user models.User

	err := r.db.Preload("Role").First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
