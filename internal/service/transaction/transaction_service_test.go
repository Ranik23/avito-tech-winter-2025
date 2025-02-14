package transaction_test

import (
	"avito/internal/logger"
	"avito/internal/models"
	transaction_mock "avito/internal/repository/mock"
	"avito/internal/service/transaction"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListSentTransactions(t *testing.T) {
	mockRepository := &transaction_mock.MockRepositoryImpl{}
	log := logger.NewLogger("debug")

	transactionService := transaction.NewTransactionService(mockRepository, log)

	mockSender := &models.User{
		ID:             1,
		Username:       "user",
		Balance:        1000,
	}

	mockReceiver1 := &models.User{
		ID:       2,
		Username: "receiver_1",
		Balance:  500,
	}

	mockReceiver2 := &models.User{
		ID:       3,
		Username: "receiver_2",
		Balance:  200,
	}


	mockTransactions := []models.Transaction{
		{
			ID:         1,
			SenderID:   mockSender.ID,
			ReceiverID: mockReceiver1.ID,
			Sender:     *mockSender,  
			Receiver:   *mockReceiver1,
			Amount:     100,
		},
		{
			ID:         2,
			SenderID:   mockSender.ID,
			ReceiverID: mockReceiver2.ID,
			Sender:     *mockSender,
			Receiver:   *mockReceiver2,
			Amount:     400,
		},
	}

	mockRepository.On("FindUserByName", mock.Anything, "user").Return(mockSender, nil)
	mockRepository.On("FindTransactions", mock.Anything).Return(mockTransactions, nil)

	transactions, err := transactionService.ListSentTransactions(context.Background(), "user")

	assert.NoError(t, err)

	assert.Len(t, transactions, len(mockTransactions))

	for i, tx := range transactions {
		assert.Equal(t, mockTransactions[i].Sender.Username, tx.FromUser)
		assert.Equal(t, mockTransactions[i].Receiver.Username, tx.ToUser)
		assert.Equal(t, mockTransactions[i].Amount, tx.Amount)
	}
}

func TestListReceivedTransactions(t *testing.T) {
	mockRepository := &transaction_mock.MockRepositoryImpl{}
	log := logger.NewLogger("debug")

	transactionService := transaction.NewTransactionService(mockRepository, log)

	mockReceiver := &models.User{
		ID:             1,
		Username:       "user",
		Balance:        1000,
	}

	mockSender1 := &models.User{
		ID:       2,
		Username: "sender_1",
		Balance:  500,
	}

	mockSender2 := &models.User{
		ID:       3,
		Username: "sender_2",
		Balance:  1000,
	}


	mockTransactions := []models.Transaction{
		{
			ID:         1,
			SenderID:   mockSender1.ID,
			ReceiverID: mockReceiver.ID,
			Sender:     *mockSender1,  
			Receiver:   *mockReceiver,
			Amount:     100,
		},
		{
			ID:         2,
			SenderID:   mockSender2.ID,
			ReceiverID: mockReceiver.ID,
			Sender:     *mockSender2,
			Receiver:   *mockReceiver,
			Amount:     400,
		},
	}

	mockRepository.On("FindUserByName", mock.Anything, "user").Return(mockReceiver, nil)
	mockRepository.On("FindTransactions", mock.Anything).Return(mockTransactions, nil)

	transactions, err := transactionService.ListReceivedTransactions(context.Background(), "user")

	assert.NoError(t, err)

	assert.Len(t, transactions, len(mockTransactions))

	for i, tx := range transactions {
		assert.Equal(t, mockTransactions[i].Sender.Username, tx.FromUser)
		assert.Equal(t, mockTransactions[i].Receiver.Username, tx.ToUser)
		assert.Equal(t, mockTransactions[i].Amount, tx.Amount)
	}
}

