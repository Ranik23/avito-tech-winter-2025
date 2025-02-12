package auth

import "context"


type AuthService interface {
	Authenticate(ctx context.Context, userName string, password string) (string, error)
	VerifyToken(ctx context.Context, tokenString string) (bool, error)
}