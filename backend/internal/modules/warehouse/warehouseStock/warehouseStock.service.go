package warehouseStock

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialPackage"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock/types"
)

type WarehouseStockService interface {
	ReceiveInventory(warehouseID uint, req types.ReceiveWarehouseDelivery) error
	TransferInventory(req types.TransferInventoryRequest) error
	GetDeliveries(filter types.WarehouseDeliveryFilter) ([]types.WarehouseDeliveryDTO, error)
	GetDeliveryByID(id uint) (*types.WarehouseDeliveryDTO, error)

	AddWarehouseStockMaterial(req types.AdjustWarehouseStock) error
	AddWarehouseStocks(warehouseID uint, req []types.AddWarehouseStockMaterial) error
	DeductFromStock(req types.AdjustWarehouseStock) error
	GetStock(query *types.GetWarehouseStockFilterQuery) ([]types.WarehouseStockResponse, error)
	GetStockMaterialDetails(stockMaterialID, warehouseID uint) (*types.WarehouseStockResponse, error)
	UpdateStock(warehouseID, stockMaterialID uint, dto types.UpdateWarehouseStockDTO) error
}

type warehouseStockService struct {
	repo              WarehouseStockRepository
	stockMaterialRepo stockMaterial.StockMaterialRepository
	barcodeRepo       barcode.BarcodeRepository
	packageRepo       stockMaterialPackage.StockMaterialPackageRepository
}

func NewWarehouseStockService(repo WarehouseStockRepository, stockMaterialRepo stockMaterial.StockMaterialRepository, barcodeRepo barcode.BarcodeRepository, packageRepo stockMaterialPackage.StockMaterialPackageRepository) WarehouseStockService {
	return &warehouseStockService{
		repo:              repo,
		stockMaterialRepo: stockMaterialRepo,
		barcodeRepo:       barcodeRepo,
		packageRepo:       packageRepo,
	}
}

func (s *warehouseStockService) ReceiveInventory(warehouseID uint, req types.ReceiveWarehouseDelivery) error {
	stockMaterialIDs := make([]uint, len(req.Materials))
	for i, material := range req.Materials {
		stockMaterialIDs[i] = material.StockMaterialID
	}
	stockMaterials, err := s.stockMaterialRepo.GetStockMaterialsByIDs(stockMaterialIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch stock materials: %w", err)
	}

	stockMaterialMap := make(map[uint]*data.StockMaterial)
	for _, stockMaterial := range stockMaterials {
		stockMaterialMap[stockMaterial.ID] = &stockMaterial
	}

	delivery := data.SupplierWarehouseDelivery{
		SupplierID:   req.SupplierID,
		WarehouseID:  warehouseID,
		DeliveryDate: time.Now(),
	}

	materials := make([]data.SupplierWarehouseDeliveryMaterial, len(req.Materials))

	for i, material := range req.Materials {
		stockMaterial, exists := stockMaterialMap[material.StockMaterialID]
		if !exists {
			return fmt.Errorf("stock material with ID %d not found", material.StockMaterialID)
		}

		packageFound := false
		for _, pkg := range stockMaterial.Packages {
			if pkg.ID == material.PackageID {
				packageFound = true
				break
			}
		}
		if !packageFound {
			return fmt.Errorf("package with ID %d not found for stock material ID %d", material.PackageID, material.StockMaterialID)
		}

		var packageToUse *data.StockMaterialPackage
		for _, pkg := range stockMaterial.Packages {
			if pkg.ID == material.PackageID {
				packageToUse = &pkg
				break
			}
		}
		if packageToUse == nil {
			return fmt.Errorf("failed to retrieve package details for package ID %d", material.PackageID)
		}

		materials[i] = data.SupplierWarehouseDeliveryMaterial{
			StockMaterialID: material.StockMaterialID,
			PackageID:       material.PackageID,
			Barcode:         stockMaterial.Barcode,
			Quantity:        material.Quantity,
			ExpirationDate:  time.Now().AddDate(0, 0, stockMaterial.ExpirationPeriodInDays),
		}
	}

	err = s.repo.RecordDeliveriesAndUpdateStock(delivery, materials, warehouseID)
	if err != nil {
		return fmt.Errorf("failed to receive inventory: %w", err)
	}

	return nil
}

