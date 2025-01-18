package service

import (
	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
)

//go:generate mockgen -source=user.go -destination=user_mock.go -package service
type UserService interface {
	GetUser(id string) (model.User, error)
	UpdateProfile(id string, userRequest dto.UserRequest) (model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
	config         config.Config
}

func newUserService(userRepository repository.UserRepository, config config.Config) UserService {
	return &userService{userRepository: userRepository, config: config}
}

func (u *userService) GetUser(id string) (model.User, error) {
	return u.userRepository.GetByID(id)
}

func (u *userService) UpdateProfile(id string, userRequest dto.UserRequest) (model.User, error) {
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
	user.Address = userRequest.Address
	user.Timezone = userRequest.Timezone

	return u.userRepository.Save(user)
}
