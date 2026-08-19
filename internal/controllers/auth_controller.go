package controllers

import (
	"net/http"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/requests"
	"backend_koperasi/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// =====================================
// POST /auth/register
// =====================================
func (ac *AuthController) Register(c *gin.Context) {

	var req requests.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		NoHP:     req.NoHP,
	}

	if err := ac.authService.Register(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Register berhasil",
		"data":    user,
	})
}

// =====================================
// POST /auth/login
// =====================================
func (ac *AuthController) Login(c *gin.Context) {

	var req requests.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	token, user, err := ac.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Simpan JWT ke HttpOnly Cookie
	c.SetCookie(
		"access_token",
		token,
		60*60*24, // 1 hari
		"/",
		"",
		false, // ubah menjadi true jika sudah HTTPS
		true,  // HttpOnly
	)

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data":    user,
	})
}

// =====================================
// POST /auth/logout
// =====================================
func (ac *AuthController) Logout(c *gin.Context) {

	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout berhasil",
	})
}

// =====================================
// GET /auth/me
// Route ini HARUS memakai AuthMiddleware()
// =====================================
func (ac *AuthController) Me(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	user, err := ac.authService.Me(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User tidak ditemukan",
		})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data user berhasil diambil",
		"data":    user,
	})
}
