package services

import (
	"errors"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/repositories"
)

type FavoriteService struct {
	favoriteRepo *repositories.FavoriteRepository
}

func NewFavoriteService(repo *repositories.FavoriteRepository) *FavoriteService {
	return &FavoriteService{
		favoriteRepo: repo,
	}
}

func (s *FavoriteService) GetUserFavorites(userID uint) ([]models.Favorite, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}
	return s.favoriteRepo.FindByUserID(userID)
}

func (s *FavoriteService) AddToFavorite(userID uint, productID uint) error {
	if userID == 0 || productID == 0 {
		return errors.New("user id and product id are required")
	}

	// Mencegah duplikasi data sebelum terkena error unique index dari database
	isFavorited, err := s.favoriteRepo.CheckIfExists(userID, productID)
	if err != nil {
		return err
	}
	if isFavorited {
		return errors.New("product is already in favorites")
	}

	newFavorite := &models.Favorite{
		UserID:    userID,
		ProductID: productID,
	}

	return s.favoriteRepo.Create(newFavorite)
}

func (s *FavoriteService) RemoveFromFavorite(userID uint, productID uint) error {
	if userID == 0 || productID == 0 {
		return errors.New("invalid parameters")
	}

	return s.favoriteRepo.DeleteByUserAndProduct(userID, productID)
}