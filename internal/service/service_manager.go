package service

import (
	"avito/internal/logger"
	"avito/internal/repository"
	"avito/internal/service/auth"
	"avito/internal/service/merch"
	"avito/internal/service/transaction"
	tokenutil "avito/internal/token"
)

type ServiceManager struct {
	TransactionService TransactionService
	MerchService       MerchService
	AuthService        AuthService
}

func NewServiceManager(storage repository.Repository,
					   logger *logger.Logger,
					   token  tokenutil.Token) *ServiceManager {
	return &ServiceManager{
		TransactionService: transaction.NewTransactionService(storage, logger),
		MerchService:       merch.NewMerchService(storage, logger),
		AuthService:        auth.NewAuthService(storage, logger, token),
	}
}
