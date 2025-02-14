package service

import (
	"avito/internal/logger"
	"avito/internal/repository"
	"avito/internal/router/handlers/responses"
	"context"
	"log/slog"
)

type transactionServiceImpl struct {
	storage repository.Repository
	logger  *logger.Logger
}

func NewTransactionService(storage repository.Repository, logger *logger.Logger) *transactionServiceImpl {
	return &transactionServiceImpl{storage: storage, logger: logger}
}

func (t *transactionServiceImpl) ListSentTransactions(ctx context.Context, userName string) ([]responses.Transaction, error) {
	transactions, err := t.storage.FindAppliedTransactions(ctx, true, userName)
	if err != nil {
		t.logger.Error("failed to find transactions", slog.String("error", err.Error()))
		return nil, err
	}

	var finalTransactions []responses.Transaction
	for _, transaction := range transactions {
		finalTransactions = append(finalTransactions, responses.Transaction{
			FromUser: transaction.Sender.Username,
			ToUser:   transaction.Receiver.Username,
			Amount:   transaction.Amount,
		})
	}
	return finalTransactions, nil
}

func (t *transactionServiceImpl) ListReceivedTransactions(ctx context.Context, userName string) ([]responses.Transaction, error) {
	transactions, err := t.storage.FindAppliedTransactions(ctx, false, userName)
	if err != nil {
		t.logger.Error("failed to find transactions", slog.String("error", err.Error()))
		return nil, err
	}

	var finalTransactions []responses.Transaction
	for _, transaction := range transactions {
		finalTransactions = append(finalTransactions, responses.Transaction{
			FromUser: transaction.Sender.Username,
			ToUser:   transaction.Receiver.Username,
			Amount:   transaction.Amount,
		})
	}
	return finalTransactions, nil
}

func (t *transactionServiceImpl) SendCoins(ctx context.Context, fromUser string, toUser string, amount int) error {
	if err := t.storage.CreateTransaction(ctx, fromUser, toUser, amount); err != nil {
		t.logger.Error("failed to create the transaction", slog.String("error", err.Error()))
		return err
	}
	return nil
}
