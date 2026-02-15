package middleware

import (
	"go-auth/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserIDKey = "auth.userId"
	ctxRoleKey   = "auth.role"
)

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing Auth tokenString",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Auth formate",
			})
			return
		}

		scheme := strings.TrimSpace(parts[0])
		tokenString := strings.TrimSpace(parts[1])

		if !strings.EqualFold(scheme, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Header Not in Bearer Formate",
			})
			return
		}

		if !strings.EqualFold(tokenString, "") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing tokenString",
			})
			return
		}

		claims, err := auth.ParseToken(jwtSecret, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or Expired Token",
			})
			return
		}

		c.Set(ctxUserIDKey, claims.Subject)
		c.Set(ctxRoleKey, claims.Role)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (string, bool) {
	res, ok := c.Get(ctxUserIDKey)
	if !ok {
		return "", false
	}

	userID, ok := res.(string)
	return userID, ok
}

func GetRole(c *gin.Context) (string, bool) {
	res, ok := c.Get(ctxRoleKey)
	if !ok {
		return "", false
	}

	role, ok := res.(string)
	return role, ok
}
