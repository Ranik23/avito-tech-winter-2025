package service

import (
	"avito/internal/logger"
	"avito/internal/repository/db"
)

type ServiceManager struct {
	TransactionService TransactionService
	MerchService       MerchService
	AuthService        AuthService
}

func NewServiceManager(storage db.Repository, logger *logger.Logger) *ServiceManager {
	return &ServiceManager{
		TransactionService: NewTransactionService(storage, logger),
		MerchService:       NewMerchService(storage, logger),
		AuthService:        NewAuthService(storage, logger),
	}
}
