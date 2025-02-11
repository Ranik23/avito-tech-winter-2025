package postgres

import (
	"avito/internal/logger"
	"avito/internal/models"
	"context"
)



type Storage interface {
	AddUser(ctx context.Context, userName string) 	error
	GetUser(ctx context.Context, userName string) 	*models.User

	AddTransaction(ctx context.Context, senderName string, receiverName string, amount int) error
	AddPurchase(ctx context.Context, purchaseName string, merchName string, price int) error
}


type StorageImpl struct {
	logger *logger.Logger
}
