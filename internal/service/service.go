package service

import (
	"avito/internal/logger"
	"avito/internal/models"
	redis "avito/internal/storage/cache"
	postgres "avito/internal/storage/db"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SendCoins(ctx context.Context, receiver string, amount int) error
	BuyItem(ctx context.Context, itemName string) error
	Authenticate(ctx context.Context, userName string, password string) (string, error)
	GetHistoryTransactions() []models.Transaction
	GetInventory() []models.Merch
}

type UserServiceImpl struct {
	logger  logger.Logger
	storage postgres.Repository
	cache   redis.Cache
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

func (u *UserServiceImpl) Authenticate(ctx context.Context, userName string, password string) (string, error) {

	user, err := u.storage.FindUserByName(ctx, userName)
	if err != nil {
		u.logger.Error("internal server error", slog.String("error", err.Error()))
		return "", err
	}

	if user == nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			u.logger.Error("failed to hash the password", slog.String("error", err.Error()))
			return "", err
		}

		token, err := generateJWT(userName)
		if err != nil {
			u.logger.Error("failed to generate token", slog.String("error", err.Error()))
			return "", err
		}

		if err := u.storage.CreateUser(ctx, userName, hashedPassword, token); err != nil {
			u.logger.Error("failed to create user", slog.String("error", err.Error()))
			return "", err
		}

		return token, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("failed to hash the password", slog.String("error", err.Error()))
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, hashedPassword); err != nil {
		u.logger.Error("invalid password", slog.String("error", err.Error()))
		return "", errors.New("invalid credentials")
	}

	return user.Token, nil
}
