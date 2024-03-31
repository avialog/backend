package service

import (
	"firebase.google.com/go/auth"
	"github.com/avialog/backend/internal/repository"
)

type AuthService interface {
}

type authService struct {
	userRepository repository.UserRepository
	authClient     *auth.Client
}

func newAuthService(userRepository repository.UserRepository, authClient *auth.Client) AuthService {
	return &authService{userRepository: userRepository, authClient: authClient}
}
