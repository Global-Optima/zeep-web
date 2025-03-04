package warehouse

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WarehouseService interface {
	AssignStoreToWarehouse(req types.AssignStoreToWarehouseRequest) error
	GetAllStoresByWarehouse(warehouseID uint, pagination *utils.Pagination) ([]types.ListStoresResponse, error)

	CreateWarehouse(req types.CreateWarehouseDTO) (*types.WarehouseDTO, error)
	GetWarehouseByID(id uint) (*types.WarehouseDTO, error)
	GetWarehouses(filter *types.WarehouseFilter) ([]types.WarehouseDTO, error)
	GetAllWarehouses(filter *types.WarehouseFilter) ([]types.WarehouseDTO, error)
	UpdateWarehouse(id uint, req types.UpdateWarehouseDTO) (*types.WarehouseDTO, error)
	DeleteWarehouse(id uint) error
}

type warehouseService struct {
	repo   WarehouseRepository
	logger *zap.SugaredLogger
}

func NewWarehouseService(repo WarehouseRepository, logger *zap.SugaredLogger) WarehouseService {
	return &warehouseService{
		repo:   repo,
		logger: logger,
	}
}

func (s *warehouseService) AssignStoreToWarehouse(req types.AssignStoreToWarehouseRequest) error {
	return s.repo.AssignStoreToWarehouse(req.StoreID, req.WarehouseID)
}

func (s *warehouseService) GetAllStoresByWarehouse(warehouseID uint, pagination *utils.Pagination) ([]types.ListStoresResponse, error) {
	stores, err := s.repo.GetAllStoresByWarehouse(warehouseID, pagination)
	if err != nil {
		s.logger.Error("failed to list stores", zap.Error(err))
		return nil, types.ErrFailedListStores
	}
	return types.ConvertToListStoresResponse(stores), nil
}

func (s *warehouseService) CreateWarehouse(req types.CreateWarehouseDTO) (*types.WarehouseDTO, error) {
	facilityAddress := types.ToFacilityAddressModel(req.FacilityAddress)
	warehouse := types.ToWarehouseModel(req, 0)

	if err := s.repo.CreateWarehouse(&warehouse, &facilityAddress); err != nil {
		s.logger.Error("failed to create warehouse", zap.Error(err))
		return nil, types.ErrFailedCreateWarehouse
	}

	createdWarehouse, err := s.repo.GetWarehouseByID(warehouse.ID)
	if err != nil {
		s.logger.Error("failed to fetch created warehouse", zap.Error(err))
		return nil, types.ErrFailedToFetchWarehouse
	}

	return types.ToWarehouseDTO(*createdWarehouse), nil
}

func (s *warehouseService) GetWarehouseByID(id uint) (*types.WarehouseDTO, error) {
	warehouse, err := s.repo.GetWarehouseByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("warehouse not found", zap.Uint("warehouseID", id))
			return nil, types.ErrWarehouseNotFound
		}
		s.logger.Error("failed to fetch warehouse", zap.Error(err))
		return nil, types.ErrFailedToFetchWarehouse
	}

	return types.ToWarehouseDTO(*warehouse), nil
}

func (s *warehouseService) GetAllWarehouses(filter *types.WarehouseFilter) ([]types.WarehouseDTO, error) {
	warehouses, err := s.repo.GetAllWarehouses(filter)
	if err != nil {
		s.logger.Error("failed to fetch warehouses", zap.Error(err))
		return nil, types.ErrFailedToFetchWarehouses
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
		s.logger.Error("failed to fetch warehouses", zap.Error(err))
		return nil, types.ErrFailedToFetchWarehouses
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
		s.logger.Error("failed to update warehouse", zap.Error(err))
		return nil, types.ErrFailedUpdateWarehouse
	}

	updatedWarehouse, err := s.repo.GetWarehouseByID(id)
	if err != nil {
		s.logger.Error("failed to fetch updated warehouse", zap.Error(err))
		return nil, types.ErrFailedToFetchWarehouse
	}

	return types.ToWarehouseDTO(*updatedWarehouse), nil
}

func (s *warehouseService) DeleteWarehouse(id uint) error {
	if err := s.repo.DeleteWarehouse(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("warehouse not found", zap.Uint("warehouseID", id))
			return types.ErrWarehouseNotFound
		}
		return types.ErrFailedDeleteWarehouse
	}

	return nil
}
