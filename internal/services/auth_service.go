package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
	"backend_koperasi/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	authRepo *repositories.AuthRepository
}

func NewAuthService(authRepo *repositories.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

// =====================================
// Register
// =====================================
func (s *AuthService) Register(user *models.User) error {

	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username is required")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}

	if strings.TrimSpace(user.NoHP) == "" {
		return errors.New("phone number is required")
	}

	if len(user.Password) < 8 {
		return errors.New("password minimum 8 karakter")
	}

	// cek email sudah digunakan
	_, err := s.authRepo.FindByEmail(user.Email)

	if err == nil {
		return errors.New("email sudah digunakan")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// hash password
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hash

	return s.authRepo.Create(user)
}

// =====================================
// Login
// =====================================
func (s *AuthService) Login(email, password string) (string, *models.User, error) {

	if strings.TrimSpace(email) == "" {
		return "", nil, errors.New("email is required")
	}

	if strings.TrimSpace(password) == "" {
		return "", nil, errors.New("password is required")
	}

	// Cari user berdasarkan email
	user, err := s.authRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("email atau password salah")
		}
		return "", nil, err
	}

	// Cek password
	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		return "", nil, errors.New("email atau password salah")
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// =====================================
// Get Profile
// =====================================
func (s *AuthService) Me(userID uint) (*models.User, error) {

	if userID == 0 {
		return nil, errors.New("invalid user id")
	}

	return s.authRepo.FindByID(userID)
}
