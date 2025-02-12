package controller

import (
	"avito/internal/service/auth"
	"avito/internal/service/purchase"
	"avito/internal/service/transaction"
)



type Controller struct {
	auth.AuthService
	transaction.TransactionService
	purchase.PurchaseService
}

func NewController(auth auth.AuthService,
					transaction transaction.TransactionService, 
					purchase purchase.PurchaseService) *Controller {
						return &Controller{
							AuthService: auth,
							TransactionService: transaction,
							PurchaseService: purchase,
						}
					}