package infrastructure

import (
	"context"
	"firebase.google.com/go/auth"
)

//go:generate mockgen -source=auth_client.go -destination=auth_client_mock.go -package infrastructure
type AuthClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type TokenExpireVerifier func(err error) bool
