package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

// ======================================
// Get All Users
// ======================================

func (r *RoleRepository) FindAll() ([]models.Roles, error) {

	var roles []models.Roles

	err := r.db.Find(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}

// ======================================
// Get Role By ID
// ======================================

func (r *RoleRepository) FindByID(id uint) (*models.Roles, error) {

	var role models.Roles

	err := r.db.First(&role, id).Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// ======================================
// Get Role By Email
// ======================================

func (r *RoleRepository) FindByEmail(email string) (*models.Roles, error) {

	var role models.Roles

	err := r.db.
		Where("email = ?", email).
		First(&role).Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// ======================================
// Get Role By Username
// ======================================

func (r *RoleRepository) FindByUsername(username string) (*models.Roles, error) {

	var role models.Roles

	err := r.db.
		Where("username = ?", username).
		First(&role).Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// ======================================
// Create Role
// ======================================

func (r *RoleRepository) Create(role *models.Roles) error {
	return r.db.Create(role).Error
}

// ======================================
// Update Role
// ======================================

func (r *RoleRepository) Update(role *models.Roles) error {
	return r.db.Save(role).Error
}

// ======================================
// Delete Role
// ======================================

func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Roles{}, id).Error
}
