package inventory

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/inventory/types"
)

type InventoryService interface {
	TransferInventory(req types.TransferInventoryRequest) error
	GetInventoryLevels(warehouseID uint) ([]types.InventoryItem, error)
	PickupStock(req types.PickupRequest) error
	GetExpiringItems(warehouseID uint, thresholdDays int) ([]types.UpcomingExpirationResponse, error)
	ExtendExpiration(req types.ExtendExpirationRequest) error
}

type inventoryService struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) ReceiveInventory(req types.ReceiveInventoryRequest) error {
	deliveries := make([]data.Delivery, len(req.Items))
	for i, item := range req.Items {
		deliveries[i] = data.Delivery{
			SKU_ID:         item.SKU_ID,
			Status:         "Pending",
			Source:         req.SupplierID,
			Target:         req.WarehouseID,
			Barcode:        generateBarcodeForDelivery(item),
			Quantity:       item.Quantity,
			DeliveryDate:   time.Now(),
			ExpirationDate: time.Now().AddDate(0, 0, 1095),
		}
	}

	if err := s.repo.LogIncomingInventory(deliveries); err != nil {
		return fmt.Errorf("failed to log incoming inventory: %w", err)
	}

	stockItems, err := types.ConvertInventoryItemsToStockRequest(req.Items, s.repo)
	if err != nil {
		return fmt.Errorf("failed to convert inventory items: %w", err)
	}

	if err := s.repo.UpdateStockLevels(req.WarehouseID, stockItems); err != nil {
		return fmt.Errorf("failed to update stock levels: %w", err)
	}

	return nil
}

func (s *inventoryService) TransferInventory(req types.TransferInventoryRequest) error {
	stockItems, err := types.ConvertInventoryItemsToStockRequest(req.Items, s.repo)
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
	stockItems, err := types.ConvertInventoryItemsToStockRequest(req.Items, s.repo)
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
