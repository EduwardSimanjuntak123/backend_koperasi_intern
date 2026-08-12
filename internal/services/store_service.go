package services

import (
	"errors"
	"strings"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type StoreService struct {
	storeRepo *repositories.StoreRepository
}

func NewStoreService(repo *repositories.StoreRepository) *StoreService {
	return &StoreService{
		storeRepo: repo,
	}
}

func (s *StoreService) GetAll() ([]models.Store, error) {
	return s.storeRepo.FindAll()
}

func (s *StoreService) GetByID(id uint) (*models.Store, error) {

	if id == 0 {
		return nil, errors.New("invalid store id")
	}

	return s.storeRepo.FindByID(id)
}

func (s *StoreService) Create(store *models.Store) error {

	if strings.TrimSpace(store.NameStore) == "" {
		return errors.New("store name is required")
	}
	if strings.TrimSpace(store.Description) == "" {
		return errors.New("store description is required")
	}
	if strings.TrimSpace(store.Logo) == "" {
		return errors.New("store logo is required")
	}
	if store.ID == 0 {
		return errors.New("store id is required")
	}

	return s.storeRepo.Create(store)
}

func (s *StoreService) Update(id uint, store *models.Store) error {

	existing, err := s.storeRepo.FindByID(id)

	if err != nil {
		return err
	}

	existing.NameStore = store.NameStore
	existing.Description = store.Description
	existing.Logo = store.Logo
	existing.ID = store.ID

	return s.storeRepo.Update(existing)
}

func (s *StoreService) Delete(id uint) error {

	_, err := s.storeRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.storeRepo.Delete(id)
}
