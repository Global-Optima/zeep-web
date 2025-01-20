package warehouseStock

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialPackage"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type WarehouseStockService interface {
	ReceiveInventory(req types.ReceiveInventoryRequest) error
	TransferInventory(req types.TransferInventoryRequest) error
	GetDeliveries(warehouseID *uint, startDate, endDate *time.Time) ([]types.DeliveryResponse, error)

	AddWarehouseStockMaterial(req types.AdjustWarehouseStock) error
	AddWarehouseStocks(warehouseID uint, req []types.AddWarehouseStockMaterial) error
	DeductFromStock(req types.AdjustWarehouseStock) error
	GetStock(query *types.GetWarehouseStockFilterQuery) ([]types.WarehouseStockResponse, error)
	GetStockMaterialDetails(stockMaterialID, warehouseID uint) (*types.WarehouseStockMaterialDetailsDTO, error)
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

func (s *warehouseStockService) ReceiveInventory(req types.ReceiveInventoryRequest) error {
	var createdStockMaterials []data.StockMaterial
	if len(req.NewItems) > 0 && req.NewItems != nil {
		var err error
		createdStockMaterials, err = s.createAndRegisterNewStockMaterials(req.SupplierID, req.NewItems)
		if err != nil {
			return fmt.Errorf("failed to create and register new stock materials: %w", err)
		}
	}

	existingStockMaterials, err := s.loadExistingStockMaterials(req.ExistingItems)
	if err != nil {
		return fmt.Errorf("failed to load existing stock materials: %w", err)
	}

	deliveries, err := s.assembleDeliveries(req, existingStockMaterials, createdStockMaterials)
	if err != nil {
		return fmt.Errorf("failed to assemble deliveries: %w", err)
	}

	if err := s.repo.RecordDeliveriesAndUpdateStock(deliveries, req.WarehouseID); err != nil {
		return fmt.Errorf("failed to log incoming inventory: %w", err)
	}

	return nil
}

func (s *warehouseStockService) assembleDeliveries(
	req types.ReceiveInventoryRequest,
	existingStockMaterials map[uint]*data.StockMaterial,
	newStockMaterials []data.StockMaterial,
) ([]data.SupplierWarehouseDelivery, error) {
	deliveries := []data.SupplierWarehouseDelivery{}
	newStockMaterialMap := make(map[string]*data.StockMaterial)

	for _, stockMaterial := range newStockMaterials {
		newStockMaterialMap[stockMaterial.Name] = &stockMaterial
	}

	for _, item := range req.ExistingItems {
		stockMaterial, found := existingStockMaterials[item.StockMaterialID]
		if !found {
			return nil, fmt.Errorf("existing stock material not found: ID %d", item.StockMaterialID)
		}

		delivery := data.SupplierWarehouseDelivery{
			StockMaterialID: stockMaterial.ID,
			SupplierID:      req.SupplierID,
			WarehouseID:     req.WarehouseID,
			Barcode:         stockMaterial.Barcode,
			Quantity:        item.Quantity,
			DeliveryDate:    time.Now(),
			ExpirationDate:  time.Now().AddDate(0, 0, stockMaterial.ExpirationPeriodInDays),
		}
		deliveries = append(deliveries, delivery)
	}

	for _, item := range req.NewItems {
		stockMaterial, found := newStockMaterialMap[item.Name]
		if !found {
			return nil, fmt.Errorf("new stock material not found: Name %s", item.Name)
		}

		delivery := data.SupplierWarehouseDelivery{
			StockMaterialID: stockMaterial.ID,
			SupplierID:      req.SupplierID,
			WarehouseID:     req.WarehouseID,
			Barcode:         stockMaterial.Barcode,
			Quantity:        item.Quantity,
			DeliveryDate:    time.Now(),
			ExpirationDate:  time.Now().AddDate(0, 0, stockMaterial.ExpirationPeriodInDays),
		}
		deliveries = append(deliveries, delivery)
	}

	return deliveries, nil
}

