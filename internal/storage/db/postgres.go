package postgres

import (
	"avito/internal/logger"
	"avito/internal/models"
	"context"
)



type Repository interface {
	CreateUser(ctx context.Context, userName string, hashedPassword []byte, token string) 						error
	CreateTransaction(ctx context.Context, senderName string, receiverName string, amount int)  	error
	CreatePurchase(ctx context.Context, purchaseName string, merchName string, price int) 			error
	FindUserByToken(ctx context.Context, tokenString string)										(*models.User, error)
	FindUserByName(ctx context.Context, userName string) 											(*models.User, error)
}


type PostgresRepositoryImpl struct {
	logger *logger.Logger
}
