package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BuyerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		roleID, exists := c.Get("role_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			return
		}

		if roleID.(uint) != 2 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Only buyers can access this resource",
			})
			return
		}

		c.Next()
	}
}
