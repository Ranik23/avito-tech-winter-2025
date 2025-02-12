package usecase

import (
	"avito/internal/logger"
	"avito/internal/models"
	"avito/internal/storage/postgres"
	"avito/internal/storage/redis"
	"context"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type UserCase interface {
	SendCoins(ctx context.Context, receiver string, amount int) 			error 
	BuyItem(ctx context.Context, itemName string) 							error
	Authenticate(ctx context.Context, userName string, password string)  	(string, error)
	GetHistoryTransactions() 												[]models.Transaction
	GetInventory()															[]models.Merch
}


type UserOperator struct {
	logger logger.Logger
	storage postgres.Storage
	cache 	redis.Cache
}


func (u *UserOperator) Authenticate(ctx context.Context, userName string, password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("failed to hash the password", slog.String("error", err.Error()))
		return "", err
	}

	u.storage.CreateUser(ctx, userName, hashedPassword)


	return "", nil
}