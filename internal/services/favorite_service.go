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

// Mengambil keranjang/wadah favorit milik satu user beserta isi produknya
func (s *FavoriteService) GetUserFavorites(userID uint) (*models.Favorite, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}
	return s.favoriteRepo.FindByUserID(userID)
}

// Menambahkan produk ke daftar favorit
func (s *FavoriteService) AddToFavorite(userID uint, productID uint) error {
	if userID == 0 || productID == 0 {
		return errors.New("user id and product id are required")
	}

	// LOGIKA BISNIS: Cek apakah user sudah mem-favoritkan produk ini sebelumnya
	isFavorited, err := s.favoriteRepo.CheckIfExists(userID, productID)
	if err != nil {
		return err
	}
	
	if isFavorited {
		return errors.New("product is already in favorites")
	}

	// Gunakan AddProduct untuk menyimpan ke tabel perantara relasi Many2Many
	return s.favoriteRepo.AddProduct(userID, productID)
}

// Menghapus dari daftar favorit
func (s *FavoriteService) RemoveFromFavorite(userID uint, productID uint) error {
	if userID == 0 || productID == 0 {
		return errors.New("user id and product id are required")
	}
	
	return s.favoriteRepo.RemoveProduct(userID, productID)
}