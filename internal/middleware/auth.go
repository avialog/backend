package middleware

import (
	"github.com/avialog/backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	authorizationHeader = "Authorization"
)

func AuthJWT(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(authorizationHeader)

		user, err := authService.ValidateToken(c, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("userID", user.ID)
		c.Next()
	}
}
