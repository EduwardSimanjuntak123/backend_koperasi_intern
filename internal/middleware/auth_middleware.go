package middleware

import (
	"net/http"
	"strings"

	"backend_koperasi/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Cek Authorization header Bearer token atau HttpOnly Cookie
		var token string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				token = authHeader
			}
		}

		if token == "" {
			cookieToken, err := c.Cookie("access_token")
			if err == nil {
				token = cookieToken
			}
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized: Token tidak ditemukan. Silakan login atau sertakan Authorization Bearer token.",
			})
			return
		}

		// Validasi token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token: " + err.Error(),
			})
			return
		}

		// Simpan data user ke context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role_id", claims.RoleID)

		c.Next()
	}
}

