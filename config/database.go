package config

import (
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	databaseURL := os.Getenv("DATABASE_URL")

	println("DATABASE_URL:", databaseURL)

	if databaseURL == "" {
		return nil, os.ErrNotExist
	}

	db, err := gorm.Open(
		postgres.Open(databaseURL),
		&gorm.Config{},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
