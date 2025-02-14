package repository

import (
	"avito/internal/models"
	"avito/internal/router/handlers/responses"
	"context"
)


type Repository interface {
	CreateUser(ctx context.Context, userName string, hashedPassword []byte, tokenString string) 	error
	CreateTransaction(ctx context.Context, senderName string, receiverName string, amount int)  	error
	CreatePurchase(ctx context.Context, purchaserName string, merchName string) 			        error

	FindUserByName(ctx context.Context, userName string) 											(*models.User, error)
	FindAppliedTransactions(ctx context.Context, sentORreceived bool, senderORbuyerName string)		([]models.Transaction, error)
	FindBoughtMerch(ctx context.Context, purchaserName string)										([]responses.InventoryItem, error)
}
