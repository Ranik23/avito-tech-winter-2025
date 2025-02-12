package purchase

import (
	"avito/internal/logger"
	"avito/internal/models"
	"avito/internal/repository/cache"
	"avito/internal/repository/db"
	"context"
	"log/slog"
)

type PurchaseService interface {
	Buy(ctx context.Context, purchaserName string, itemName string) error
	GetMerch(ctx context.Context) ([]models.Merch, error)
}

type PurchaseServiceImpl struct {
	storage db.Repository
	cache   cache.Cache
	logger  *logger.Logger
}

func NewPurchaseServiceImpl(storage db.Repository, cache cache.Cache, logger *logger.Logger) *PurchaseServiceImpl {
	return &PurchaseServiceImpl{
		storage: storage,
		cache:   cache,
		logger:  logger,
	}
}

func (p *PurchaseServiceImpl) GetMerch(ctx context.Context) ([]models.Merch, error) {
	merch, err := p.storage.FindMerch(ctx)
	if err != nil {
		p.logger.Error("failed to get the merch", slog.String("error", err.Error()))
		return nil, err
	}
	return merch, nil
}

func (p *PurchaseServiceImpl) Buy(ctx context.Context, purchaserName string, itemName string) error {
	if err := p.storage.CreatePurchase(ctx, purchaserName, itemName); err != nil {
		p.logger.Error("failed to create the purchase", slog.String("error", err.Error()))
		return err
	}

	return nil
}
