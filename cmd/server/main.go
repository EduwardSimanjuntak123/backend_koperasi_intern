package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"backend_koperasi/config"
	"backend_koperasi/migrations"
)

func main() {
	// Connect database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Run migration
	if err := migrations.Run(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Initialize Gin
	r := gin.Default()

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Backend Koperasi API",
		})
	})

	// Run server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
