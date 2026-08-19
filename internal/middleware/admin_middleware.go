package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		roleID := c.MustGet("role_id").(uint)

		// RoleID = 1 adalah Admin
		if roleID != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden",
			})
			return
		}

		c.Next()
	}
}
