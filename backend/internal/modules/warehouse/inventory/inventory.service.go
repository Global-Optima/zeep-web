package inventory

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/inventory/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
)

type InventoryService interface {
	ReceiveInventory(req types.ReceiveInventoryRequest) error
	TransferInventory(req types.TransferInventoryRequest) error
	GetInventoryLevels(filter *types.GetInventoryLevelsFilterQuery) ([]types.InventoryItem, error)
	PickupStock(req types.PickupRequest) error
	GetExpiringItems(warehouseID uint, thresholdDays int) ([]types.UpcomingExpirationResponse, error)
	ExtendExpiration(req types.ExtendExpirationRequest) error
	GetDeliveries(warehouseID *uint, startDate, endDate *time.Time) ([]types.DeliveryResponse, error)
}

type inventoryService struct {
	repo              InventoryRepository
	stockMaterialRepo stockMaterial.StockMaterialRepository
	barcodeRepo       barcode.BarcodeRepository
	packageRepo       PackageRepository
}

func NewInventoryService(repo InventoryRepository, stockMaterialRepo stockMaterial.StockMaterialRepository, barcodeRepo barcode.BarcodeRepository, packageRepo PackageRepository) InventoryService {
	return &inventoryService{
		repo:              repo,
		stockMaterialRepo: stockMaterialRepo,
		barcodeRepo:       barcodeRepo,
		packageRepo:       packageRepo,
	}
}

func (s *inventoryService) ReceiveInventory(req types.ReceiveInventoryRequest) error {
	existingStockMaterials, newStockMaterialItems, err := s.separateExistingAndNewStockMaterials(req.Items)
	if err != nil {
		return err
	}

	createdStockMaterials, err := s.createAndRegisterNewStockMaterials(req.SupplierID, newStockMaterialItems)
	if err != nil {
		return err
	}

	fullItems := s.mergeStockMaterialItems(req.Items, createdStockMaterials)

	deliveries, err := s.assembleDeliveries(req, existingStockMaterials, createdStockMaterials, fullItems)
	if err != nil {
		return err
	}

	if err := s.repo.RecordDeliveriesAndUpdateStock(deliveries, req.WarehouseID); err != nil {
		return fmt.Errorf("failed to log incoming inventory: %w", err)
	}

	return nil
}

func (s *inventoryService) TransferInventory(req types.TransferInventoryRequest) error {
	stockItems, err := s.repo.ConvertInventoryItemsToStockRequest(req.Items)
	if err != nil {
		return fmt.Errorf("failed to convert inventory items: %w", err)
	}

	if err := s.repo.TransferStock(req.SourceWarehouseID, req.TargetWarehouseID, stockItems); err != nil {
		return fmt.Errorf("failed to transfer stock: %w", err)
	}

	return nil
}

func (s *inventoryService) GetInventoryLevels(filter *types.GetInventoryLevelsFilterQuery) ([]types.InventoryItem, error) {
	stocks, err := s.repo.GetInventoryLevels(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch inventory levels: %w", err)
	}

	return types.StocksToInventoryItems(stocks), nil
}

func (s *inventoryService) PickupStock(req types.PickupRequest) error {
	stockItems, err := s.repo.ConvertInventoryItemsToStockRequest(req.Items)
	if err != nil {
		return fmt.Errorf("failed to convert inventory items: %w", err)
	}

	if err := s.repo.PickupStock(req.StoreWarehouseID, stockItems); err != nil {
		return fmt.Errorf("failed to handle store pickup: %w", err)
	}

	return nil
}

func (s *inventoryService) GetExpiringItems(warehouseID uint, thresholdDays int) ([]types.UpcomingExpirationResponse, error) {
	deliveries, err := s.repo.GetExpiringItems(warehouseID, thresholdDays)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expiring items: %w", err)
	}

	return types.ExpiringItemsToResponses(deliveries), nil
}

func (s *inventoryService) ExtendExpiration(req types.ExtendExpirationRequest) error {
	var delivery data.SupplierWarehouseDelivery
	if err := s.repo.GetDeliveryByID(req.DeliveryID, &delivery); err != nil {
		return fmt.Errorf("failed to fetch delivery: %w", err)
	}

	if err := types.ValidateExpirationDays(req.AddDays); err != nil {
		return err
	}

	newExpirationDate := delivery.ExpirationDate.AddDate(0, 0, req.AddDays)

	if err := s.repo.ExtendExpiration(req.DeliveryID, newExpirationDate); err != nil {
		return fmt.Errorf("failed to extend expiration date: %w", err)
	}

	return nil
}

func (s *inventoryService) GetDeliveries(warehouseID *uint, startDate, endDate *time.Time) ([]types.DeliveryResponse, error) {
	deliveries, err := s.repo.GetDeliveries(warehouseID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deliveries: %w", err)
	}

	return types.DeliveriesToDeliveryResponses(deliveries), nil
}

