package middleware

import (
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

func AuthJWT(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(authorizationHeader)
		token = strings.Replace(token, "Bearer ", "", 1)
		user, err := authService.ValidateToken(c, token)
		if err != nil {
			if errors.Is(err, dto.ErrNotAuthorized) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("userID", user.ID)
		c.Next()
	}
}
