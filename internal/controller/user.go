package controller

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController interface {
	GetUser(*gin.Context)
	UpdateProfile(*gin.Context)
}

type userController struct {
	userService service.UserService
}

func newUserController(userService service.UserService) UserController {
	return &userController{userService: userService}
}

func (u *userController) GetUser(ctx *gin.Context) {
	// TODO: add getting user id from JWT token
	userID := uint(1)

	user, err := u.userService.GetUser(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponse := u.adaptUser(user)

	ctx.JSON(http.StatusOK, userResponse)
}

func (u *userController) UpdateProfile(ctx *gin.Context) {
	// TODO: add getting user id from JWT token
	userID := uint(1)

	var userRequest dto.UserRequest
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.userService.UpdateProfile(userID, userRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponse := u.adaptUser(user)

	ctx.JSON(http.StatusOK, userResponse)
}

func (u *userController) adaptUser(user model.User) dto.UserResponse {
	return dto.UserResponse{
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
}
