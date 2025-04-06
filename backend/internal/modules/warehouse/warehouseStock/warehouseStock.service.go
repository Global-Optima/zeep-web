package warehouseStock

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock/types"
	"go.uber.org/zap"
)

type WarehouseStockService interface {
	ReceiveInventory(warehouseID uint, req types.ReceiveWarehouseDelivery) error
	GetDeliveries(filter types.WarehouseDeliveryFilter) ([]types.WarehouseDeliveryDTO, error)
	GetDeliveryByID(id uint) (*types.WarehouseDeliveryDTO, error)

	AddWarehouseStockMaterial(req types.AdjustWarehouseStock) error
	AddWarehouseStocks(warehouseID uint, req []types.AddWarehouseStockMaterial) error
	DeductFromStock(req types.AdjustWarehouseStock) error
	GetStock(query *types.GetWarehouseStockFilterQuery) ([]types.WarehouseStockResponse, error)
	GetStockMaterialDetails(stockMaterialID uint, filter *contexts.WarehouseContextFilter) (*types.WarehouseStockResponse, error)
	UpdateStock(warehouseID, stockMaterialID uint, dto types.UpdateWarehouseStockDTO) error

	CheckStockNotifications(warehouseID uint, stock data.WarehouseStock) error

	GetAvailableToAddStockMaterials(storeID uint, query *types.AvailableStockMaterialFilter) ([]stockMaterialTypes.StockMaterialsDTO, error)
}

type warehouseStockService struct {
	repo                WarehouseStockRepository
	stockMaterialRepo   stockMaterial.StockMaterialRepository
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewWarehouseStockService(repo WarehouseStockRepository,
	stockMaterialRepo stockMaterial.StockMaterialRepository,
	notificationService notifications.NotificationService,
	logger *zap.SugaredLogger,
) WarehouseStockService {
	return &warehouseStockService{
		repo:                repo,
		stockMaterialRepo:   stockMaterialRepo,
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
		s.logger.Errorf("failed to fetch stock materials: %v", err)
		return types.ErrFetchStockMaterials
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
			s.logger.Errorf("stock material with ID %d not found", material.StockMaterialID)
			return types.ErrStockMaterialNotFound
		}
		materials[i] = data.SupplierWarehouseDeliveryMaterial{
			StockMaterialID: material.StockMaterialID,
			Barcode:         stockMaterial.Barcode,
			Quantity:        material.Quantity,
			Price:           material.Price,
			ExpirationDate:  time.Now().AddDate(0, 0, stockMaterial.ExpirationPeriodInDays),
		}
	}

	err = s.repo.RecordDeliveriesAndUpdateStock(delivery, materials, warehouseID)
	if err != nil {
		s.logger.Errorf("failed to receive inventory: %v", err)
		return types.ErrFailedToRecordDeliveries
	}

	return nil
}

func (s *warehouseStockService) GetDeliveries(filter types.WarehouseDeliveryFilter) ([]types.WarehouseDeliveryDTO, error) {
	deliveries, err := s.repo.GetDeliveries(filter)
	if err != nil {
		s.logger.Errorf("failed to fetch deliveries: %v", err)
		return nil, types.ErrFetchDeliveries
	}

	return types.DeliveriesToDeliveryResponses(deliveries), nil
}

func (s *warehouseStockService) GetDeliveryByID(id uint) (*types.WarehouseDeliveryDTO, error) {
	var delivery data.SupplierWarehouseDelivery
	if err := s.repo.GetDeliveryByID(id, &delivery); err != nil {
		s.logger.Errorf("failed to fetch delivery by ID %d: %v", id, err)
		return nil, types.ErrFetchDelivery
	}

	response := types.ToDeliveryResponse(delivery)
	return &response, nil
}

func (s *warehouseStockService) AddWarehouseStockMaterial(req types.AdjustWarehouseStock) error {
	if err := s.repo.AddToWarehouseStock(req.WarehouseID, req.StockMaterialID, req.Quantity); err != nil {
		s.logger.Errorf("failed to add warehouse stock material: %v", err)
		return types.ErrAddWarehouseStockMaterial
	}
	return nil
}

func (s *warehouseStockService) DeductFromStock(req types.AdjustWarehouseStock) error {
	stock, err := s.repo.DeductFromWarehouseStock(req.WarehouseID, req.StockMaterialID, req.Quantity)
	if err != nil {
		s.logger.Errorf("failed to deduct from stock: %v", err)
		return types.ErrDeductFromStock
	}

	if err = s.checkStockAndNotify(stock); err != nil {
		s.logger.Errorf("failed to check stock and notify: %v", err)
	}
	return nil
}

func (s *warehouseStockService) GetStock(query *types.GetWarehouseStockFilterQuery) ([]types.WarehouseStockResponse, error) {
	stocks, err := s.repo.GetWarehouseStock(query)
	if err != nil {
		s.logger.Errorf("failed to fetch warehouse stocks: %v", err)
		return nil, types.ErrFetchStock
	}

	responses := make([]types.WarehouseStockResponse, len(stocks))
	for i, stock := range stocks {
		responses[i] = types.ToWarehouseStockResponse(stock)
	}
	return responses, nil
}

