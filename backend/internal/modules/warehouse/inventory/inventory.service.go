package inventory

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/inventory/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/sku"
)

type InventoryService interface {
	ReceiveInventory(req types.ReceiveInventoryRequest) error
	TransferInventory(req types.TransferInventoryRequest) error
	GetInventoryLevels(warehouseID uint) ([]types.InventoryItem, error)
	PickupStock(req types.PickupRequest) error
	GetExpiringItems(warehouseID uint, thresholdDays int) ([]types.UpcomingExpirationResponse, error)
	ExtendExpiration(req types.ExtendExpirationRequest) error
	GetDeliveries(warehouseID *uint, startDate, endDate *time.Time) ([]types.DeliveryResponse, error)
}

type inventoryService struct {
	repo        InventoryRepository
	skuRepo     sku.SKURepository
	barcodeRepo barcode.BarcodeRepository
	packageRepo PackageRepository
}

func NewInventoryService(repo InventoryRepository, skuRepo sku.SKURepository, barcodeRepo barcode.BarcodeRepository, packageRepo PackageRepository) InventoryService {
	return &inventoryService{
		repo:        repo,
		skuRepo:     skuRepo,
		barcodeRepo: barcodeRepo,
		packageRepo: packageRepo,
	}
}

func (s *inventoryService) ReceiveInventory(req types.ReceiveInventoryRequest) error {
	existingSKUs, newSKUItems, err := s.separateExistingAndNewSKUs(req.Items)
	if err != nil {
		return err
	}

	createdSKUs, err := s.createAndRegisterNewSKUs(req.SupplierID, newSKUItems)
	if err != nil {
		return err
	}

	fullItems := s.mergeSKUItems(req.Items, createdSKUs)

	deliveries, err := s.createDeliveries(req, existingSKUs, createdSKUs, fullItems)
	if err != nil {
		return err
	}

	if err := s.repo.LogAndUpdateStock(deliveries, req.WarehouseID); err != nil {
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

func (s *inventoryService) GetInventoryLevels(warehouseID uint) ([]types.InventoryItem, error) {
	stocks, err := s.repo.GetInventoryLevels(warehouseID)
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
	var delivery data.Delivery
	if err := s.repo.GetDeliveryByID(req.DeliveryID, &delivery); err != nil {
		return fmt.Errorf("failed to fetch delivery: %w", err)
	}

	if err := types.ValidateExpirationDate(req.NewExpirationDate, delivery.ExpirationDate); err != nil {
		return err
	}

	if err := s.repo.ExtendExpiration(req.DeliveryID, req.NewExpirationDate); err != nil {
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
func (s *inventoryService) createDeliveries(
	req types.ReceiveInventoryRequest,
	existingSKUs map[uint]*data.StockMaterial,
	newSKUs []data.StockMaterial,
	fullItems []types.InventoryItem,
) ([]data.Delivery, error) {
	deliveries := []data.Delivery{}
	newSKUMap := make(map[uint]*data.StockMaterial)

	for _, sku := range newSKUs {
		newSKUMap[sku.ID] = &sku
	}

	for _, item := range fullItems {
		var sku *data.StockMaterial
		var found bool

		if item.SKU_ID != 0 {
			sku, found = existingSKUs[item.SKU_ID]
			if !found {
				sku = newSKUMap[item.SKU_ID]
				found = sku != nil
			}
		}

		if !found {
			return nil, fmt.Errorf("failed to find SKU for item: %v", item)
		}

		if err := s.ensureSupplierMaterialAssociation(req.SupplierID, sku.ID); err != nil {
			return nil, err
		}

		expirationPeriod := sku.ExpirationPeriodInDays
		if item.Expiration != nil {
			expirationPeriod = *item.Expiration
		}

		delivery := data.Delivery{
			StockMaterialID: sku.ID,
			SupplierID:      req.SupplierID,
			WarehouseID:     req.WarehouseID,
			Barcode:         sku.Barcode,
			Quantity:        item.Quantity,
			DeliveryDate:    time.Now(),
			ExpirationDate:  time.Now().AddDate(0, 0, expirationPeriod),
		}
		deliveries = append(deliveries, delivery)
	}

	return deliveries, nil
}

func (s *inventoryService) createAndRegisterNewSKUs(supplierID uint, items []types.InventoryItem) ([]data.StockMaterial, error) {
	newSKUs := []data.StockMaterial{}

	for _, item := range items {
		expirationPeriod := 1095
		if item.Expiration != nil {
			expirationPeriod = *item.Expiration
		}

		newSKU := data.StockMaterial{
			Name:                   *item.Name,
			Description:            *item.Description,
			SafetyStock:            *item.SafetyStock,
			ExpirationFlag:         *item.ExpirationFlag,
			UnitID:                 *item.UnitID,
			Category:               *item.Category,
			ExpirationPeriodInDays: expirationPeriod,
			IsActive:               true,
		}

		if err := s.skuRepo.CreateSKU(&newSKU); err != nil {
			return nil, fmt.Errorf("failed to create new SKU %s: %w", *item.Name, err)
		}

		if err := s.ensureSupplierMaterialAssociation(supplierID, newSKU.ID); err != nil {
			return nil, fmt.Errorf("failed to associate supplier %d with SKU %d: %w", supplierID, newSKU.ID, err)
		}

		barcode, err := s.barcodeRepo.GenerateAndAssignBarcode(newSKU.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate barcode for SKU %d: %w", newSKU.ID, err)
		}
		newSKU.Barcode = barcode

		newPackage := data.Package{
			StockMaterialID: newSKU.ID,
			PackageSize:     item.Package.PackageSize,
			PackageUnitID:   item.Package.PackageUnitID,
		}
		if err := s.packageRepo.CreatePackage(&newPackage); err != nil {
			return nil, fmt.Errorf("failed to create package for SKU %d: %w", newSKU.ID, err)
		}

		newSKUs = append(newSKUs, newSKU)
	}

	return newSKUs, nil
}

func (s *inventoryService) separateExistingAndNewSKUs(items []types.InventoryItem) (map[uint]*data.StockMaterial, []types.InventoryItem, error) {
	existingSKUs := make(map[uint]*data.StockMaterial)
	newSKUItems := []types.InventoryItem{}
	skuIDs := []uint{}

	for _, item := range items {
		if item.SKU_ID != 0 {
			skuIDs = append(skuIDs, item.SKU_ID)
		} else {
			newSKUItems = append(newSKUItems, item)
		}
	}

	if len(skuIDs) > 0 {
		skus, err := s.skuRepo.GetSKUsByIDs(skuIDs)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch existing SKUs: %w", err)
		}
		for _, sku := range skus {
			existingSKUs[sku.ID] = &sku
		}
	}

	return existingSKUs, newSKUItems, nil
}

func (s *inventoryService) mergeSKUItems(items []types.InventoryItem, createdSKUs []data.StockMaterial) []types.InventoryItem {
	createdSKUMap := make(map[string]uint)
	for _, sku := range createdSKUs {
		createdSKUMap[sku.Name] = sku.ID
	}

	fullItems := make([]types.InventoryItem, len(items))
	for i, item := range items {
		if item.SKU_ID == 0 {
			if id, found := createdSKUMap[*item.Name]; found {
				item.SKU_ID = id
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
