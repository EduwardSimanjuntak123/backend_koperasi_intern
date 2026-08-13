package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type StoreMemberService struct {
	storeMemberRepo *repositories.StoreMemberRepository
}

func NewStoreMemberService(repo *repositories.StoreMemberRepository) *StoreMemberService {
	return &StoreMemberService{
		storeMemberRepo: repo,
	}
}

func (s *StoreMemberService) GetAll() ([]models.StoreMember, error) {
	return s.storeMemberRepo.FindAll()
}

func (s *StoreMemberService) GetByID(id uint) (*models.StoreMember, error) {

	if id == 0 {
		return nil, errors.New("invalid store member id")
	}

	return s.storeMemberRepo.FindByID(id)
}

func (s *StoreMemberService) Create(storeMember *models.StoreMember) error {
	if strings.TrimSpace(storeMember.Status) == "" {
		return errors.New("store member status is required")
	}

	if storeMember.JoinedAt.IsZero() {
		return errors.New("joined_at is required")
	}

	if storeMember.StoreID == 0 {
		return errors.New("store_id is required")
	}

	if storeMember.UserID == 0 {
		return errors.New("user_id is required")
	}

	if storeMember.RoleID == 0 {
		return errors.New("role_id is required")
	}

	return s.storeMemberRepo.Create(storeMember)
}

func (s *StoreMemberService) Update(id uint, storeMember *models.StoreMember) error {

	existing, err := s.storeMemberRepo.FindByID(id)

	if err != nil {
		return err
	}

	existing.Status = storeMember.Status
	existing.JoinedAt = storeMember.JoinedAt
	existing.StoreID = storeMember.StoreID
	existing.UserID = storeMember.UserID
	existing.RoleID = storeMember.RoleID

	return s.storeMemberRepo.Update(existing)
}

func (s *StoreMemberService) Delete(id uint) error {

	_, err := s.storeMemberRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.storeMemberRepo.Delete(id)
}
