package db

import (
	"avito/internal/models"
	"context"
)


type Repository interface {
	CreateUser(ctx context.Context, userName string, hashedPassword []byte, tokenString string) 	error
	CreateTransaction(ctx context.Context, senderName string, receiverName string, amount int)  	error
	CreatePurchase(ctx context.Context, purchaserName string, merchName string) 			        error
	FindUserByName(ctx context.Context, userName string) 											(*models.User, error)
	FindTransactions(ctx context.Context)														 	([]models.Transaction, error)
	FindMerch(ctx context.Context)																	([]models.Merch, error)
}
