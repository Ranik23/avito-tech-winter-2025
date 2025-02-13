package service

import (
	"avito/config"
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
	"gorm.io/gorm"
)



type Service interface {
	GetAppliedTransactions(ctx context.Context, userName string) 				([]models.Transaction, error)
	SendCoins(ctx context.Context, fromUser string, toUser string, amount int) 	error
	Buy(ctx context.Context, purchaserName string, itemName string) 			error
	GetBoughtMerch(ctx context.Context, purchaserName string) 					([]string, error)
	Authenticate(ctx context.Context, userName string, password string) 		(string, error)
	VerifyToken(ctx context.Context, tokenString string) 						(*models.User, error)
}

type ServiceImpl struct {
	storage 	db.Repository
	cache 		cache.Cache
	logger 		*logger.Logger
	cfg			*config.Config
}

func NewServiceImpl(storage db.Repository, cache cache.Cache, logger *logger.Logger, config *config.Config) *ServiceImpl {
	return &ServiceImpl{
		storage: storage,
		cache: cache,
		logger: logger,
		cfg: config,
	}
}

func (a *ServiceImpl) GetAppliedTransactions(ctx context.Context, userName string) ([]models.Transaction, error) {
	transactions, err := a.storage.FindAppliedTransactions(ctx, userName)
	if err != nil {
		a.logger.Error("failed to find transactions", slog.String("error", err.Error()))
		return nil, err 
	}
	return transactions, nil
}

func (a *ServiceImpl) SendCoins(ctx context.Context, fromUser string, toUser string, amount int) error {
	if err := a.storage.CreateTransaction(ctx, fromUser, toUser, amount); err != nil {
		a.logger.Error("faield to create the transaction", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (a *ServiceImpl) GetBoughtMerch(ctx context.Context, purchaserName string) ([]string, error) {
	merch, err := a.storage.FindBoughtMerch(ctx, purchaserName)
	if err != nil {
		a.logger.Error("failed to get the merch", slog.String("error", err.Error()))
		return nil, err
	}
	return merch, nil
}

func (a *ServiceImpl) Buy(ctx context.Context, purchaserName string, itemName string) error {
	if err := a.storage.CreatePurchase(ctx, purchaserName, itemName); err != nil {
		a.logger.Error("failed to create the purchase", slog.String("error", err.Error()))
		return err
	}

	return nil
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

func (a *ServiceImpl) Authenticate(ctx context.Context, userName string, password string) (string, error) {
	user, err := a.storage.FindUserByName(ctx, userName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound){
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

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)); err != nil {
		a.logger.Error("invalid password", slog.String("error", err.Error()))
		return "", errors.New("invalid credentials")
	}

	return user.Token, nil
}

func (a *ServiceImpl) VerifyToken(ctx context.Context, tokenString string) (*models.User, error) {
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
