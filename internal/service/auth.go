package service

import (
	"context"
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
	userService    UserService
	authClient     *auth.Client
}

func newAuthService(userRepository repository.UserRepository, authClient *auth.Client) AuthService {
	return &authService{userRepository: userRepository, authClient: authClient}
}

func (a *authService) ValidateToken(ctx context.Context, token string) (model.User, error) {
	response, err := a.authClient.VerifyIDToken(ctx, token)
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
	}

	user, err := a.userService.GetUser(response.UID)
	email := response.Claims["email"].(string)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// Jeżeli nie istnieje -> utwórz go podając email z claimsów i id z uid
// jeżeli user istnieje  sprawdź czy email się zmienił  jeżeli tak to zaktualizuj
// zwróceń
