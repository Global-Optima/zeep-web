package warehouse

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
)

type WarehouseService interface {
	AssignStoreToWarehouse(req types.AssignStoreToWarehouseRequest) error
	ReassignStore(storeID uint, req types.ReassignStoreRequest) error
	GetAllStoresByWarehouse(warehouseID uint) ([]types.ListStoresResponse, error)
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

func (s *warehouseService) GetAllStoresByWarehouse(warehouseID uint) ([]types.ListStoresResponse, error) {
	stores, err := s.repo.GetAllStoresByWarehouse(warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list stores: %w", err)
	}
	return types.ConvertToListStoresResponse(stores), nil
}
