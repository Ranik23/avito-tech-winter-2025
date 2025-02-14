package transaction

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
	user, err := t.storage.FindUserByName(ctx, userName)
	if err != nil {
		t.logger.Error("failed to find the user", slog.String("error", err.Error()))
		return nil, err
	}

	transactions, err := t.storage.FindTransactions(ctx)
	if err != nil {
		t.logger.Error("failed to find transactions", slog.String("error", err.Error()))
		return nil, err
	}

	var finalTransactions []responses.Transaction
	for _, transaction := range transactions {
		if transaction.SenderID == user.ID {
			finalTransactions = append(finalTransactions, responses.Transaction{
				FromUser: transaction.Sender.Username,
				ToUser:   transaction.Receiver.Username,
				Amount:   transaction.Amount,
			})
		}
	}
	return finalTransactions, nil
}

func (t *transactionServiceImpl) ListReceivedTransactions(ctx context.Context, userName string) ([]responses.Transaction, error) {
	user, err := t.storage.FindUserByName(ctx, userName)
	if err != nil {
		t.logger.Error("failed to find the user", slog.String("error", err.Error()))
		return nil, err
	}

	transactions, err := t.storage.FindTransactions(ctx)
	if err != nil {
		t.logger.Error("failed to find transactions", slog.String("error", err.Error()))
		return nil, err
	}

	var finalTransactions []responses.Transaction
	for _, transaction := range transactions {
		if transaction.ReceiverID == user.ID {
			finalTransactions = append(finalTransactions, responses.Transaction{
				FromUser: transaction.Sender.Username,
				ToUser:   transaction.Receiver.Username,
				Amount:   transaction.Amount,
			})
		}
	}
	return finalTransactions, nil
}

func (t *transactionServiceImpl) SendCoins(ctx context.Context, fromUser string, toUser string, amount int) error {	

	if err := t.storage.UpdateBalance(ctx, fromUser, -amount); err != nil {
		t.logger.Error("failed to update the balance", slog.String("error", err.Error()))
		return err
	}
	if err := t.storage.UpdateBalance(ctx, toUser, amount); err != nil {
		t.logger.Error("failed to update the balance", slog.String("error", err.Error()))
		return err
	}
	
	if err := t.storage.CreateTransaction(ctx, fromUser, toUser, amount); err != nil {
		t.logger.Error("failed to create the transaction", slog.String("error", err.Error()))
		return err
	}
	return nil
}
