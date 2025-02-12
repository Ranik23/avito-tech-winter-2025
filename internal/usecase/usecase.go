package usecase


import (
	"avito/internal/storage/postgres"
)

type UserCase interface {
	SendCoins(receiver string, amount int) error 
	BuyItem(itemName string) error
	Authenticate(userName string, password string) (string, error)
}


type UserOperator struct {
	strg postgres.Storage
}