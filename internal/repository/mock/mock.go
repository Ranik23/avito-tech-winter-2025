package mock

import (
	"avito/internal/models"
	"avito/internal/router/handlers/responses"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockRepositoryImpl struct {
	mock.Mock
}

func (m *MockRepositoryImpl) CreateUser(ctx context.Context, userName string, hashedPassword []byte, tokenString string) error {
	args := m.Called(ctx, userName, hashedPassword, tokenString)
	return args.Error(0)
}

func (m *MockRepositoryImpl) CreateTransaction(ctx context.Context, senderName string, receiverName string, amount int) error {
	args := m.Called(ctx, senderName, receiverName, amount)
	return args.Error(0)
}

func (m *MockRepositoryImpl) CreatePurchase(ctx context.Context, purchaserName string, merchName string) error {
	args := m.Called(ctx, purchaserName, merchName)
	return args.Error(0)
}

func (m *MockRepositoryImpl) FindUserByName(ctx context.Context, userName string) (*models.User, error) {
	args := m.Called(ctx, userName)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoryImpl) FindBoughtMerch(ctx context.Context, purchaserName string) ([]responses.InventoryItem, error) {
	args := m.Called(ctx, purchaserName)
	if args.Get(0) != nil {
		return args.Get(0).([]responses.InventoryItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoryImpl) FindTransactions(ctx context.Context) ([]models.Transaction, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockRepositoryImpl) FindMerchByName(ctx context.Context, merchName string) (*models.Merch, error) {
	args := m.Called(ctx, merchName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Merch), args.Error(1)
}

func (m *MockRepositoryImpl) UpdateBalance(ctx context.Context, userName string, amount int) error {
	args := m.Called(ctx, userName, amount)
	return args.Error(0)
}