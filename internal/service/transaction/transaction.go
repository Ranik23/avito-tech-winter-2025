package transaction

import (
	"avito/internal/models"
	"context"
)


type TransactionService interface {
	Transactions(ctx context.Context) []models.Transaction
}