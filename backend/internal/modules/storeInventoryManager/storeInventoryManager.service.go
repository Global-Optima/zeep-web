package storeInventoryManager

import (
	"go.uber.org/zap"
)

type StoreInventoryManagerService interface {
}

type storeInventoryManagerService struct {
	repo               StoreInventoryManagerRepository
	transactionManager TransactionManager
	logger             *zap.SugaredLogger
}

func NewStoreInventoryManagerService(
	repo StoreInventoryManagerRepository,
	transactionManager TransactionManager,
	logger *zap.SugaredLogger,
) StoreInventoryManagerService {
	return &storeInventoryManagerService{
		repo:               repo,
		transactionManager: transactionManager,
		logger:             logger,
	}
}
