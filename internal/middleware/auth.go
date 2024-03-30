package middleware

import (
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	authorizationHeader = "Authorization"
)

// Gin middleware for JWT auth

func AuthJWT(client *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := client.VerifyIDToken(c.Request.Context(), c.GetHeader(authorizationHeader))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("userID", token)
		c.Next()
	}
}