func (s *warehouseStockService) TransferInventory(req types.TransferInventoryRequest) error {
	stockItems, err := s.repo.ConvertInventoryItemsToStockRequest(req.Items)
	if err != nil {
		return fmt.Errorf("failed to convert inventory items: %w", err)
	}

	if err := s.repo.TransferStock(req.SourceWarehouseID, req.TargetWarehouseID, stockItems); err != nil {
		return fmt.Errorf("failed to transfer stock: %w", err)
	}

	return nil
}

func (s *warehouseStockService) GetDeliveries(filter types.WarehouseDeliveryFilter) ([]types.WarehouseDeliveryDTO, error) {
	deliveries, err := s.repo.GetDeliveries(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deliveries: %w", err)
	}

	return types.DeliveriesToDeliveryResponses(deliveries), nil
}

func (s *warehouseStockService) GetDeliveryByID(id uint) (*types.WarehouseDeliveryDTO, error) {
	var delivery data.SupplierWarehouseDelivery
	err := s.repo.GetDeliveryByID(id, &delivery)
	if err != nil {
		return nil, err
	}

	response := types.ToDeliveryResponse(delivery)
	return &response, nil
}

func (s *warehouseStockService) AddWarehouseStockMaterial(req types.AdjustWarehouseStock) error {
	return s.repo.AddToWarehouseStock(req.WarehouseID, req.StockMaterialID, req.Quantity)
}

func (s *warehouseStockService) DeductFromStock(req types.AdjustWarehouseStock) error {
	return s.repo.DeductFromWarehouseStock(req.WarehouseID, req.StockMaterialID, req.Quantity)
}

func (s *warehouseStockService) GetStock(query *types.GetWarehouseStockFilterQuery) ([]types.WarehouseStockResponse, error) {
	stocks, err := s.repo.GetWarehouseStock(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch warehouse stocks: %w", err)
	}

	responses := make([]types.WarehouseStockResponse, len(stocks))
	for i, stock := range stocks {
		responses[i] = types.ToWarehouseStockResponse(stock)
	}

	return responses, nil
}

func (s *warehouseStockService) AddWarehouseStocks(warehouseID uint, req []types.AddWarehouseStockMaterial) error {
	if len(req) == 0 {
		return fmt.Errorf("stocks cannot be empty")
	}

	var stocks []data.WarehouseStock
	for _, dto := range req {
		stocks = append(stocks, data.WarehouseStock{
			WarehouseID:     warehouseID,
			StockMaterialID: dto.StockMaterialID,
			Quantity:        dto.Quantity,
		})
	}

	return s.repo.AddWarehouseStocks(warehouseID, stocks)
}

func (s *warehouseStockService) GetStockMaterialDetails(stockMaterialID, warehouseID uint) (*types.WarehouseStockResponse, error) {
	aggregatedStock, err := s.repo.GetWarehouseStockMaterialDetails(stockMaterialID, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock material details: %w", err)
	}

	details := types.ToWarehouseStockResponse(*aggregatedStock)

	return &details, nil
}

func (s *warehouseStockService) UpdateStock(warehouseID, stockMaterialID uint, dto types.UpdateWarehouseStockDTO) error {
	// Validate input
	if dto.Quantity == nil && dto.ExpirationDate == nil {
		return fmt.Errorf("nothing to update")
	}

	// Update quantity if provided
	if dto.Quantity != nil {
		err := s.repo.UpdateStockQuantity(stockMaterialID, warehouseID, *dto.Quantity)
		if err != nil {
			return fmt.Errorf("failed to update stock quantity: %w", err)
		}
	}

	// Update expiration date if provided
	if dto.ExpirationDate != nil {
		err := s.repo.UpdateExpirationDate(stockMaterialID, warehouseID, *dto.ExpirationDate)
		if err != nil {
			return fmt.Errorf("failed to update expiration date: %w", err)
		}
	}

	return nil
}
