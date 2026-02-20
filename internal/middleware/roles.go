package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := GetRole(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		if !strings.EqualFold(role, "admin") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "No access Without Admin",
			})
			return
		}

		c.Next()
	}
}
