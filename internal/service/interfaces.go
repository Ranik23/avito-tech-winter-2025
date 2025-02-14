package service

import (
	"avito/internal/models"
	"avito/internal/router/handlers/responses"
	"context"
)

type TransactionService interface {
	ListSentTransactions(ctx context.Context, userName string) ([]responses.Transaction, error)
	ListReceivedTransactions(ctx context.Context, userName string) ([]responses.Transaction, error)
	SendCoins(ctx context.Context, fromUser string, toUser string, amount int) error
}

type MerchService interface {
	Buy(ctx context.Context, purchaserName string, itemName string) error
	FetchBoughtMerch(ctx context.Context, purchaserName string) ([]responses.InventoryItem, error)
}

type AuthService interface {
	Authenticate(ctx context.Context, userName string, password string) (string, error)
	VerifyToken(ctx context.Context, tokenString string) (*models.User, error)
}
