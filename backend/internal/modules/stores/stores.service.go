package stores

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type StoreService interface {
	CreateStore(storeDTO types.StoreDTO) (*types.StoreDTO, error)
	GetAllStores(searchTerm string) ([]types.StoreDTO, error)
	GetStoreByID(storeID uint) (*types.StoreDTO, error)
	UpdateStore(storeDTO types.StoreDTO) (*types.StoreDTO, error)
	DeleteStore(storeID uint, hardDelete bool) error
}

type storeService struct {
	repo StoreRepository
}

func NewStoreService(repo StoreRepository) StoreService {
	return &storeService{repo: repo}
}

func (s *storeService) CreateStore(storeDTO types.StoreDTO) (*types.StoreDTO, error) {
	if err := validStoreDTO(storeDTO, true, false); err != nil {
		return nil, err
	}

	var facilityAddress *data.FacilityAddress
	existingFacilityAddress, err := s.repo.GetFacilityAddressByAddress(storeDTO.FacilityAddress.Address)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error checking existing facility address: %w", err)
	}

	if existingFacilityAddress != nil {
		facilityAddress = existingFacilityAddress
	} else {
		facilityAddress = &data.FacilityAddress{
			Address:   storeDTO.FacilityAddress.Address,
			Longitude: &storeDTO.FacilityAddress.Longitude,
			Latitude:  &storeDTO.FacilityAddress.Latitude,
		}
		facilityAddress, err = s.repo.CreateFacilityAddress(facilityAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to create facility address: %w", err)
		}
	}

	if facilityAddress != nil {
		storeDTO.FacilityAddress.ID = facilityAddress.ID
	}

	store := mapToStoreEntity(storeDTO)

	createdStore, err := s.repo.CreateStore(store)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return mapToStoreDTO(*createdStore), nil
}

func (s *storeService) GetAllStores(searchTerm string) ([]types.StoreDTO, error) {
	stores, err := s.repo.GetAllStores(searchTerm)
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

func (s *storeService) UpdateStore(storeDTO types.StoreDTO) (*types.StoreDTO, error) {
	if err := validStoreDTO(storeDTO, false, true); err != nil {
		return nil, err
	}

	store, err := s.repo.GetStoreByID(storeDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, err
	}

	updateStoreFields(store, storeDTO)
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
		Status:          store.Status,
		ContactPhone:    store.ContactPhone,
		ContactEmail:    store.ContactEmail,
		StoreHours:      store.StoreHours,
		FacilityAddress: facilityAddress,
	}
}

func mapToStoreEntity(dto types.StoreDTO) *data.Store {
	facilityAddressID := dto.FacilityAddress.ID

	return &data.Store{
		Name:              dto.Name,
		IsFranchise:       dto.IsFranchise,
		Status:            dto.Status,
		ContactPhone:      dto.ContactPhone,
		ContactEmail:      dto.ContactEmail,
		StoreHours:        dto.StoreHours,
		FacilityAddressID: facilityAddressID,
	}
}

func validStoreDTO(storeDTO types.StoreDTO, isCreate, isPartialUpdate bool) error {
	if !isPartialUpdate && storeDTO.Name == "" {
		return errors.New("store name cannot be empty")
	}

	if storeDTO.ContactPhone != "" && !utils.IsValidPhone(storeDTO.ContactPhone) {
		return errors.New("invalid phone number format")
	}

	if storeDTO.ContactEmail != "" && !utils.IsValidEmail(storeDTO.ContactEmail) {
		return errors.New("invalid email format")
	}

	if !isPartialUpdate && storeDTO.StoreHours == "" {
		return errors.New("store hours cannot be empty")
	}

	if storeDTO.FacilityAddress != nil {
		if storeDTO.FacilityAddress.Address == "" {
			return errors.New("facility address cannot be empty")
		}

		if storeDTO.FacilityAddress.Latitude != 0 && !utils.IsValidLatitude(storeDTO.FacilityAddress.Latitude) {
			return errors.New("invalid latitude format")
		}
		if storeDTO.FacilityAddress.Longitude != 0 && !utils.IsValidLongitude(storeDTO.FacilityAddress.Longitude) {
			return errors.New("invalid longitude format")
		}
	}

	return nil
}

func updateStoreFields(store *data.Store, storeDTO types.StoreDTO) {
	if storeDTO.Name != "" {
		store.Name = storeDTO.Name
	}
	if storeDTO.IsFranchise {
		store.IsFranchise = storeDTO.IsFranchise
	}
	if storeDTO.ContactPhone != "" {
		store.ContactPhone = storeDTO.ContactPhone
	}
	if storeDTO.ContactEmail != "" {
		store.ContactEmail = storeDTO.ContactEmail
	}
	if storeDTO.StoreHours != "" {
		store.StoreHours = storeDTO.StoreHours
	}
	if storeDTO.FacilityAddress != nil {
		if storeDTO.FacilityAddress.Address != "" {
			store.FacilityAddress.Address = storeDTO.FacilityAddress.Address
		}
		if storeDTO.FacilityAddress.Latitude != 0 {
			store.FacilityAddress.Latitude = &storeDTO.FacilityAddress.Latitude
		}
		if storeDTO.FacilityAddress.Longitude != 0 {
			store.FacilityAddress.Longitude = &storeDTO.FacilityAddress.Longitude
		}
	}
	if storeDTO.Status != "" {
		store.Status = storeDTO.Status
	}
}

func safeFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
