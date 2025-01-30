package warehouseStock

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock/types"
	"go.uber.org/zap"
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

	CheckStockNotifications(warehouseID uint, stock data.WarehouseStock) error

	GetAvailableToAddStockMaterials(storeID uint, query *types.AvailableStockMaterialFilter) ([]stockMaterialTypes.StockMaterialsDTO, error)
}

type warehouseStockService struct {
	repo                WarehouseStockRepository
	stockMaterialRepo   stockMaterial.StockMaterialRepository
	barcodeRepo         barcode.BarcodeRepository
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewWarehouseStockService(repo WarehouseStockRepository,
	stockMaterialRepo stockMaterial.StockMaterialRepository,
	barcodeRepo barcode.BarcodeRepository,
	notificationService notifications.NotificationService,
	logger *zap.SugaredLogger) WarehouseStockService {
	return &warehouseStockService{
		repo:                repo,
		stockMaterialRepo:   stockMaterialRepo,
		barcodeRepo:         barcodeRepo,
		notificationService: notificationService,
		logger:              logger,
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

		materials[i] = data.SupplierWarehouseDeliveryMaterial{
			StockMaterialID: material.StockMaterialID,
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
	stock, err := s.repo.DeductFromWarehouseStock(req.WarehouseID, req.StockMaterialID, req.Quantity)
	if err != nil {
		return fmt.Errorf("failed to deduct from stock: %w", err)
	}

	err = s.checkStockAndNotify(stock)
	if err != nil {
		s.logger.Errorf("failed to check stock and notify: %w", err)
	}

	return nil
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
	if dto.Quantity == nil && dto.ExpirationDate == nil {
		return fmt.Errorf("nothing to update")
	}

	if dto.Quantity != nil {
		stock, err := s.repo.UpdateStockQuantity(stockMaterialID, warehouseID, *dto.Quantity)
		if err != nil {
			return fmt.Errorf("failed to update stock quantity: %w", err)
		}

		err = s.checkStockAndNotify(stock)
		if err != nil {
			s.logger.Errorf("failed to check stock and notify: %w", err)
		}
	}

	if dto.ExpirationDate != nil {
		err := s.repo.UpdateExpirationDate(stockMaterialID, warehouseID, *dto.ExpirationDate)
		if err != nil {
			return fmt.Errorf("failed to update expiration date: %w", err)
		}
	}

	return nil
}

func (s *warehouseStockService) CheckStockNotifications(warehouseID uint, stock data.WarehouseStock) error {
	if stock.Quantity < stock.StockMaterial.SafetyStock {
		details := &details.StoreWarehouseRunOutDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID:           warehouseID,
				FacilityName: stock.Warehouse.Name,
			},
			StockItem:   stock.StockMaterial.Name,
			StockItemID: stock.ID,
		}
		err := s.notificationService.NotifyStoreWarehouseRunOut(details)
		if err != nil {
			return fmt.Errorf("failed to send warehouse runout notification: %v", err)
		}
	}

	closestExpirationDate := stock.UpdatedAt.Add(time.Duration(stock.StockMaterial.ExpirationPeriodInDays) * 24 * time.Hour)

	if closestExpirationDate.Before(time.Now().Add(7 * 24 * time.Hour)) { // Expiration within 7 days
		details := &details.StockExpirationDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID:           warehouseID,
				FacilityName: stock.Warehouse.Name,
			},
			ItemName:       stock.StockMaterial.Name,
			ExpirationDate: closestExpirationDate.Format("2006-01-02"),
		}
		err := s.notificationService.NotifyStockExpiration(details)
		if err != nil {
			return fmt.Errorf("failed to send stock expiration notification: %v", err)
		}
	}

	return nil
}

func (s *warehouseStockService) checkStockAndNotify(stock *data.WarehouseStock) error {
	if stock.Quantity < stock.StockMaterial.SafetyStock {
		details := &details.OutOfStockDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID:           stock.WarehouseID,
				FacilityName: stock.Warehouse.Name,
			},
			ItemName: stock.StockMaterial.Name,
		}

		err := s.notificationService.NotifyOutOfStock(details)
		if err != nil {
			return fmt.Errorf("failed to send out of stock notification: %w", err)
		}
	}

	return nil
}

func (s *warehouseStockService) GetAvailableToAddStockMaterials(
	storeID uint,
	query *types.AvailableStockMaterialFilter,
) ([]stockMaterialTypes.StockMaterialsDTO, error) {

	stocks, err := s.repo.GetAvailableToAddStockMaterials(storeID, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch available stock materials: %w", err)
	}

	stockMaterialResponses := make([]stockMaterialTypes.StockMaterialsDTO, 0)
	for _, stockMaterial := range stocks {
		stockMaterialResponses = append(stockMaterialResponses, *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&stockMaterial))
	}

	return stockMaterialResponses, nil
}
