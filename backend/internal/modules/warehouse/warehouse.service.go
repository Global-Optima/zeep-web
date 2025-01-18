package warehouse

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type WarehouseService interface {
	AssignStoreToWarehouse(req types.AssignStoreToWarehouseRequest) error
	ReassignStore(storeID uint, req types.ReassignStoreRequest) error
	GetAllStoresByWarehouse(warehouseID uint, pagination *utils.Pagination) ([]types.ListStoresResponse, error)

	CreateWarehouse(req types.CreateWarehouseDTO) (*types.WarehouseResponse, error)
	GetWarehouseByID(id uint) (*types.WarehouseResponse, error)
	GetAllWarehouses(pagination *utils.Pagination) ([]types.WarehouseResponse, error)
	UpdateWarehouse(id uint, req types.UpdateWarehouseDTO) (*types.WarehouseResponse, error)
	DeleteWarehouse(id uint) error

	AddToStock(req types.AdjustWarehouseStockRequest) error
	DeductFromStock(req types.AdjustWarehouseStockRequest) error
	GetStock(query *types.GetWarehouseStockFilterQuery) ([]types.WarehouseStockResponse, error)
	GetStockMaterialDetails(stockMaterialID, warehouseID uint) (*types.StockMaterialDetailsDTO, error)
	ResetStock(req types.ResetWarehouseStockRequest) error
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

func (s *warehouseService) CreateWarehouse(req types.CreateWarehouseDTO) (*types.WarehouseResponse, error) {
	facilityAddress := types.ToFacilityAddressModel(req.FacilityAddress)
	warehouse := types.ToWarehouseModel(req, 0)

	if err := s.repo.CreateWarehouse(&warehouse, &facilityAddress); err != nil {
		return nil, fmt.Errorf("failed to create warehouse: %w", err)
	}

	createdWarehouse, err := s.repo.GetWarehouseByID(warehouse.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created warehouse: %w", err)
	}

	return types.ToWarehouseResponse(*createdWarehouse), nil
}

func (s *warehouseService) GetWarehouseByID(id uint) (*types.WarehouseResponse, error) {
	warehouse, err := s.repo.GetWarehouseByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("warehouse with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch warehouse: %w", err)
	}

	return types.ToWarehouseResponse(*warehouse), nil
}

func (s *warehouseService) GetAllWarehouses(pagination *utils.Pagination) ([]types.WarehouseResponse, error) {
	warehouses, err := s.repo.GetAllWarehouses(pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all warehouses: %w", err)
	}

	responses := make([]types.WarehouseResponse, len(warehouses))
	for i, warehouse := range warehouses {
		responses[i] = *types.ToWarehouseResponse(warehouse)
	}

	return responses, nil
}

func (s *warehouseService) UpdateWarehouse(id uint, req types.UpdateWarehouseDTO) (*types.WarehouseResponse, error) {
	warehouse, err := s.repo.GetWarehouseByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("warehouse with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch warehouse: %w", err)
	}

	warehouse.Name = req.Name

	if err := s.repo.UpdateWarehouse(warehouse); err != nil {
		return nil, fmt.Errorf("failed to update warehouse: %w", err)
	}

	updatedWarehouse, err := s.repo.GetWarehouseByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated warehouse: %w", err)
	}

	return types.ToWarehouseResponse(*updatedWarehouse), nil
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

func (s *warehouseService) AddToStock(req types.AdjustWarehouseStockRequest) error {
	return s.repo.AddToWarehouseStock(req.WarehouseID, req.StockMaterialID, req.Quantity)
}

func (s *warehouseService) DeductFromStock(req types.AdjustWarehouseStockRequest) error {
	return s.repo.DeductFromWarehouseStock(req.WarehouseID, req.StockMaterialID, req.Quantity)
}

func (s *warehouseService) GetStock(query *types.GetWarehouseStockFilterQuery) ([]types.WarehouseStockResponse, error) {
	stocks, err := s.repo.GetWarehouseStock(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch warehouse stocks: %w", err)
	}

	responses := make([]types.WarehouseStockResponse, len(stocks))
	for i, stock := range stocks {
		if stock.StockMaterial.Package == nil {
			return nil, fmt.Errorf("package measures not found for StockMaterialID %d", stock.StockMaterialID)
		}

		packageMeasures, err := utils.ReturnPackageMeasureForStockMaterial(stock.StockMaterial, stock.TotalQuantity)
		if err != nil {
			return nil, err
		}

		responses[i] = types.WarehouseStockResponse{
			StockMaterial: types.StockMaterialResponse{
				ID:             stock.StockMaterialID,
				Name:           stock.StockMaterial.Name,
				Description:    stock.StockMaterial.Description,
				Category:       stock.StockMaterial.StockMaterialCategory.Name,
				SafetyStock:    stock.StockMaterial.SafetyStock,
				Barcode:        stock.StockMaterial.Barcode,
				PackageMeasure: packageMeasures,
			},
			TotalQuantity:          stock.TotalQuantity,
			EarliestExpirationDate: stock.EarliestExpirationDate,
		}
	}

	return responses, nil
}

func (s *warehouseService) GetStockMaterialDetails(stockMaterialID, warehouseID uint) (*types.StockMaterialDetailsDTO, error) {
	aggregatedStock, deliveries, err := s.repo.GetWarehouseStockMaterialDetails(stockMaterialID, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock material details: %w", err)
	}

	deliveriesDTO := make([]types.StockMaterialDeliveryDTO, len(deliveries))
	for i, delivery := range deliveries {
		deliveriesDTO[i] = types.StockMaterialDeliveryDTO{
			Supplier:     delivery.Supplier.Name,
			Quantity:     delivery.Quantity,
			DeliveryDate: delivery.DeliveryDate,
			ExpiresOn:    delivery.ExpirationDate,
		}
	}

	packageMeasure, err := utils.ReturnPackageMeasureForStockMaterial(
		aggregatedStock.StockMaterial,
		aggregatedStock.TotalQuantity,
	)
	if err != nil {
		return nil, err
	}

	details := &types.StockMaterialDetailsDTO{
		ID:                     aggregatedStock.StockMaterial.ID,
		Name:                   aggregatedStock.StockMaterial.Name,
		Description:            aggregatedStock.StockMaterial.Description,
		Category:               aggregatedStock.StockMaterial.StockMaterialCategory.Name,
		SafetyStock:            aggregatedStock.StockMaterial.SafetyStock,
		ExpirationFlag:         aggregatedStock.StockMaterial.ExpirationFlag,
		ExpirationInDays:       aggregatedStock.StockMaterial.ExpirationPeriodInDays,
		PackageMeasure:         packageMeasure,
		TotalQuantity:          aggregatedStock.TotalQuantity,
		EarliestExpirationDate: aggregatedStock.EarliestExpirationDate,
		Deliveries:             deliveriesDTO,
	}

	return details, nil
}

func (s *warehouseService) ResetStock(req types.ResetWarehouseStockRequest) error {
	stocks := []data.WarehouseStock{}
	for _, stock := range req.Stocks {
		stocks = append(stocks, data.WarehouseStock{
			WarehouseID:     req.WarehouseID,
			StockMaterialID: stock.StockMaterialID,
			Quantity:        stock.Quantity,
		})
	}
	return s.repo.ResetWarehouseStock(req.WarehouseID, stocks)
}
