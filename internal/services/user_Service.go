package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo     *repositories.UserRepository
	buildingRepo *repositories.BuildingRepository
	floorRepo    *repositories.FloorRepository
}

func NewUserService(
	repo *repositories.UserRepository,
	buildingRepo *repositories.BuildingRepository,
	floorRepo *repositories.FloorRepository,
) *UserService {
	return &UserService{
		userRepo:     repo,
		buildingRepo: buildingRepo,
		floorRepo:    floorRepo,
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

func (s *UserService) GetByID(id string) (*models.User, error) {

	if strings.TrimSpace(id) == "" {
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

	if strings.TrimSpace(user.NoHP) == "" {
		return errors.New("phone number is required")
	}
	if user.RoleID == 0 {
		return errors.New("role id is required")
	}

	// Validasi opsional Gedung & Lantai
	if user.BuildingID != nil && strings.TrimSpace(*user.BuildingID) == "" {
		user.BuildingID = nil
	}
	if user.FloorID != nil && strings.TrimSpace(*user.FloorID) == "" {
		user.FloorID = nil
	}

	if user.BuildingID != nil {
		building, err := s.buildingRepo.FindByID(*user.BuildingID)
		if err != nil || building == nil {
			return errors.New("gedung tidak ditemukan")
		}
	}

	if user.FloorID != nil {
		floor, err := s.floorRepo.FindByID(*user.FloorID)
		if err != nil || floor == nil {
			return errors.New("lantai tidak ditemukan")
		}
		if user.BuildingID != nil && floor.BuildingID != *user.BuildingID {
			return errors.New("lantai tidak berada pada gedung yang dipilih")
		}
		if user.BuildingID == nil {
			user.BuildingID = &floor.BuildingID
		}
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

	// Generate ID
	userID, err := s.userRepo.GenerateNextID(user.RoleID)
	if err != nil {
		return err
	}
	user.ID = userID

	return s.userRepo.Create(user)
}

// ======================================
// Update User
// ======================================

func (s *UserService) Update(id string, user *models.User) error {

	existing, err := s.userRepo.FindByID(id)

	if err != nil {
		return err
	}

	existing.Name = user.Name
	existing.Username = user.Username
	existing.Email = user.Email
	existing.NoHP = user.NoHP
	existing.RoleID = user.RoleID

	// Validasi opsional Gedung & Lantai
	if user.BuildingID != nil && strings.TrimSpace(*user.BuildingID) == "" {
		user.BuildingID = nil
	}
	if user.FloorID != nil && strings.TrimSpace(*user.FloorID) == "" {
		user.FloorID = nil
	}

	if user.BuildingID != nil {
		building, err := s.buildingRepo.FindByID(*user.BuildingID)
		if err != nil || building == nil {
			return errors.New("gedung tidak ditemukan")
		}
	}

	if user.FloorID != nil {
		floor, err := s.floorRepo.FindByID(*user.FloorID)
		if err != nil || floor == nil {
			return errors.New("lantai tidak ditemukan")
		}
		if user.BuildingID != nil && floor.BuildingID != *user.BuildingID {
			return errors.New("lantai tidak berada pada gedung yang dipilih")
		}
		if user.BuildingID == nil {
			user.BuildingID = &floor.BuildingID
		}
	}

	existing.BuildingID = user.BuildingID
	existing.FloorID = user.FloorID

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

func (s *UserService) Delete(id string) error {

	_, err := s.userRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}
