package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type RolesService struct {
	RolesRepo *repositories.RoleRepository
}

func NewRolesService(repo *repositories.RoleRepository) *RolesService {
	return &RolesService{
		RolesRepo: repo,
	}
}

func (s *RolesService) GetAll() ([]models.Roles, error) {
	return s.RolesRepo.FindAll()
}

func (s *RolesService) GetByID(id uint) (*models.Roles, error) {

	if id == 0 {
		return nil, errors.New("invalid Roles id")
	}

	return s.RolesRepo.FindByID(id)
}

func (s *RolesService) Create(roles *models.Roles) error {

	if strings.TrimSpace(roles.Name) == "" {
		return errors.New("Roles name is required")
	}
	if roles.StoreID == 0 {
		return errors.New("store id is required")
	}

	return s.RolesRepo.Create(roles)
}

func (s *RolesService) Update(id uint, roles *models.Roles) error {

	existing, err := s.RolesRepo.FindByID(id)

	if err != nil {
		return err
	}

	existing.Name = roles.Name
	existing.StoreID = roles.StoreID

	return s.RolesRepo.Update(existing)
}

func (s *RolesService) Delete(id uint) error {

	_, err := s.RolesRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.RolesRepo.Delete(id)
}
