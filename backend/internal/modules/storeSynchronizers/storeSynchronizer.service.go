package storeSynchronizers

import (
	"go.uber.org/zap"
)

type StoreSynchronizeService interface {
	SynchronizeStoreInventory(storeID uint) error
	IsSynchronizedStore(storeID uint) (bool, error)
}

type storeSynchronizeService struct {
	repo               StoreSynchronizeRepository
	transactionManager TransactionManager
	logger             *zap.SugaredLogger
}

func NewStoreSynchronizeService(
	repo StoreSynchronizeRepository,
	transactionManager TransactionManager,
	logger *zap.SugaredLogger,
) StoreSynchronizeService {
	return &storeSynchronizeService{
		repo:               repo,
		transactionManager: transactionManager,
		logger:             logger,
	}
}

func (s *storeSynchronizeService) SynchronizeStoreInventory(storeID uint) error {
	err := s.transactionManager.SynchronizeStoreInventory(storeID)
	if err != nil {
		s.logger.Error("Error synchronizing store inventory", zap.Error(err))
		return err
	}
	return nil
}

func (s *storeSynchronizeService) IsSynchronizedStore(storeID uint) (bool, error) {
	isSync, err := s.repo.IsSynchronizedStore(storeID)
	if err != nil {
		s.logger.Error("Error checking if store is synchronized", zap.Error(err))
		return false, err
	}
	return isSync, nil
}
