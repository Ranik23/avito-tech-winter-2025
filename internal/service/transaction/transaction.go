package transaction

import (
	"avito/internal/logger"
	"avito/internal/models"
	"avito/internal/repository/cache"
	"avito/internal/repository/db"
	"context"
	"log/slog"
)


type TransactionService interface {
	GetTransactions(ctx context.Context) ([]models.Transaction, error)
	SendCoins(ctx context.Context, fromUser string, toUser string, amount int) error
}


type TransactionServiceImpl struct {
	storage 		db.Repository
	cache 			cache.Cache
	logger 			*logger.Logger
}

func NewTransactionServiceImpl(storage db.Repository, cache cache.Cache, logger *logger.Logger) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		storage: storage,
		cache: cache,
		logger: logger,
	}
}

func (t *TransactionServiceImpl) GetTransactions(ctx context.Context) ([]models.Transaction, error) {
	transactions, err := t.storage.FindTransactions(ctx)
	if err != nil {
		t.logger.Error("failed to find transactions", slog.String("error", err.Error()))
		return nil, err 
	}
	return transactions, nil
}

func (t *TransactionServiceImpl) SendCoins(ctx context.Context, fromUser string, toUser string, amount int) error {
	if err := t.storage.CreateTransaction(ctx, fromUser, toUser, amount); err != nil {
		t.logger.Error("faield to create the transaction", slog.String("error", err.Error()))
		return err
	}

	if err := t.storage.UpdateAmount(ctx, fromUser, -amount); err != nil {
		t.logger.Error("failed to update the amount", slog.String("error", err.Error()))
		return err
	}

	if err := t.storage.UpdateAmount(ctx, toUser, amount); err != nil {
		t.logger.Error("failed to update the amount", slog.String("error", err.Error()))
		return err
	}
	return nil
}

