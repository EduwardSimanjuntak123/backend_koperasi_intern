package migrations

import (
	"log"

	"backend_koperasi/internal/models"

	"gorm.io/gorm"
)

// Run menjalankan seluruh migration database.
// Digunakan saat aplikasi pertama kali dijalankan
// atau ketika ada penambahan struktur tabel.
func Run(db *gorm.DB) error {

	// =====================================================
	// DEVELOPMENT ONLY
	// =====================================================
	// Hapus seluruh tabel sebelum membuat ulang.
	// Jangan digunakan di PRODUCTION karena akan
	// menghapus seluruh data yang ada.
	// =====================================================
	err := db.Migrator().DropTable(
		&models.Order_Item{},
		&models.Order{},
		&models.Cart_Item{},
		&models.Cart{},
		&models.Favorite{},
		&models.Product{},
		&models.CategoryProduct{},
		&models.Brand{},
		&models.Supplier{},
		&models.Unit{},
		&models.PaymentMethod{},
		&models.PointLocation{},
		&models.StoreMember{},
		&models.Store{},
		&models.Roles{},
		&models.User{},
	)

	if err != nil {
		log.Println("Drop table failed:", err)
		return err
	}

	// =====================================================
	// AutoMigrate akan:
	// - Membuat tabel jika belum ada
	// - Menambahkan kolom baru jika ada perubahan model
	// - Tidak menghapus kolom lama
	// =====================================================
	err = db.AutoMigrate(
		&models.Roles{},
		&models.User{},
		// ===== Master Data =====
		&models.Brand{},
		&models.Supplier{},
		&models.Unit{},
		&models.CategoryProduct{},

		// ===== User & Store =====

		&models.Store{},
		&models.StoreMember{},

		// ===== Product =====
		&models.Product{},

		// ===== Cart =====
		&models.Cart{},
		&models.Cart_Item{},
		&models.Favorite{},

		// ===== Payment =====
		&models.PointLocation{},
		&models.PaymentMethod{},

		// ===== Transaction =====
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
