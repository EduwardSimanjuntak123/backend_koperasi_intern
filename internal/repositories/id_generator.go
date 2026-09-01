package repositories

import (
	"fmt"
	"time"

	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

func generateNextPrefixedID(db *gorm.DB, model interface{}, prefix string) (string, error) {
	var result struct {
		MaxNumber int
	}

	err := db.
		Model(model).
		Select("COALESCE(MAX(CAST(SPLIT_PART(id, '-', 2) AS INTEGER)), 0) AS max_number").
		Where("id LIKE ?", prefix+"-%").
		Scan(&result).Error

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%03d", prefix, result.MaxNumber+1), nil
}

func generateInvoiceNumber(db *gorm.DB) (string, error) {
	today := time.Now().Format("20060102")
	prefix := fmt.Sprintf("INV/%s/", today)

	var count int64
	err := db.
		Model(&models.Order{}).
		Where("invoice_number LIKE ?", prefix+"%").
		Count(&count).Error

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}

