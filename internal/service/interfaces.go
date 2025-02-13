package service

import (
	"avito/internal/models"
	"avito/internal/router/handlers/responses"
	"context"
)

// TransactionService - отвечает за переводы
type TransactionService interface {
	GetSentTransactions(ctx context.Context, userName string) ([]responses.Transaction, error)
	GetReceivedTransactions(ctx context.Context, userName string) ([]responses.Transaction, error)
	SendCoins(ctx context.Context, fromUser string, toUser string, amount int) error
}

// MerchService - отвечает за покупки
type MerchService interface {
	Buy(ctx context.Context, purchaserName string, itemName string) error
	GetBoughtMerch(ctx context.Context, purchaserName string) ([]responses.InventoryItem, error)
}

// AuthService - отвечает за аутентификацию
type AuthService interface {
	Authenticate(ctx context.Context, userName string, password string) (string, error)
	VerifyToken(ctx context.Context, tokenString string) (*models.User, error)
}
