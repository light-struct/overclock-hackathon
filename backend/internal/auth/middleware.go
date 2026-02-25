package auth

import (
	"net/http"
	"strings"

	"backend/internal/config"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey = "userID"
	ContextRoleKey   = "role"
)

// JWTMiddleware проверяет Bearer-токен и кладет userID и role в контекст Gin.
func JWTMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			return
		}

		tokenStr := strings.TrimSpace(authHeader[7:])
		claims, err := ParseToken(cfg.JWTSecret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextRoleKey, claims.Role)
		c.Next()
	}
}

