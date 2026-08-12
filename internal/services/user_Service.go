package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

// ======================================
// Get All Users
// ======================================

func (s *UserService) GetAll() ([]models.User, error) {
	return s.userRepo.FindAll()
}

// ======================================
// Get User By ID
// ======================================

func (s *UserService) GetByID(id uint) (*models.User, error) {

	if id == 0 {
		return nil, errors.New("invalid user id")
	}

	return s.userRepo.FindByID(id)
}

// ======================================
// Create User
// ======================================

func (s *UserService) Create(user *models.User) error {

	// Validasi
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username is required")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}

	if strings.TrimSpace(user.Password) == "" {
		return errors.New("password is required")
	}

	if strings.TrimSpace(user.No_hp) == "" {
		return errors.New("phone number is required")
	}

	// Username sudah digunakan?
	existingUsername, _ := s.userRepo.FindByUsername(user.Username)
	if existingUsername != nil {
		return errors.New("username already exists")
	}

	// Email sudah digunakan?
	existingEmail, _ := s.userRepo.FindByEmail(user.Email)
	if existingEmail != nil {
		return errors.New("email already exists")
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return s.userRepo.Create(user)
}

// ======================================
// Update User
// ======================================

func (s *UserService) Update(id uint, user *models.User) error {

	existing, err := s.userRepo.FindByID(id)

	if err != nil {
		return err
	}

	existing.Name = user.Name
	existing.Username = user.Username
	existing.Email = user.Email
	existing.No_hp = user.No_hp

	// Update password hanya jika diisi
	if strings.TrimSpace(user.Password) != "" {

		hash, err := bcrypt.GenerateFromPassword(
			[]byte(user.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			return err
		}

		existing.Password = string(hash)
	}

	return s.userRepo.Update(existing)
}

// ======================================
// Delete User
// ======================================

func (s *UserService) Delete(id uint) error {

	_, err := s.userRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}
