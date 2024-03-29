package controller

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController interface {
	GetProfile(*gin.Context)
	UpdateProfile(*gin.Context)
}

type userController struct {
	userService service.UserService
}

func newUserController(userService service.UserService) UserController {
	return &userController{userService: userService}
}

func (u *userController) GetProfile(c *gin.Context) {
	// there will be getting user id from JWT token
	userID := uint(1)
	// -----------------------
	user, err := u.userService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // or should I return status not found (404) ?
		return
	}
	userResponse := dto.UserResponse{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		AvatarURL:    user.AvatarURL,
		SignatureURL: user.SignatureURL,
		Country:      user.Country,
		Phone:        user.Phone,
		Street:       user.Street,
		City:         user.City,
		Company:      user.Company,
		Timezone:     user.Timezone,
	}
	c.JSON(http.StatusOK, userResponse)
}

func (u *userController) UpdateProfile(c *gin.Context) {
	// there will be getting user id from JWT token
	userID := uint(1)
	// -----------------------

	var userRequest dto.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.userService.UpdateProfile(userID, userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponse := dto.UserResponse{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		AvatarURL:    user.AvatarURL,
		SignatureURL: user.SignatureURL,
		Country:      user.Country,
		Phone:        user.Phone,
		Street:       user.Street,
		City:         user.City,
		Company:      user.Company,
		Timezone:     user.Timezone,
	}
	c.JSON(200, userResponse)
}
