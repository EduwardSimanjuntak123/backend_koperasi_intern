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

func (s *RolesService) GetByID(id string) (*models.Roles, error) {

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("invalid Roles id")
	}

	return s.RolesRepo.FindByID(id)
}

func (s *RolesService) Create(roles *models.Roles) error {

	if strings.TrimSpace(roles.Name) == "" {
		return errors.New("Roles name is required")
	}

	return s.RolesRepo.Create(roles)
}

func (s *RolesService) Update(id string, roles *models.Roles) error {

	existing, err := s.RolesRepo.FindByID(id)

	if err != nil {
		return err
	}

	existing.Name = roles.Name

	return s.RolesRepo.Update(existing)
}

func (s *RolesService) Delete(id string) error {

	_, err := s.RolesRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.RolesRepo.Delete(id)
}
