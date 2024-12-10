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
}

func NewInventoryService(repo InventoryRepository, skuRepo sku.SKURepository, barcodeRepo barcode.BarcodeRepository) InventoryService {
	return &inventoryService{
		repo:        repo,
		skuRepo:     skuRepo,
		barcodeRepo: barcodeRepo,
	}
}

func (s *inventoryService) ReceiveInventory(req types.ReceiveInventoryRequest) error {
	existingSKUs := make(map[uint]*data.SKU)
	newSKUItems := []types.InventoryItem{}
	deliveries := make([]data.Delivery, len(req.Items))

	skuIDs := make([]uint, 0)
	for _, item := range req.Items {
		if item.SKU_ID != 0 {
			skuIDs = append(skuIDs, item.SKU_ID)
		} else {
			newSKUItems = append(newSKUItems, item)
		}
	}

	if len(skuIDs) > 0 {
		skus, err := s.skuRepo.GetSKUsByIDs(skuIDs)
		if err != nil {
			return fmt.Errorf("failed to fetch existing SKUs: %w", err)
		}
		for _, sku := range skus {
			existingSKUs[sku.ID] = &sku
		}
	}

	newSKUs := []data.SKU{}
	for _, item := range newSKUItems {
		newSKUs = append(newSKUs, data.SKU{
			Name:             *item.Name,
			Description:      *item.Description,
			SafetyStock:      *item.SafetyStock,
			ExpirationFlag:   *item.ExpirationFlag,
			Quantity:         item.Quantity,
			SupplierID:       req.SupplierID,
			UnitID:           *item.UnitID,
			Category:         *item.Category,
			ExpirationPeriod: 1095,
			IsActive:         true,
		})
	}

	if len(newSKUs) > 0 {
		if err := s.skuRepo.CreateSKUs(newSKUs); err != nil {
			return fmt.Errorf("failed to create new SKUs: %w", err)
		}

		for i := range newSKUs {
			barcode, err := s.barcodeRepo.GenerateAndAssignBarcode(newSKUs[i].ID)
			if err != nil {
				return fmt.Errorf("failed to generate barcode for SKU %d: %w", newSKUs[i].ID, err)
			}
			newSKUs[i].Barcode = barcode
		}
	}

	for i, item := range req.Items {
		if item.SKU_ID != 0 {

			sku := existingSKUs[item.SKU_ID]
			deliveries[i] = data.Delivery{
				SKU_ID:         sku.ID,
				Source:         req.SupplierID,
				Target:         req.WarehouseID,
				Barcode:        sku.Barcode,
				Quantity:       item.Quantity,
				DeliveryDate:   time.Now(),
				ExpirationDate: time.Now().AddDate(0, 0, 1095),
			}
		} else {

			for _, sku := range newSKUs {
				if sku.Name == *item.Name {
					deliveries[i] = data.Delivery{
						SKU_ID:         sku.ID,
						Source:         req.SupplierID,
						Target:         req.WarehouseID,
						Barcode:        sku.Barcode,
						Quantity:       item.Quantity,
						DeliveryDate:   time.Now(),
						ExpirationDate: time.Now().AddDate(0, 0, 1095),
					}
					break
				}
			}
		}
	}

	if err := s.repo.LogIncomingInventory(deliveries); err != nil {
		return fmt.Errorf("failed to log incoming inventory: %w", err)
	}

	stockItems, err := s.repo.ConvertInventoryItemsToStockRequest(req.Items)
	if err != nil {
		return fmt.Errorf("failed to convert inventory items: %w", err)
	}

	if err := s.repo.UpdateStockLevels(req.WarehouseID, stockItems); err != nil {
		return fmt.Errorf("failed to update stock levels: %w", err)
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

	response := make([]types.InventoryItem, len(stocks))
	for i, stock := range stocks {
		response[i] = types.InventoryItem{
			SKU_ID:   stock.IngredientID,
			Quantity: stock.Quantity,
		}
	}

	return response, nil
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

	response := make([]types.UpcomingExpirationResponse, len(deliveries))
	for i, delivery := range deliveries {
		response[i] = types.UpcomingExpirationResponse{
			SKU_ID:         delivery.SKU_ID,
			Name:           delivery.SKU.Name,
			ExpirationDate: delivery.ExpirationDate,
			Quantity:       delivery.Quantity,
		}
	}

	return response, nil
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

	response := make([]types.DeliveryResponse, len(deliveries))
	for i, delivery := range deliveries {
		response[i] = types.DeliveryResponse{
			ID:             delivery.ID,
			SKU_ID:         delivery.SKU_ID,
			Source:         delivery.Source,
			Target:         delivery.Target,
			Barcode:        delivery.Barcode,
			Quantity:       delivery.Quantity,
			DeliveryDate:   delivery.DeliveryDate,
			ExpirationDate: delivery.ExpirationDate,
		}
	}

	return response, nil
}

func (s *inventoryService) logAudit(action string, skuID uint, deliveryID *uint, quantity float64, unit string, userID uint) error {
	auditLog := data.AuditLog{
		Action:        action,
		DeliveryID:    deliveryID,
		SKU_ID:        skuID,
		Quantity:      quantity,
		UnitOfMeasure: unit,
		PerformedBy:   userID,
		PerformedAt:   time.Now(),
	}
	return s.repo.CreateAuditLog(auditLog)
}
