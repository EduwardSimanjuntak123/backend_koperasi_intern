package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{
		db: db,
	}
}

func (r *StoreRepository) FindAll() ([]models.Store, error) {

	var stores []models.Store

	err := r.db.
		Find(&stores).Error

	if err != nil {
		return nil, err
	}

	return stores, nil
}

func (r *StoreRepository) FindByID(id uint) (*models.Store, error) {

	var store models.Store

	err := r.db.
		First(&store, id).Error

	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (r *StoreRepository) FindBySlug(slug string) (*models.Store, error) {

	var store models.Store

	err := r.db.
		Where("slug = ?", slug).
		First(&store).Error

	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (r *StoreRepository) Create(store *models.Store) error {
	return r.db.Create(store).Error
}

func (r *StoreRepository) Update(store *models.Store) error {
	return r.db.Save(store).Error
}

func (r *StoreRepository) Delete(id uint) error {
	return r.db.Delete(&models.Store{}, id).Error
}
