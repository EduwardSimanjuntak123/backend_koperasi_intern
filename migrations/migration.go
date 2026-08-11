package migrations

import (
	"log"

	"gorm.io/gorm"

	"backend_koperasi/internal/models"
)

func Run(db *gorm.DB) error {
	err := db.AutoMigrate(

		&models.Category{},
		&models.Product{},
	)

	if err != nil {
		log.Println("Migration failed:", err)
		return err
	}

	log.Println("Migration completed successfully")

	return nil
}
