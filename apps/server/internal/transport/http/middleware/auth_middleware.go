package middleware

import (
	"errors"
	"strings"

	jwtpkg "github.com/ablikhanovrm/pastebin/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const ContextUserIDKey = "userID"

func AuthMiddleware(jwtManager *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtManager.Parse(tokenStr)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) || strings.Contains(err.Error(), "expired") {
				c.AbortWithStatusJSON(401, gin.H{"error": "token expired"})
				return
			}

			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Next()
	}
}