func (s *warehouseStockService) AddWarehouseStocks(warehouseID uint, req []types.AddWarehouseStockMaterial) error {
	if len(req) == 0 {
		return types.ErrEmptyStocks
	}

	var stocks []data.WarehouseStock
	for _, dto := range req {
		stocks = append(stocks, data.WarehouseStock{
			WarehouseID:     warehouseID,
			StockMaterialID: dto.StockMaterialID,
			Quantity:        dto.Quantity,
		})
	}
	if err := s.repo.AddWarehouseStocks(warehouseID, stocks); err != nil {
		s.logger.Errorf("failed to add warehouse stocks: %v", err)
		return types.ErrAddWarehouseStockMaterial
	}
	return nil
}

func (s *warehouseStockService) GetStockMaterialDetails(stockMaterialID uint, filter *contexts.WarehouseContextFilter) (*types.WarehouseStockResponse, error) {
	aggregatedStock, err := s.repo.GetWarehouseStockMaterialDetails(stockMaterialID, filter)
	if err != nil {
		s.logger.Errorf("failed to fetch stock material details: %v", err)
		return nil, types.ErrFetchStockMaterialDetails
	}
	details := types.ToWarehouseStockResponse(*aggregatedStock)
	return &details, nil
}

func (s *warehouseStockService) UpdateStock(warehouseID, stockMaterialID uint, dto types.UpdateWarehouseStockDTO) error {
	if dto.Quantity == nil && dto.ExpirationDate == nil {
		return types.ErrNothingToUpdate
	}

	if dto.Quantity != nil {
		stock, err := s.repo.UpdateStockQuantity(stockMaterialID, warehouseID, *dto.Quantity)
		if err != nil {
			s.logger.Errorf("failed to update stock quantity: %v", err)
			return types.ErrUpdateStockQuantity
		}
		if err = s.checkStockAndNotify(stock); err != nil {
			s.logger.Errorf("failed to check stock and notify: %v", err)
		}
	}

	if dto.ExpirationDate != nil {
		if err := s.repo.UpdateExpirationDate(stockMaterialID, warehouseID, *dto.ExpirationDate); err != nil {
			s.logger.Errorf("failed to update expiration date: %v", err)
			return types.ErrUpdateExpiration
		}
	}
	return nil
}

func (s *warehouseStockService) CheckStockNotifications(warehouseID uint, stock data.WarehouseStock) error {
	// Check for low stock notification
	if stock.Quantity < stock.StockMaterial.SafetyStock {
		details := &details.OutOfStockDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID:           warehouseID,
				FacilityName: stock.Warehouse.Name,
			},
			ItemName: stock.StockMaterial.Name,
		}
		if err := s.notificationService.NotifyOutOfStock(details); err != nil {
			return fmt.Errorf("failed to send warehouse runout notification: %v", err)
		}
	}

	closestExpirationDate, err := s.repo.FindEarliestExpirationDateForStock(stock.StockMaterialID, &contexts.WarehouseContextFilter{
		WarehouseID: &warehouseID,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch earliest expiration date for stock: %v", err)
	}
	if closestExpirationDate == nil {
		defaultClosestExpirationDate := stock.UpdatedAt.Add(time.Duration(stock.StockMaterial.ExpirationPeriodInDays) * 24 * time.Hour)
		closestExpirationDate = &defaultClosestExpirationDate
	}
	if closestExpirationDate.Before(time.Now().Add(7 * 24 * time.Hour)) { // Expiration within 7 days
		expDetails := &details.WarehouseStockExpirationDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID:           warehouseID,
				FacilityName: stock.Warehouse.Name,
			},
			ItemName:       stock.StockMaterial.Name,
			ExpirationDate: closestExpirationDate.Format("2006-01-02"),
		}
		if err := s.notificationService.NotifyWarehouseStockExpiration(expDetails); err != nil {
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
		if err := s.notificationService.NotifyOutOfStock(details); err != nil {
			return fmt.Errorf("failed to send out of stock notification: %w", err)
		}
	}
	return nil
}

func (s *warehouseStockService) GetAvailableToAddStockMaterials(storeID uint, query *types.AvailableStockMaterialFilter) ([]stockMaterialTypes.StockMaterialsDTO, error) {
	stocks, err := s.repo.GetAvailableToAddStockMaterials(storeID, query)
	if err != nil {
		s.logger.Errorf("failed to fetch available stock materials: %v", err)
		return []stockMaterialTypes.StockMaterialsDTO{}, types.ErrFetchStockMaterials
	}

	stockMaterialResponses := make([]stockMaterialTypes.StockMaterialsDTO, 0)
	for _, stockMaterial := range stocks {
		stockMaterialResponses = append(stockMaterialResponses, *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&stockMaterial))
	}
	return stockMaterialResponses, nil
}
