package purchase

import (
	"avito/internal/models"
	"context"
)


type PurchaseService interface {
	BuyItem(ctx context.Context, itemName string) error
	Inventory() []models.Merch
}