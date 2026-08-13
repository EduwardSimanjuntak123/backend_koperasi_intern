package main

import (
	"log"

	"backend_koperasi/config"
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/repositories"
	"backend_koperasi/internal/routes"
	"backend_koperasi/internal/services"
	"backend_koperasi/migrations"

	"github.com/gin-gonic/gin"
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

	// Health Check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Backend Koperasi API",
		})
	})

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productController := controllers.NewProductController(productService)

	categoryRepo := repositories.NewCategoryProductRepository(db)
	categoryService := services.NewCategoryProductService(categoryRepo)
	categoryController := controllers.NewCategoryProductController(categoryService)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	rolesRepo := repositories.NewRoleRepository(db)
	rolesService := services.NewRolesService(rolesRepo)
	rolesController := controllers.NewRolesController(rolesService)

	storeRepo := repositories.NewStoreRepository(db)
	storeService := services.NewStoreService(storeRepo)
	storeController := controllers.NewStoreController(storeService)

	storeMemberRepo := repositories.NewStoreMemberRepository(db)
	storeMemberService := services.NewStoreMemberService(storeMemberRepo)
	storeMemberController := controllers.NewStoreMemberController(storeMemberService)

	routes.RegisterRoutes(
		r,
		productController,
		categoryController,
		userController,
		rolesController,
		storeController,
		storeMemberController,
	)

	// Run server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
