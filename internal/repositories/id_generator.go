package repositories

import (
	"fmt"

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
