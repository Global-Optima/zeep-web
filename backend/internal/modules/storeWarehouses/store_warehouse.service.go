package storeWarehouses

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type StoreWarehouseService interface {
	AddStock(storeId uint, dto *types.AddStockDTO) (uint, error)
	AddMultipleStock(storeId uint, dto *types.AddMultipleStockDTO) error
	GetStockList(storeId uint, query *types.GetStockQuery) ([]types.StockDTO, error)
	GetStockById(storeId, stockId uint) (*types.StockDTO, error)
	UpdateStockById(storeId, stockId uint, input *types.UpdateStockDTO) error
	DeleteStockById(storeId, stockId uint) error
}

type storeWarehouseService struct {
	repo   StoreWarehouseRepository
	logger *zap.SugaredLogger
}

func NewStoreWarehouseService(repo StoreWarehouseRepository, logger *zap.SugaredLogger) StoreWarehouseService {
	return &storeWarehouseService{
		repo:   repo,
		logger: logger,
	}
}

func (s *storeWarehouseService) AddMultipleStock(storeId uint, dto *types.AddMultipleStockDTO) error {
	// Start a transaction
	err := s.repo.WithTransaction(func(txRepo storeWarehouseRepository) error {
		for _, stock := range dto.IngredientStocks {
			// Add or update stock for each ingredient
			err := txRepo.AddOrUpdateStock(storeId, &stock)
			if err != nil {
				wrappedErr := utils.WrapError("error adding/updating stock element", err)
				s.logger.Error(wrappedErr)
				return wrappedErr
			}
		}
		return nil
	})

	if err != nil {
		wrappedErr := utils.WrapError("error adding multiple stock elements", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *storeWarehouseService) AddStock(storeId uint, dto *types.AddStockDTO) (uint, error) {
	id, err := s.repo.AddStock(storeId, dto)
	if err != nil {
		wrappedErr := utils.WrapError("error adding new stock element", err)
		s.logger.Error(wrappedErr)
		return 0, err
	}

	return id, nil
}

func (s *storeWarehouseService) GetStockList(storeId uint, query *types.GetStockQuery) ([]types.StockDTO, error) {
	stockList, err := s.repo.GetStockList(storeId, query)
	if err != nil {
		wrappedErr := utils.WrapError("error getting store stock list", err)
		s.logger.Error(wrappedErr)
		return nil, err
	}

	dtos := make([]types.StockDTO, len(stockList))
	for i, stock := range stockList {
		dtos[i] = types.MapToStockDTO(stock)
	}

	return dtos, nil
}

func (s *storeWarehouseService) GetStockById(storeId, stockId uint) (*types.StockDTO, error) {
	stock, err := s.repo.GetStockById(storeId, stockId)
	if err != nil {
		wrappedErr := utils.WrapError("error getting stock", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	stockDto := types.MapToStockDTO(*stock)

	return &stockDto, nil
}

func (s *storeWarehouseService) UpdateStockById(storeId, stockId uint, input *types.UpdateStockDTO) error {
	err := s.repo.UpdateStock(storeId, stockId, input)
	if err != nil {
		wrappedErr := utils.WrapError("error updating stock", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *storeWarehouseService) DeleteStockById(storeId, stockId uint) error {
	err := s.repo.DeleteStockById(storeId, stockId)
	if err != nil {
		wrappedErr := utils.WrapError("error deleting stock", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}
