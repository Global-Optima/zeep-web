package stores

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
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
	if storeDTO.Name == "" {
		return nil, errors.New("store name cannot be empty")
	}
	if storeDTO.FacilityAddress == nil || storeDTO.FacilityAddress.Address == "" {
		return nil, errors.New("facility address is required")
	}

	store := mapToStoreEntity(storeDTO)
	createdStore, err := s.repo.CreateStore(store)
	if err != nil {
		return nil, err
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
	if storeDTO.ID == 0 {
		return nil, errors.New("store ID is required for update")
	}
	if storeDTO.Name == "" {
		return nil, errors.New("store name cannot be empty")
	}

	store, err := s.repo.GetStoreByID(storeDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, err
	}

	store.Name = storeDTO.Name
	store.IsFranchise = storeDTO.IsFranchise
	if storeDTO.FacilityAddress != nil {
		store.FacilityAddressID = &storeDTO.FacilityAddress.ID
	}
	store.Status = storeDTO.Status

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
	return &types.StoreDTO{
		ID:          store.ID,
		Name:        store.Name,
		IsFranchise: store.IsFranchise,
		Status:      store.Status,
		FacilityAddress: &types.FacilityAddressDTO{
			ID:      store.FacilityAddress.ID,
			Address: store.FacilityAddress.Address,
		},
	}
}

func mapToStoreEntity(dto types.StoreDTO) *data.Store {
	return &data.Store{
		Name:        dto.Name,
		IsFranchise: dto.IsFranchise,
		Status:      dto.Status,
		FacilityAddressID: func() *uint {
			if dto.FacilityAddress != nil {
				return &dto.FacilityAddress.ID
			}
			return nil
		}(),
	}
}
