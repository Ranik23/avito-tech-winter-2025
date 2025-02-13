package service

import (
	"avito/internal/logger"
	"avito/internal/repository/db"
	"avito/internal/router/handlers/responses"
	"context"
	"log/slog"
)

type merchServiceImpl struct {
	storage db.Repository
	logger  *logger.Logger
}

func NewMerchService(storage db.Repository, logger *logger.Logger) *merchServiceImpl {
	return &merchServiceImpl{storage: storage, logger: logger}
}

func (m *merchServiceImpl) GetBoughtMerch(ctx context.Context, purchaserName string) ([]responses.InventoryItem, error) {
	merch, err := m.storage.FindBoughtMerch(ctx, purchaserName)
	if err != nil {
		m.logger.Error("failed to get the merch", slog.String("error", err.Error()))
		return nil, err
	}
	return merch, nil
}

func (m *merchServiceImpl) Buy(ctx context.Context, purchaserName string, itemName string) error {
	if err := m.storage.CreatePurchase(ctx, purchaserName, itemName); err != nil {
		m.logger.Error("failed to create the purchase", slog.String("error", err.Error()))
		return err
	}
	return nil
}
