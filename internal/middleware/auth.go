package middleware

import (
	"github.com/avialog/backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	authorizationHeader = "Authorization"
)

// Gin middleware for JWT auth

func AuthJWT(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idToken := c.Request.Header.Get("Authorization")
		user, err := authService.ValidateToken(c, idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Przechowaj lokalne ID użytkownika w kontekście, aby można go było użyć w kolejnych handlerach
		c.Set("userID", user.ID)

		c.Next()
	}
}
