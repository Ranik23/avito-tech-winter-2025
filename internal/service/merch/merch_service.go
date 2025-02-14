package merch

import (
	"avito/internal/logger"
	"avito/internal/repository"
	"avito/internal/router/handlers/responses"
	"context"
	"log/slog"
)


// ПИСАТЬ ТЕСТЫ ТУТ - БЕСПОЛЕЗНО. ТУТ НЕ ЛОГИКИ
type merchServiceImpl struct {
	storage repository.Repository
	logger  *logger.Logger
}

func NewMerchService(storage repository.Repository, logger *logger.Logger) *merchServiceImpl {
	return &merchServiceImpl{storage: storage, logger: logger}
}

func (m *merchServiceImpl) FetchBoughtMerch(ctx context.Context, purchaserName string) ([]responses.InventoryItem, error) {
	merch, err := m.storage.FindBoughtMerch(ctx, purchaserName)
	if err != nil {
		m.logger.Error("failed to get the merch", slog.String("error", err.Error()))
		return nil, err
	}
	return merch, nil
}

func (m *merchServiceImpl) Buy(ctx context.Context, purchaserName string, itemName string) error {

	merch, err := m.storage.FindMerchByName(ctx, itemName)
	if err != nil {
		m.logger.Error("failed to get the merch item", slog.String("error", err.Error()))
		return err
	}

	if err := m.storage.UpdateBalance(ctx, purchaserName, -merch.Price); err != nil {
		m.logger.Error("failed to update the balance", slog.String("error", err.Error()))
		return err
	}

	if err := m.storage.CreatePurchase(ctx, purchaserName, itemName); err != nil {
		m.logger.Error("failed to create the purchase", slog.String("error", err.Error()))
		return err
	}
	return nil
}
