package postgres

import (
	"avito/internal/logger"
	"context"
)



type Storage interface {
	CreateUser(ctx context.Context, userName string, hashedPassword []byte) 					error
	CreateTransaction(ctx context.Context, senderName string, receiverName string, amount int)  error
	CreatePurchase(ctx context.Context, purchaseName string, merchName string, price int) 		error
}


type PostgresStorage struct {
	logger *logger.Logger
}
