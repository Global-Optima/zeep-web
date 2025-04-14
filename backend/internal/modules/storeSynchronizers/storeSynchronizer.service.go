package storeSynchronizers

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeSynchronizers/types"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"time"
)

type StoreSynchronizeService interface {
	SynchronizeStoreInventory(storeID uint) error
	GetSynchronizationStatus(storeID uint) (*types.SynchronizationStatus, error)
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
	start := time.Now()
	logrus.Info("+++++++++++++++++++++SYNCHRONIZATION STARTS++++++++++++++++++++++++")
	err := s.transactionManager.SynchronizeStoreInventory(storeID)
	if err != nil {
		s.logger.Error("Error synchronizing store inventory", zap.Error(err))
		return err
	}
	logrus.Infof("++++++++++++++++++++SYNCHRONIZATION FINISHED IN %v++++++++++++++++++++", time.Since(start))
	return nil
}

func (s *storeSynchronizeService) GetSynchronizationStatus(storeID uint) (*types.SynchronizationStatus, error) {
	syncStatus, err := s.transactionManager.GetSynchronizationStatus(storeID)
	if err != nil {
		s.logger.Error("Error checking if store is synchronized", zap.Error(err))
		return nil, err
	}
	return syncStatus, nil
}
