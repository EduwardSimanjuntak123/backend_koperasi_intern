package services

import (
	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type FavoriteService struct {
	repo *repositories.FavoriteRepository
}

func NewFavoriteService(repo *repositories.FavoriteRepository) *FavoriteService {
	return &FavoriteService{
		repo: repo,
	}
}

func (s *FavoriteService) GetFavorites(userID uint) ([]models.Favorite, error) {
	return s.repo.FindByUserID(userID)
}

func (s *FavoriteService) Add(userID, productID uint) error {
	return s.repo.AddProduct(userID, productID)
}

func (s *FavoriteService) Remove(userID, productID uint) error {
	return s.repo.RemoveProduct(userID, productID)
}
