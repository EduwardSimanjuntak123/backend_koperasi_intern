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

func (s *StoreService) GetByID(id string) (*models.Store, error) {

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("invalid store id")
	}

	return s.storeRepo.FindByID(id)
}

// Create membuat toko baru. userID diambil dari token JWT (bukan dari body request).
func (s *StoreService) Create(store *models.Store, userID string) error {

	if strings.TrimSpace(store.NameStore) == "" {
		return errors.New("store name is required")
	}
	if strings.TrimSpace(store.Description) == "" {
		return errors.New("store description is required")
	}
	if strings.TrimSpace(store.Logo) == "" {
		return errors.New("store logo is required")
	}
	if strings.TrimSpace(userID) == "" {
		return errors.New("user tidak terautentikasi")
	}

	// Ambil user ID dari token, bukan dari body request
	store.UserID = userID

	// Generate ID
	id, err := s.storeRepo.GenerateNextID()
	if err != nil {
		return err
	}
	store.ID = id

	return s.storeRepo.Create(store)
}

func (s *StoreService) Update(id string, store *models.Store) error {

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

func (s *StoreService) Delete(id string) error {

	_, err := s.storeRepo.FindByID(id)

	if err != nil {
		return err
	}

	return s.storeRepo.Delete(id)
}
