package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// ======================================
// Get All Users
// ======================================

func (r *UserRepository) FindAll() ([]models.User, error) {

	var users []models.User

	err := r.db.Preload("Roles").Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// ======================================
// Get User By ID
// ======================================

func (r *UserRepository) FindByID(id uint) (*models.User, error) {

	var user models.User

	err := r.db.Preload("Roles").First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ======================================
// Get User By Email
// ======================================

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {

	var user models.User

	err := r.db.Preload("Roles").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ======================================
// Get User By Username
// ======================================

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {

	var user models.User

	err := r.db.Preload("Roles").
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ======================================
// Create User
// ======================================

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// ======================================
// Update User
// ======================================

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// ======================================
// Delete User
// ======================================

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