// Helper methods
func (s *inventoryService) assembleDeliveries(
	req types.ReceiveInventoryRequest,
	existingStockMaterials map[uint]*data.StockMaterial,
	newStockMaterials []data.StockMaterial,
	fullItems []types.InventoryItem,
) ([]data.SupplierWarehouseDelivery, error) {
	deliveries := []data.SupplierWarehouseDelivery{}
	newStockMaterialMap := make(map[uint]*data.StockMaterial)

	for _, stockMaterial := range newStockMaterials {
		newStockMaterialMap[stockMaterial.ID] = &stockMaterial
	}

	for _, item := range fullItems {
		var stockMaterial *data.StockMaterial
		var found bool

		if item.StockMaterialID != 0 {
			stockMaterial, found = existingStockMaterials[item.StockMaterialID]
			if !found {
				stockMaterial = newStockMaterialMap[item.StockMaterialID]
				found = stockMaterial != nil
			}
		}

		if !found {
			return nil, fmt.Errorf("failed to find StockMaterial for item: %v", item)
		}

		if err := s.ensureSupplierMaterialAssociation(req.SupplierID, stockMaterial.ID); err != nil {
			return nil, err
		}

		expirationPeriod := stockMaterial.ExpirationPeriodInDays
		if item.Expiration != nil {
			expirationPeriod = *item.Expiration
		}

		delivery := data.SupplierWarehouseDelivery{
			StockMaterialID: stockMaterial.ID,
			SupplierID:      req.SupplierID,
			WarehouseID:     req.WarehouseID,
			Barcode:         stockMaterial.Barcode,
			Quantity:        item.Quantity,
			DeliveryDate:    time.Now(),
			ExpirationDate:  time.Now().AddDate(0, 0, expirationPeriod),
		}
		deliveries = append(deliveries, delivery)
	}

	return deliveries, nil
}

func (s *inventoryService) createAndRegisterNewStockMaterials(supplierID uint, items []types.InventoryItem) ([]data.StockMaterial, error) {
	newStockMaterials := []data.StockMaterial{}

	for _, item := range items {
		expirationPeriod := 1095
		if item.Expiration != nil {
			expirationPeriod = *item.Expiration
		}

		newStockMaterial := data.StockMaterial{
			Name:                   *item.Name,
			Description:            *item.Description,
			SafetyStock:            *item.SafetyStock,
			ExpirationFlag:         *item.ExpirationFlag,
			UnitID:                 *item.UnitID,
			Category:               *item.Category,
			ExpirationPeriodInDays: expirationPeriod,
			IsActive:               true,
		}

		if err := s.stockMaterialRepo.CreateStockMaterial(&newStockMaterial); err != nil {
			return nil, fmt.Errorf("failed to create new StockMaterial %s: %w", *item.Name, err)
		}
		item.StockMaterialID = newStockMaterial.ID

		if err := s.ensureSupplierMaterialAssociation(supplierID, newStockMaterial.ID); err != nil {
			return nil, fmt.Errorf("failed to associate supplier %d with StockMaterial %d: %w", supplierID, newStockMaterial.ID, err)
		}

		barcode, err := s.barcodeRepo.GenerateAndAssignBarcode(newStockMaterial.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate barcode for StockMaterial %d: %w", newStockMaterial.ID, err)
		}
		newStockMaterial.Barcode = barcode

		newPackage := types.ValidatePackage(item)
		if newPackage == nil {
			return nil, fmt.Errorf("package data not provided for new stockMaterial")
		}

		if err := s.packageRepo.CreatePackage(newPackage); err != nil {
			return nil, fmt.Errorf("failed to create package for StockMaterial %d: %w", newStockMaterial.ID, err)
		}

		newStockMaterials = append(newStockMaterials, newStockMaterial)
	}

	return newStockMaterials, nil
}

func (s *inventoryService) separateExistingAndNewStockMaterials(items []types.InventoryItem) (map[uint]*data.StockMaterial, []types.InventoryItem, error) {
	existingStockMaterials := make(map[uint]*data.StockMaterial)
	newStockMaterialItems := []types.InventoryItem{}
	stockMaterialIDs := []uint{}

	for _, item := range items {
		if item.StockMaterialID != 0 {
			stockMaterialIDs = append(stockMaterialIDs, item.StockMaterialID)
		} else {
			newStockMaterialItems = append(newStockMaterialItems, item)
		}
	}

	if len(stockMaterialIDs) > 0 {
		stockMaterials, err := s.stockMaterialRepo.GetStockMaterialsByIDs(stockMaterialIDs)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch existing StockMaterials: %w", err)
		}
		for _, stockMaterial := range stockMaterials {
			existingStockMaterials[stockMaterial.ID] = &stockMaterial
		}
	}

	return existingStockMaterials, newStockMaterialItems, nil
}

func (s *inventoryService) mergeStockMaterialItems(items []types.InventoryItem, createdStockMaterials []data.StockMaterial) []types.InventoryItem {
	createdStockMaterialMap := make(map[string]uint)
	for _, stockMaterial := range createdStockMaterials {
		createdStockMaterialMap[stockMaterial.Name] = stockMaterial.ID
	}

	fullItems := make([]types.InventoryItem, len(items))
	for i, item := range items {
		if item.StockMaterialID == 0 {
			if id, found := createdStockMaterialMap[*item.Name]; found {
				item.StockMaterialID = id
			}
		}
		fullItems[i] = item
	}

	return fullItems
}

func (s *inventoryService) ensureSupplierMaterialAssociation(supplierID, stockMaterialID uint) error {
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
