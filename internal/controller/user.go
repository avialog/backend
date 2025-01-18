package controller

import (
	"net/http"

	"github.com/avialog/backend/internal/common"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/util"
	"github.com/gin-gonic/gin"
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

// GetUser godoc
//
// @Summary Get a user
// @Description Get a user by userID from the token
// @Tags profile
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object}      dto.UserResponse
// @Failure 500 {object}      util.HTTPError
// @Router  /profile [get]
func (u *userController) GetUser(ctx *gin.Context) {
	userID := ctx.GetString(common.UserID)

	user, err := u.userService.GetUser(userID)
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	userResponse := u.adaptUser(user)

	ctx.JSON(http.StatusOK, userResponse)
}

// UpdateProfile godoc
//
// @Summary Update user profile
// @Description Update user profile information
// @Tags profile
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   userRequest       body     dto.UserRequest true       "User profile information to update"
// @Success 200 {object}      dto.UserResponse
// @Failure 400 {object}      util.HTTPError
// @Failure 500 {object}      util.HTTPError
// @Router  /profile [put]
func (u *userController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetString(common.UserID)

	var userRequest dto.UserRequest
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := u.userService.UpdateProfile(userID, userRequest)
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
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
		Address:      user.Address,
		Timezone:     user.Timezone,
	}
}
