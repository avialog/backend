package middleware

import (
	"firebase.google.com/go/auth"
	"github.com/avialog/backend/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	authorizationHeader = "Authorization"
)

// Gin middleware for JWT auth

func AuthJWT(client *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idToken := c.Request.Header.Get("Authorization")
		token, err := client.VerifyIDToken(c, idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Wyszukaj użytkownika w bazie danych za pomocą token.UID
		var user model.User
		if err := db.Where("firebase_id = ?", token.UID).First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		// Przechowaj lokalne ID użytkownika w kontekście, aby można go było użyć w kolejnych handlerach
		c.Set("userID", user.ID)

		c.Next()
	}
}
