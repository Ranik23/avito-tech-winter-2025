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
	FindBoughtMerch(ctx context.Context, purchaserName string)										([]responses.InventoryItem, error)
	FindTransactions(ctx context.Context)															([]models.Transaction, error)
	FindMerchByName(ctx context.Context, merchName string)											(*models.Merch, error)
	UpdateBalance(ctx context.Context, userName string, amount int) 								error							
}


// UpdateBalance(ctx context.Context, userName string, amount int) error
// 