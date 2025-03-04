package stores

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
	"gorm.io/gorm"
)

type StoreService interface {
	CreateStore(storeDTO *types.CreateStoreDTO) (uint, error)
	GetAllStores(filter *types.StoreFilter) ([]types.StoreDTO, error)
	GetAllStoresForNotifications() ([]types.StoreDTO, error)
	GetStoreByID(storeID uint) (*types.StoreDTO, error)
	GetStores(filter *types.StoreFilter) ([]types.StoreDTO, error)
	UpdateStore(storeId uint, storeDTO *types.UpdateStoreDTO) error
	DeleteStore(storeID uint, hardDelete bool) error
}

type storeService struct {
	repo   StoreRepository
	logger *zap.SugaredLogger
}

func NewStoreService(repo StoreRepository, logger *zap.SugaredLogger) StoreService {
	return &storeService{
		repo:   repo,
		logger: logger,
	}
}

func (s *storeService) CreateStore(createStoreDto *types.CreateStoreDTO) (uint, error) {
	existingFacilityAddress, err := s.repo.GetFacilityAddressByAddress(createStoreDto.FacilityAddress.Address)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr := fmt.Errorf("error checking existing facility address: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	store, err := types.CreateStoreFields(createStoreDto)
	if err != nil {
		wrappedErr := fmt.Errorf("error creating store: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	if existingFacilityAddress != nil {
		store.FacilityAddress = *existingFacilityAddress
	}

	id, err := s.repo.CreateStore(store)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create store: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *storeService) GetAllStores(filter *types.StoreFilter) ([]types.StoreDTO, error) {
	stores, err := s.repo.GetAllStores(filter)
	if err != nil {
		return nil, err
	}

	storeDTOs := make([]types.StoreDTO, len(stores))
	for i, store := range stores {
		storeDTOs[i] = *types.MapToStoreDTO(&store)
	}

	return storeDTOs, nil
}

func (s *storeService) GetAllStoresForNotifications() ([]types.StoreDTO, error) {
	stores, err := s.repo.GetAllStoresForNotifications()
	if err != nil {
		return nil, err
	}

	storeDTOs := make([]types.StoreDTO, len(stores))
	for i, store := range stores {
		storeDTOs[i] = *types.MapToStoreDTO(&store)
	}

	return storeDTOs, nil
}

func (s *storeService) GetStoreByID(storeID uint) (*types.StoreDTO, error) {
	if storeID == 0 {
		return nil, errors.New("store ID cannot be zero")
	}

	store, err := s.repo.GetStoreByID(storeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStoreNotFound
		}
		return nil, err
	}

	return types.MapToStoreDTO(store), nil
}

func (s *storeService) GetStores(filter *types.StoreFilter) ([]types.StoreDTO, error) {
	stores, err := s.repo.GetStores(filter)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get stores: %w", err)
		return nil, wrappedErr
	}

	storeDTOs := make([]types.StoreDTO, len(stores))
	for i, store := range stores {
		storeDTOs[i] = *types.MapToStoreDTO(&store)
	}
	return storeDTOs, nil
}

func (s *storeService) UpdateStore(storeID uint, updateStoreDto *types.UpdateStoreDTO) error {
	updateModels, err := types.UpdateStoreFields(updateStoreDto)
	if err != nil {
		wrappedErr := fmt.Errorf("validation failed trying to update store: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	err = s.repo.UpdateStore(storeID, updateModels)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update store: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *storeService) DeleteStore(storeID uint, hardDelete bool) error {
	if storeID == 0 {
		wrappedErr := fmt.Errorf("store ID is required for deletion")
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	if err := s.repo.DeleteStore(storeID, hardDelete); err != nil {
		wrappedErr := fmt.Errorf("failed to delete store: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}
