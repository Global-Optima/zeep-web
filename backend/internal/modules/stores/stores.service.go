package stores

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
	"gorm.io/gorm"
)

type StoreService interface {
	CreateStore(storeDTO types.CreateStoreDTO) (*types.StoreDTO, error)
	GetAllStores(filter *types.StoreFilter) ([]types.StoreDTO, error)
	GetAllStoresForNotifications() ([]types.StoreDTO, error)
	GetStoreByID(storeID uint) (*types.StoreDTO, error)
	GetStores(filter *types.StoreFilter) ([]types.StoreDTO, error)
	UpdateStore(storeId uint, storeDTO types.UpdateStoreDTO) (*types.StoreDTO, error)
	DeleteStore(storeID uint, hardDelete bool) error
}

type storeService struct {
	repo StoreRepository
}

func NewStoreService(repo StoreRepository) StoreService {
	return &storeService{repo: repo}
}

func (s *storeService) CreateStore(createStoreDto types.CreateStoreDTO) (*types.StoreDTO, error) {
	var facilityAddress *data.FacilityAddress
	existingFacilityAddress, err := s.repo.GetFacilityAddressByAddress(createStoreDto.FacilityAddress.Address)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking existing facility address: %w", err)
	}

	if existingFacilityAddress != nil {
		facilityAddress = existingFacilityAddress
	} else {
		facilityAddress = &data.FacilityAddress{
			Address:   createStoreDto.FacilityAddress.Address,
			Longitude: createStoreDto.FacilityAddress.Longitude,
			Latitude:  createStoreDto.FacilityAddress.Latitude,
		}
		facilityAddress, err = s.repo.CreateFacilityAddress(facilityAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to create facility address: %w", err)
		}
	}

	store := &data.Store{
		Name:              createStoreDto.Name,
		FranchiseeID:      createStoreDto.FranchiseID,
		ContactPhone:      createStoreDto.ContactPhone,
		ContactEmail:      createStoreDto.ContactEmail,
		StoreHours:        createStoreDto.StoreHours,
		FacilityAddressID: facilityAddress.ID,
	}

	createdStore, err := s.repo.CreateStore(store)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return types.MapToStoreDTO(createdStore), nil
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
			return nil, errors.New("store not found")
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

func (s *storeService) UpdateStore(storeId uint, updateStoreDto types.UpdateStoreDTO) (*types.StoreDTO, error) {
	store, err := s.repo.GetStoreByID(storeId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, err
	}

	types.UpdateStoreFields(store, updateStoreDto)
	updatedStore, err := s.repo.UpdateStore(store)
	if err != nil {
		return nil, err
	}

	return types.MapToStoreDTO(updatedStore), nil
}

func (s *storeService) DeleteStore(storeID uint, hardDelete bool) error {
	if storeID == 0 {
		return errors.New("store ID is required for deletion")
	}

	return s.repo.DeleteStore(storeID, hardDelete)
}