func (s *warehouseStockService) createAndRegisterNewStockMaterials(supplierID uint, items []types.NewWarehouseStockMaterial) ([]data.StockMaterial, error) {
	newStockMaterials := []data.StockMaterial{}

	for _, item := range items {
		expirationPeriod := 1095
		if item.ExpirationInDays != nil {
			expirationPeriod = *item.ExpirationInDays
		}

		newStockMaterial := data.StockMaterial{
			Name:                   item.Name,
			Description:            item.Description,
			SafetyStock:            item.SafetyStock,
			UnitID:                 item.UnitID,
			CategoryID:             item.CategoryID,
			ExpirationPeriodInDays: expirationPeriod,
			IngredientID:           item.IngredientID, // Linking to ingredient
			IsActive:               true,
		}

		if err := s.stockMaterialRepo.CreateStockMaterial(&newStockMaterial); err != nil {
			return nil, fmt.Errorf("failed to create stock material %s: %w", item.Name, err)
		}

		if err := s.ensureSupplierMaterialAssociation(supplierID, newStockMaterial.ID); err != nil {
			return nil, fmt.Errorf("failed to associate supplier %d with stock material %d: %w", supplierID, newStockMaterial.ID, err)
		}

		pkg := data.StockMaterialPackage{
			Size:   item.Package.Size,
			UnitID: item.Package.UnitID,
		}

		packageData := types.ValidatePackage(newStockMaterial.ID, pkg)
		if packageData == nil {
			return nil, fmt.Errorf("invalid package data for stock material %s", item.Name)
		}
		if err := s.packageRepo.Create(packageData); err != nil {
			return nil, fmt.Errorf("failed to create package for stock material %s: %w", item.Name, err)
		}

		newStockMaterials = append(newStockMaterials, newStockMaterial)
	}

	return newStockMaterials, nil
}

func (s *warehouseStockService) loadExistingStockMaterials(items []types.ExistingWarehouseStockMaterial) (map[uint]*data.StockMaterial, error) {
	stockMaterialIDs := []uint{}
	for _, item := range items {
		stockMaterialIDs = append(stockMaterialIDs, item.StockMaterialID)
	}

	stockMaterials, err := s.stockMaterialRepo.GetStockMaterialsByIDs(stockMaterialIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock materials: %w", err)
	}

	existingStockMaterials := make(map[uint]*data.StockMaterial)
	for _, stockMaterial := range stockMaterials {
		existingStockMaterials[stockMaterial.ID] = &stockMaterial
	}

	return existingStockMaterials, nil
}

func (s *warehouseStockService) ensureSupplierMaterialAssociation(supplierID, stockMaterialID uint) error {
	exists, err := s.repo.SupplierMaterialExists(supplierID, stockMaterialID)
	if err != nil {
		return fmt.Errorf("failed to check supplier-material association: %w", err)
	}

	if !exists {
		association := data.SupplierMaterial{
			SupplierID:      supplierID,
			StockMaterialID: stockMaterialID,
		}

		if err := s.repo.CreateSupplierMaterial(&association); err != nil {
			return fmt.Errorf("failed to create supplier-material association: %w", err)
		}
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

func (s *warehouseStockService) GetDeliveries(warehouseID *uint, startDate, endDate *time.Time) ([]types.DeliveryResponse, error) {
	deliveries, err := s.repo.GetDeliveries(warehouseID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deliveries: %w", err)
	}

	return types.DeliveriesToDeliveryResponses(deliveries), nil
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
		if stock.StockMaterial.Package == nil {
			return nil, fmt.Errorf("package measures not found for StockMaterialID %d", stock.StockMaterialID)
		}

		packageMeasures, err := utils.ReturnPackageMeasureForStockMaterialWithQuantity(stock.StockMaterial, stock.TotalQuantity)
		if err != nil {
			return nil, err
		}

		responses[i] = types.ToWarehouseStockResponse(stock, packageMeasures)
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

func (s *warehouseStockService) GetStockMaterialDetails(stockMaterialID, warehouseID uint) (*types.WarehouseStockMaterialDetailsDTO, error) {
	aggregatedStock, deliveries, err := s.repo.GetWarehouseStockMaterialDetails(stockMaterialID, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock material details: %w", err)
	}

	packageMeasure, err := utils.ReturnPackageMeasureForStockMaterialWithQuantity(
		aggregatedStock.StockMaterial,
		aggregatedStock.TotalQuantity,
	)
	if err != nil {
		return nil, err
	}

	details := types.ToStockMaterialDetails(*aggregatedStock, packageMeasure, deliveries)

	return &details, nil
}

func (s *warehouseStockService) UpdateStock(warehouseID, stockMaterialID uint, dto types.UpdateWarehouseStockDTO) error {
	stock, err := s.repo.GetWarehouseStockByID(warehouseID, stockMaterialID)
	if err != nil {
		return fmt.Errorf("failed to fetch warehouse stock: %w", err)
	}

	if err := s.repo.UpdateExpirationDate(stock.StockMaterialID, stock.WarehouseID, dto.ExpirationDate); err != nil {
		return fmt.Errorf("failed to update expiration date: %w", err)
	}

	if err := s.repo.UpdateStockQuantity(stock.ID, dto.Quantity); err != nil {
		return fmt.Errorf("failed to update stock quantity: %w", err)
	}

	return nil
}
