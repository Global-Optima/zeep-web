package storeStocks

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"

	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type StoreStockService interface {
	GetAvailableIngredientsToAdd(storeID uint, filter *ingredientTypes.IngredientFilter) ([]ingredientTypes.IngredientDTO, error)
	AddStock(storeId uint, dto *types.AddStoreStockDTO) (uint, error)
	AddMultipleStock(storeId uint, dto *types.AddMultipleStoreStockDTO) ([]uint, error)
	GetStockList(storeId uint, query *types.GetStockFilterQuery) ([]types.StoreStockDTO, error)
	GetStockListByIDs(storeId uint, stockIds []uint) ([]types.StoreStockDTO, error)
	GetStockById(stockId uint, filter *contexts.StoreContextFilter) (*types.StoreStockDTO, error)
	UpdateStockById(storeId, stockId uint, input *types.UpdateStoreStockDTO) error
	DeleteStockById(storeId, stockId uint) error

	CheckStockNotifications(storeID uint, stock data.StoreStock) error
}

type storeStockService struct {
	repo                StoreStockRepository
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewStoreStockService(repo StoreStockRepository, notificationService notifications.NotificationService, logger *zap.SugaredLogger) StoreStockService {
	return &storeStockService{
		repo:                repo,
		notificationService: notificationService,
		logger:              logger,
	}
}

func (s *storeStockService) AddMultipleStock(storeId uint, dto *types.AddMultipleStoreStockDTO) ([]uint, error) {
	var IDs []uint

	// Start a transaction
	err := s.repo.WithTransaction(func(txRepo storeStockRepository) error {
		for _, stock := range dto.IngredientStocks {
			// Add or update stock for each ingredient
			id, err := txRepo.AddOrSaveStock(storeId, &stock)
			if err != nil {
				wrappedErr := utils.WrapError("error adding/updating stock element", err)
				s.logger.Error(wrappedErr)
				return wrappedErr
			}

			err = s.checkStockAndNotify(storeId, id)
			if err != nil {
				s.logger.Errorf("failed to check stock and notify: %v", err)
			}

			IDs = append(IDs, id)
		}
		return nil
	})
	if err != nil {
		wrappedErr := utils.WrapError("error adding multiple stock elements", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return IDs, nil
}

func (s *storeStockService) AddStock(storeId uint, dto *types.AddStoreStockDTO) (uint, error) {
	id, err := s.repo.AddStock(storeId, dto)
	if err != nil {
		wrappedErr := utils.WrapError("error adding new stock element", err)
		s.logger.Error(wrappedErr)
		return 0, err
	}

	err = s.checkStockAndNotify(storeId, id)
	if err != nil {
		s.logger.Errorf("failed to check stock and notify: %v", err)
	}

	return id, nil
}

func (s *storeStockService) GetStockList(storeId uint, query *types.GetStockFilterQuery) ([]types.StoreStockDTO, error) {
	stockList, err := s.repo.GetStockList(storeId, query)
	if err != nil {
		wrappedErr := utils.WrapError("error getting store stock list", err)
		s.logger.Error(wrappedErr)
		return nil, err
	}

	dtos := make([]types.StoreStockDTO, len(stockList))
	for i, stock := range stockList {
		dtos[i] = types.MapToStockDTO(stock)
	}

	return dtos, nil
}

func (s *storeStockService) GetStockListByIDs(storeId uint, IDs []uint) ([]types.StoreStockDTO, error) {
	stockList, err := s.repo.GetStockListByIDs(storeId, IDs)
	if err != nil {
		wrappedErr := utils.WrapError("error getting store stock list", err)
		s.logger.Error(wrappedErr)
		return nil, err
	}

	dtos := make([]types.StoreStockDTO, len(stockList))
	for i, stock := range stockList {
		dtos[i] = types.MapToStockDTO(stock)
	}

	return dtos, nil
}

func (s *storeStockService) GetStockById(stockId uint, filter *contexts.StoreContextFilter) (*types.StoreStockDTO, error) {
	stock, err := s.repo.GetStockById(stockId, filter)
	if err != nil {
		wrappedErr := utils.WrapError("error getting stock", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	stockDto := types.MapToStockDTO(*stock)

	return &stockDto, nil
}

func (s *storeStockService) UpdateStockById(storeId, stockId uint, input *types.UpdateStoreStockDTO) error {
	storeStock, err := s.repo.GetRawStockByID(storeId, stockId)
	if err != nil {
		return err
	}

	err = types.UpdateToStoreStockModel(input, storeStock)
	if err != nil {
		wrappedErr := utils.WrapError("error updating store stock", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	err = s.repo.SaveStock(storeId, stockId, storeStock)
	if err != nil {
		wrappedErr := utils.WrapError("error updating stock", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	err = s.checkStockAndNotify(storeId, stockId)
	if err != nil {
		s.logger.Errorf("failed to check stock and notify: %v", err)
	}

	return nil
}

func (s *storeStockService) DeleteStockById(storeId, stockId uint) error {
	storeStock, err := s.repo.GetRawStockByID(storeId, stockId)
	if err != nil {
		wrappedErr := utils.WrapError("error getting store stock", err)
		s.logger.Error(wrappedErr)
		return err
	}
	if storeStock.Quantity > 0 {
		s.logger.Error("error deleting store stock, quantity must be 0")
		return types.ErrStockIsInUse
	}

	isInUse, err := s.repo.CheckStoreStockUsage(stockId)
	if err != nil {
		wrappedErr := utils.WrapError("error checking store stock usage", err)
		s.logger.Error(wrappedErr)
		return err
	}
	if isInUse {

		s.logger.Error("error deleting store stock, used in store")
		return types.ErrStockIsInUse
	}

	err = s.repo.DeleteStockById(storeId, stockId)
	if err != nil {
		wrappedErr := utils.WrapError("error deleting stock", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *storeStockService) GetAvailableIngredientsToAdd(storeID uint, filter *ingredientTypes.IngredientFilter) ([]ingredientTypes.IngredientDTO, error) {
	ingredients, err := s.repo.GetAvailableIngredientsToAdd(storeID, filter)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch available ingredients to add for store: %d: %w", storeID, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return ingredientTypes.ConvertToIngredientResponseDTOs(ingredients), nil
}

func (s *storeStockService) CheckStockNotifications(storeID uint, stock data.StoreStock) error {
	if stock.Quantity <= stock.LowStockThreshold {
		details := &details.StoreWarehouseRunOutDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID:           storeID,
				FacilityName: stock.Store.Name,
			},
			StockItem:   stock.Ingredient.Name,
			StockItemID: stock.ID,
		}
		err := s.notificationService.NotifyStoreWarehouseRunOut(details)
		if err != nil {
			s.logger.Errorf("failed to send store warehouse runout notification: %v", err)
		}
	}

	earliestExpirationDate, err := s.repo.FindEarliestExpirationForIngredient(stock.Ingredient.ID, storeID)
	if err != nil {
		s.logger.Errorf("failed to fetch earliest expiration date for ingredient %d: %v", stock.Ingredient.ID, err)
	}

	var closestExpirationDate time.Time
	if earliestExpirationDate != nil {
		closestExpirationDate = *earliestExpirationDate
	} else {
		closestExpirationDate = stock.UpdatedAt.Add(time.Duration(stock.Ingredient.ExpirationInDays) * 24 * time.Hour)
	}

	if closestExpirationDate.Before(time.Now().Add(7 * 24 * time.Hour)) { // Expiration within 7 days
		details := &details.StoreStockExpirationDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID:           storeID,
				FacilityName: stock.Store.Name,
			},
			ItemName:       stock.Ingredient.Name,
			ExpirationDate: closestExpirationDate.Format("2006-01-02"),
		}

		err := s.notificationService.NotifyStoreStockExpiration(details)
		if err != nil {
			s.logger.Errorf("failed to send stock expiration notification: %v", err)
		}
	}
	return nil
}

func (s *storeStockService) checkStockAndNotify(storeId, stockId uint) error {
	updatedStock, err := s.repo.GetStockById(stockId, &contexts.StoreContextFilter{StoreID: &storeId})
	if err != nil {
		s.logger.Errorf("failed to fetch stock for %d: %v", stockId, err)
		return err
	}

	if updatedStock == nil {
		s.logger.Errorf("stock with ID %d not found", stockId)
		return fmt.Errorf("stock with ID %d not found", stockId)
	}

	if updatedStock.Quantity <= updatedStock.LowStockThreshold {
		details := &details.StoreWarehouseRunOutDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID:           storeId,
				FacilityName: updatedStock.Store.Name,
			},
			StockItem:   updatedStock.Ingredient.Name,
			StockItemID: updatedStock.IngredientID,
		}

		err := s.notificationService.NotifyStoreWarehouseRunOut(details)
		if err != nil {
			s.logger.Errorf("failed to send store warehouse runout notification: %v", err)
			return err
		}
	}

	return nil
}
