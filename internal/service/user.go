package service

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
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
	user, err := u.userRepository.GetByID(id)
	if err != nil {
		return model.User{}, err
	}
	user.FirstName = userRequest.FirstName
	user.LastName = userRequest.LastName
	user.AvatarURL = userRequest.AvatarURL
	user.SignatureURL = userRequest.SignatureURL
	user.Country = userRequest.Country
	user.Phone = userRequest.Phone
	user.Street = userRequest.Street
	user.City = userRequest.City
	user.Company = userRequest.Company
	user.Timezone = userRequest.Timezone

	return u.userRepository.Save(user)
}

//email nie do zmiany wywialić
