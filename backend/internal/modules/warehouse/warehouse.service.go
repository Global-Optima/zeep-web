package warehouse

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type WarehouseService interface {
	AssignStoreToWarehouse(req types.AssignStoreToWarehouseRequest) error
	ReassignStore(storeID uint, req types.ReassignStoreRequest) error
	GetAllStoresByWarehouse(warehouseID uint, pagination *utils.Pagination) ([]types.ListStoresResponse, error)

	CreateWarehouse(req types.CreateWarehouseDTO) (*types.WarehouseDTO, error)
	GetWarehouseByID(id uint) (*types.WarehouseDTO, error)
	GetWarehouses(filter *types.WarehouseFilter) ([]types.WarehouseDTO, error)
	GetAllWarehouses(filter *types.WarehouseFilter) ([]types.WarehouseDTO, error)
	UpdateWarehouse(id uint, req types.UpdateWarehouseDTO) (*types.WarehouseDTO, error)
	DeleteWarehouse(id uint) error
}

type warehouseService struct {
	repo WarehouseRepository
}

func NewWarehouseService(repo WarehouseRepository) WarehouseService {
	return &warehouseService{repo: repo}
}

func (s *warehouseService) AssignStoreToWarehouse(req types.AssignStoreToWarehouseRequest) error {
	return s.repo.AssignStoreToWarehouse(req.StoreID, req.WarehouseID)
}

func (s *warehouseService) ReassignStore(storeID uint, req types.ReassignStoreRequest) error {
	return s.repo.ReassignStoreToWarehouse(storeID, req.WarehouseID)
}

func (s *warehouseService) GetAllStoresByWarehouse(warehouseID uint, pagination *utils.Pagination) ([]types.ListStoresResponse, error) {
	stores, err := s.repo.GetAllStoresByWarehouse(warehouseID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list stores: %w", err)
	}
	return types.ConvertToListStoresResponse(stores), nil
}

func (s *warehouseService) CreateWarehouse(req types.CreateWarehouseDTO) (*types.WarehouseDTO, error) {
	facilityAddress := types.ToFacilityAddressModel(req.FacilityAddress)
	warehouse := types.ToWarehouseModel(req, 0)

	if err := s.repo.CreateWarehouse(&warehouse, &facilityAddress); err != nil {
		return nil, fmt.Errorf("failed to create warehouse: %w", err)
	}

	createdWarehouse, err := s.repo.GetWarehouseByID(warehouse.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created warehouse: %w", err)
	}

	return types.ToWarehouseDTO(*createdWarehouse), nil
}

func (s *warehouseService) GetWarehouseByID(id uint) (*types.WarehouseDTO, error) {
	warehouse, err := s.repo.GetWarehouseByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("warehouse with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch warehouse: %w", err)
	}

	return types.ToWarehouseDTO(*warehouse), nil
}

func (s *warehouseService) GetAllWarehouses(filter *types.WarehouseFilter) ([]types.WarehouseDTO, error) {
	warehouses, err := s.repo.GetAllWarehouses(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all warehouses: %w", err)
	}

	responses := make([]types.WarehouseDTO, len(warehouses))
	for i, warehouse := range warehouses {
		responses[i] = *types.ToWarehouseDTO(warehouse)
	}

	return responses, nil
}

func (s *warehouseService) GetWarehouses(filter *types.WarehouseFilter) ([]types.WarehouseDTO, error) {
	warehouses, err := s.repo.GetWarehouses(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch warehouses: %w", err)
	}

	responses := make([]types.WarehouseDTO, len(warehouses))
	for i, warehouse := range warehouses {
		responses[i] = *types.ToWarehouseDTO(warehouse)
	}

	return responses, nil
}

func (s *warehouseService) UpdateWarehouse(id uint, dto types.UpdateWarehouseDTO) (*types.WarehouseDTO, error) {
	warehouse := types.UpdateWarehouseToModel(&dto)

	if err := s.repo.UpdateWarehouse(id, warehouse); err != nil {
		return nil, fmt.Errorf("failed to update warehouse: %w", err)
	}

	updatedWarehouse, err := s.repo.GetWarehouseByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated warehouse: %w", err)
	}

	return types.ToWarehouseDTO(*updatedWarehouse), nil
}

func (s *warehouseService) DeleteWarehouse(id uint) error {
	if err := s.repo.DeleteWarehouse(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("warehouse with ID %d not found", id)
		}
		return fmt.Errorf("failed to delete warehouse: %w", err)
	}

	return nil
}
