package service

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(id uint) (model.User, error)
	UpdateProfile(id uint, userRequest dto.UserRequest) (model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
	config         dto.Config
}

func newUserService(userRepository repository.UserRepository, config dto.Config) UserService {
	return &userService{userRepository: userRepository, config: config}
}

func (u *userService) GetProfile(id uint) (model.User, error) {
	return u.userRepository.GetByID(id)
}

func (u *userService) UpdateProfile(id uint, userRequest dto.UserRequest) (model.User, error) {
	user := model.User{
		Model:        gorm.Model{ID: id},
		FirstName:    userRequest.FirstName,
		LastName:     userRequest.LastName,
		Email:        userRequest.Email,
		AvatarURL:    userRequest.AvatarURL,
		SignatureURL: userRequest.SignatureURL,
		Country:      userRequest.Country,
		Phone:        userRequest.Phone,
		Street:       userRequest.Street,
		City:         userRequest.City,
		Company:      userRequest.Company,
		Timezone:     userRequest.Timezone,
	}
	return u.userRepository.Update(user)
}
