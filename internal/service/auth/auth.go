package auth

import (
	"avito/internal/apperror"
	"avito/internal/logger"
	"avito/internal/models"
	"avito/internal/repository/cache"
	"avito/internal/repository/db"
	"context"
	"errors"
	"log/slog"
	"time"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)


type AuthService interface {
	Authenticate(ctx context.Context, userName string, password string) (string, error)
	VerifyToken(ctx context.Context, tokenString string) 				(*models.User, error)
}


type AuthServiceImpl struct {
	storage 	db.Repository
	cache   	cache.Cache
	logger 		*logger.Logger
}

func NewAuthServiceImpl(strg db.Repository, cache cache.Cache, logger *logger.Logger) *AuthServiceImpl {
	return &AuthServiceImpl{
		storage: strg,
		cache: cache,
		logger: logger,
	}
}

var jwtSecret = []byte("your_secret_key")

func generateJWT(userName string) (string, error) {
	claims := jwt.MapClaims{
		"user_name": userName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}


func (a *AuthServiceImpl) Authenticate(ctx context.Context, userName string, password string) (string, error) {
	user, err := a.storage.FindUserByName(ctx, userName)
	if err != nil {
		a.logger.Error("internal server error", slog.String("error", err.Error()))
		return "", err
	}

	if user == nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			a.logger.Error("failed to hash the password", slog.String("error", err.Error()))
			return "", err
		}

		token, err := generateJWT(userName)
		if err != nil {
			a.logger.Error("failed to generate token", slog.String("error", err.Error()))
			return "", err
		}

		if err := a.storage.CreateUser(ctx, userName, hashedPassword, token); err != nil {
			a.logger.Error("failed to create user", slog.String("error", err.Error()))
			return "", err
		}

		return token, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error("failed to hash the password", slog.String("error", err.Error()))
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, hashedPassword); err != nil {
		a.logger.Error("invalid password", slog.String("error", err.Error()))
		return "", errors.New("invalid credentials")
	}

	return user.Token, nil
}

func (a *AuthServiceImpl) VerifyToken(ctx context.Context, tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.ErrInvalidToken
		}
		return jwtSecret, nil
	})

	if err != nil {
		a.logger.Error("invalid token", slog.String("error", err.Error()))
		return nil, apperror.ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userName, ok := claims["user_name"].(string)
		if !ok {
			return nil, apperror.ErrInvalidToken
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, apperror.ErrInvalidToken
		}

		if time.Now().Unix() > int64(exp) {
			return nil, apperror.ErrExpiredToken
		}

		user, err := a.storage.FindUserByName(ctx, userName)
		if err != nil {
			a.logger.Error("failed to find the user", slog.String("error", err.Error()))
			return nil, err
		}

		return user, nil
	}

	return nil, apperror.ErrInvalidToken
}