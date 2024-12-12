package storeWarehouses

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type StoreWarehouseService interface {
	AddStock(storeId uint, dto types.AddStockDTO) (uint, error)
	GetStockList(storeId uint, query types.GetStockQuery) ([]types.StockDTO, error)
	GetStockById(storeId, stockId uint) (*types.StockDTO, error)
	UpdateStockById(storeId, stockId uint, input types.UpdateStockDTO) error
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

func (s *storeWarehouseService) AddStock(storeId uint, dto types.AddStockDTO) (uint, error) {
	id, err := s.repo.AddStock(storeId, dto)
	if err != nil {
		wrappedErr := utils.WrapError("error adding new stock element", err)
		s.logger.Error(wrappedErr)
		return 0, err
	}

	return id, nil
}

func (s *storeWarehouseService) GetStockList(storeId uint, query types.GetStockQuery) ([]types.StockDTO, error) {
	ingredients, err := s.repo.GetStockList(storeId, query)
	if err != nil {
		wrappedErr := utils.WrapError("error getting store stock list", err)
		s.logger.Error(wrappedErr)
		return nil, err
	}

	dtos := make([]types.StockDTO, len(ingredients))
	for i, ingredient := range ingredients {
		dtos[i] = *mapToStoreWarehouseIngredient(&ingredient)
	}

	return dtos, nil
}

func (s *storeWarehouseService) GetStockById(storeId, stockId uint) (*types.StockDTO, error) {
	ingredient, err := s.repo.GetStockById(storeId, stockId)
	if err != nil {
		wrappedErr := utils.WrapError("error getting store ingredient", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return mapToStoreWarehouseIngredient(ingredient), nil
}

func (s *storeWarehouseService) UpdateStockById(storeId, stockId uint, input types.UpdateStockDTO) error {
	updateFields, err := types.PrepareUpdateFields(input)
	if err != nil {
		wrappedErr := utils.WrapError("error preparing update fields", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return s.repo.PartialUpdateStock(storeId, stockId, updateFields)
}

func mapToStoreWarehouseIngredient(storeIngredient *data.StoreWarehouseStock) *types.StockDTO {
	dto := &types.StockDTO{
		ID:                storeIngredient.IngredientID,
		Name:              storeIngredient.Ingredient.Name,
		CurrentStock:      storeIngredient.Quantity,
		LowStockThreshold: storeIngredient.LowStockThreshold,
		LowStockAlert:     storeIngredient.Quantity < storeIngredient.LowStockThreshold,
	}

	return dto
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
