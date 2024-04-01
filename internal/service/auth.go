package service

import (
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
)

type AuthService interface {
	ValidateToken(ctx context.Context, token string) (model.User, error)
}

type authService struct {
	userRepository repository.UserRepository
	authClient     *auth.Client
}

func newAuthService(userRepository repository.UserRepository, authClient *auth.Client) AuthService {
	return &authService{userRepository: userRepository, authClient: authClient}
}

func (a *authService) ValidateToken(ctx context.Context, token string) (model.User, error) {
	var newUser model.User

	response, err := a.authClient.VerifyIDToken(ctx, token)
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
	}
	userEmail := response.Claims["email"].(string)

	user, err := a.userRepository.GetByID(response.UID)

	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			newUser, err = a.userRepository.Create(model.User{
				ID:    response.UID,
				Email: userEmail,
			})
			if err != nil {
				return model.User{}, err // internal error
			}
			return newUser, nil
		}
		return model.User{}, err // internal error
	}

	if user.Email != userEmail {
		user.Email = userEmail

		_, err = a.userRepository.Save(user)
		if err != nil {
			return model.User{}, err
		}
	}

	return user, nil
}
