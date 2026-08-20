package main

import (
	"log"
	"time"

	"backend_koperasi/config"
	"backend_koperasi/internal/controllers"
	"backend_koperasi/internal/repositories"
	"backend_koperasi/internal/routes"
	"backend_koperasi/internal/services"
	"backend_koperasi/migrations"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// ==============================
	// Connect Database
	// ==============================
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// ==============================
	// Run Migration
	// ==============================
	if err := migrations.Run(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// ==============================
	// Initialize Gin
	// ==============================
	r := gin.Default()

	// ==============================
	// CORS
	// ==============================
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000", // React
			"http://localhost:5173", // Vite
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ==============================
	// Health Check
	// ==============================
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Backend Koperasi API",
		})
	})

	// ==============================
	// Repository
	// ==============================
	authRepo := repositories.NewAuthRepository(db)
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryProductRepository(db)
	userRepo := repositories.NewUserRepository(db)
	rolesRepo := repositories.NewRoleRepository(db)
	storeRepo := repositories.NewStoreRepository(db)
	storeMemberRepo := repositories.NewStoreMemberRepository(db)
	favoriteRepo := repositories.NewFavoriteRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	brandRepo := repositories.NewBrandRepository(db)
	unitRepo := repositories.NewUnitRepository(db)

	// ==============================
	// Service
	// ==============================
	authService := services.NewAuthService(authRepo)
	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryProductService(categoryRepo)
	userService := services.NewUserService(userRepo)
	rolesService := services.NewRolesService(rolesRepo)
	storeService := services.NewStoreService(storeRepo)
	storeMemberService := services.NewStoreMemberService(storeMemberRepo)
	favoriteService := services.NewFavoriteService(favoriteRepo)
	cartService := services.NewCartService(cartRepo)
	brandService := services.NewBrandService(brandRepo)
	unitService := services.NewUnitService(unitRepo)

	// ==============================
	// Controller
	// ==============================
	authController := controllers.NewAuthController(authService)
	productController := controllers.NewProductController(productService)
	categoryController := controllers.NewCategoryProductController(categoryService)
	userController := controllers.NewUserController(userService)
	rolesController := controllers.NewRolesController(rolesService)
	storeController := controllers.NewStoreController(storeService)
	storeMemberController := controllers.NewStoreMemberController(storeMemberService)
	favoriteController := controllers.NewFavoriteController(favoriteService)
	cartController := controllers.NewCartController(cartService)
	brandController := controllers.NewBrandController(brandService)
	unitController := controllers.NewUnitController(unitService)

	// ==============================
	// Routes
	// ==============================
	routes.RegisterRoutes(
		r,

		productController,
		categoryController,
		userController,
		rolesController,
		storeController,
		storeMemberController,
		favoriteController,
		authController,
		cartController,
		unitController,
		brandController,
	)

	// ==============================
	// Run Server
	// ==============================
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
