package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func Middleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c.Request)
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func extractToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization") // "Bearer <token>"
	if len(bearer) > 7 && strings.HasPrefix(bearer, "Bearer ") {
		return bearer[7:]
	}
	return ""
}
