package migrations

import (
	"log"

	"gorm.io/gorm"

	"backend_koperasi/internal/models"
)

func Run(db *gorm.DB) error {
	err := db.AutoMigrate(

		&models.CategoryProduct{},
		&models.User{},
		&models.Store{},
		&models.StoreMember{},
		&models.Roles{},
		&models.Product{},
		&models.Cart{},
		&models.Cart_Item{},
		&models.Favorite{},
		&models.PointLocation{},
		&models.PaymentMethod{},

		&models.Order{},
		&models.Order_Item{},
	)

	if err != nil {
		log.Println("Migration failed:", err)
		return err
	}

	log.Println("Migration completed successfully")

	return nil
}
