package migrations

import (
	"log"
	"os"

	"gorm.io/gorm"

	"backend_koperasi/internal/models"
)

func Run(db *gorm.DB) error {
	if os.Getenv("DROP_ALL_TABLES") == "true" {
		log.Println("Dropping all tables in public schema...")

		if err := db.Exec(`
			DO $$
			DECLARE
				r RECORD;
			BEGIN
				FOR r IN (
					SELECT tablename
					FROM pg_tables
					WHERE schemaname = 'public'
				) LOOP
					EXECUTE 'DROP TABLE IF EXISTS "' || r.tablename || '" CASCADE';
				END LOOP;
			END $$;
		`).Error; err != nil {
			log.Println("Drop tables failed:", err)
			return err
		}
	}

	log.Println("Running AutoMigrate...")

	err := db.AutoMigrate(
		&models.Roles{},
		&models.User{},
		&models.Store{},
		&models.Brand{},
		&models.Unit{},
		&models.CategoryProduct{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Favorite{},

		&models.Building{},
		&models.Floor{},
		&models.Courier{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.OrderStatusHistory{},
	)

	if err != nil {
		log.Println("Migration failed:", err)
		return err
	}

	if err := db.Exec(`
		ALTER TABLE products
		DROP COLUMN IF EXISTS min_stock,
		DROP COLUMN IF EXISTS max_stock
	`).Error; err != nil {
		log.Println("Removing obsolete stock columns failed:", err)
		return err
	}

	log.Println("Migration completed successfully")

	return nil
}
