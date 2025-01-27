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
	if err != nil && err != gorm.ErrRecordNotFound {
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
		IsFranchise:       createStoreDto.IsFranchise,
		ContactPhone:      createStoreDto.ContactPhone,
		ContactEmail:      createStoreDto.ContactEmail,
		StoreHours:        createStoreDto.StoreHours,
		FacilityAddressID: facilityAddress.ID,
	}

	createdStore, err := s.repo.CreateStore(store)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return mapToStoreDTO(*createdStore), nil
}

func (s *storeService) GetAllStores(filter *types.StoreFilter) ([]types.StoreDTO, error) {
	stores, err := s.repo.GetAllStores(*filter)
	if err != nil {
		return nil, err
	}

	storeDTOs := make([]types.StoreDTO, len(stores))
	for i, store := range stores {
		storeDTOs[i] = *mapToStoreDTO(store)
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
		storeDTOs[i] = *mapToStoreDTO(store)
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

	return mapToStoreDTO(*store), nil
}

func (s *storeService) UpdateStore(storeId uint, updateStoreDto types.UpdateStoreDTO) (*types.StoreDTO, error) {
	store, err := s.repo.GetStoreByID(storeId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, err
	}

	updateStoreFields(store, updateStoreDto)
	updatedStore, err := s.repo.UpdateStore(store)
	if err != nil {
		return nil, err
	}

	return mapToStoreDTO(*updatedStore), nil
}

func (s *storeService) DeleteStore(storeID uint, hardDelete bool) error {
	if storeID == 0 {
		return errors.New("store ID is required for deletion")
	}

	return s.repo.DeleteStore(storeID, hardDelete)
}

func mapToStoreDTO(store data.Store) *types.StoreDTO {
	facilityAddress := &types.FacilityAddressDTO{
		ID:        store.FacilityAddress.ID,
		Address:   store.FacilityAddress.Address,
		Longitude: safeFloat(store.FacilityAddress.Longitude),
		Latitude:  safeFloat(store.FacilityAddress.Latitude),
	}

	return &types.StoreDTO{
		ID:              store.ID,
		Name:            store.Name,
		IsFranchise:     store.IsFranchise,
		ContactPhone:    store.ContactPhone,
		ContactEmail:    store.ContactEmail,
		StoreHours:      store.StoreHours,
		FacilityAddress: facilityAddress,
	}
}

func updateStoreFields(store *data.Store, dto types.UpdateStoreDTO) {
	if dto.Name != "" {
		store.Name = dto.Name
	}
	if dto.IsFranchise {
		store.IsFranchise = dto.IsFranchise
	}
	if dto.ContactPhone != "" {
		store.ContactPhone = dto.ContactPhone
	}
	if dto.ContactEmail != "" {
		store.ContactEmail = dto.ContactEmail
	}
	if dto.StoreHours != "" {
		store.StoreHours = dto.StoreHours
	}
	if dto.FacilityAddress.Address != "" {
		store.FacilityAddress.Address = dto.FacilityAddress.Address
	}
	if dto.FacilityAddress.Latitude != nil {
		store.FacilityAddress.Latitude = dto.FacilityAddress.Latitude
	}
	if dto.FacilityAddress.Longitude != nil {
		store.FacilityAddress.Longitude = dto.FacilityAddress.Longitude
	}
}

func safeFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
