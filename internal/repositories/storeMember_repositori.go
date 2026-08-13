package repositories

import (
	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

type StoreMemberRepository struct {
	db *gorm.DB
}

func NewStoreMemberRepository(db *gorm.DB) *StoreMemberRepository {
	return &StoreMemberRepository{
		db: db,
	}
}

func (r *StoreMemberRepository) FindAll() ([]models.StoreMember, error) {

	var storeMembers []models.StoreMember

	err := r.db.
		Preload("Categories").
		Find(&storeMembers).Error

	if err != nil {
		return nil, err
	}

	return storeMembers, nil
}

func (r *StoreMemberRepository) FindByID(id uint) (*models.StoreMember, error) {

	var storeMember models.StoreMember

	err := r.db.
		Preload("Categories").
		First(&storeMember, id).Error

	if err != nil {
		return nil, err
	}

	return &storeMember, nil
}

func (r *StoreMemberRepository) FindBySlug(slug string) (*models.StoreMember, error) {

	var storeMember models.StoreMember

	err := r.db.
		Where("slug = ?", slug).
		First(&storeMember).Error

	if err != nil {
		return nil, err
	}

	return &storeMember, nil
}

func (r *StoreMemberRepository) Create(storeMember *models.StoreMember) error {
	return r.db.Create(storeMember).Error
}

func (r *StoreMemberRepository) Update(storeMember *models.StoreMember) error {
	return r.db.Save(storeMember).Error
}

func (r *StoreMemberRepository) Delete(id uint) error {
	return r.db.Delete(&models.StoreMember{}, id).Error
}
