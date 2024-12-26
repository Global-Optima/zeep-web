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
	GetInventoryLevels(filter *types.GetInventoryLevelsFilterQuery) (*types.InventoryLevelsResponse, error)
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
	var createdStockMaterials []data.StockMaterial
	if len(req.NewItems) > 0 {
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

func (s *inventoryService) GetInventoryLevels(filter *types.GetInventoryLevelsFilterQuery) (*types.InventoryLevelsResponse, error) {
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

func (s *inventoryService) assembleDeliveries(
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

func (s *inventoryService) createAndRegisterNewStockMaterials(supplierID uint, items []types.NewInventoryItem) ([]data.StockMaterial, error) {
	newStockMaterials := []data.StockMaterial{}

	for _, item := range items {

		expirationPeriod := 1095
		if item.ExpirationInDays != nil {
			expirationPeriod = *item.ExpirationInDays
		}

		newStockMaterial := data.StockMaterial{
			Name:                   item.Name,
			Description:            *item.Description,
			SafetyStock:            item.SafetyStock,
			ExpirationFlag:         item.ExpirationFlag,
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
		if err := s.packageRepo.CreatePackage(packageData); err != nil {
			return nil, fmt.Errorf("failed to create package for stock material %s: %w", item.Name, err)
		}

		newStockMaterials = append(newStockMaterials, newStockMaterial)
	}

	return newStockMaterials, nil
}

func (s *inventoryService) loadExistingStockMaterials(items []types.ExistingInventoryItem) (map[uint]*data.StockMaterial, error) {
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
