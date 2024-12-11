package storeWarehouses

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type StoreWarehouseService interface {
	GetStoreWarehouseStockList(query types.GetStoreWarehouseStockQuery) ([]types.StoreWarehouseIngredientDTO, error)
	GetStoreWarehouseStockById(storeId, ingredientId uint) (*types.StoreWarehouseIngredientDTO, error)
	UpdateStoreWarehouseStockById(id uint, input types.UpdateStoreWarehouseIngredientDTO) error
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

func (s *storeWarehouseService) GetStoreWarehouseStockList(query types.GetStoreWarehouseStockQuery) ([]types.StoreWarehouseIngredientDTO, error) {
	ingredients, err := s.repo.GetStoreWarehouseStockList(query)
	if err != nil {
		s.logger.Errorf("error getting store ingredients: %v", err)
		return nil, err
	}

	dtos := make([]types.StoreWarehouseIngredientDTO, len(ingredients))
	for i, ingredient := range ingredients {
		dtos[i] = *mapToStoreWarehouseIngredient(&ingredient)
	}

	return dtos, nil
}

func (s *storeWarehouseService) GetStoreWarehouseStockById(storeId, ingredientId uint) (*types.StoreWarehouseIngredientDTO, error) {
	ingredient, err := s.repo.GetStoreWarehouseStockById(storeId, ingredientId)
	if err != nil {
		errMsg := "error getting store ingredient"
		s.logger.Error(errors.Wrapf(err, errMsg))
		return nil, errors.Wrapf(err, errMsg)
	}

	return mapToStoreWarehouseIngredient(ingredient), nil
}

func (s *storeWarehouseService) UpdateStoreWarehouseStockById(id uint, input types.UpdateStoreWarehouseIngredientDTO) error {
	updateFields, err := types.PrepareUpdateFields(input)
	if err != nil {
		errMsg := "error preparing update fields"
		s.logger.Error(errors.Wrapf(err, errMsg))
		return errors.Wrapf(err, errMsg)
	}

	return s.repo.PartialUpdateStoreWarehouseStock(id, updateFields)
}

func mapToStoreWarehouseIngredient(storeIngredient *data.StoreWarehouseStock) *types.StoreWarehouseIngredientDTO {
	dto := &types.StoreWarehouseIngredientDTO{
		ID:                storeIngredient.IngredientID,
		Name:              storeIngredient.Ingredient.Name,
		CurrentStock:      storeIngredient.Quantity,
		LowStockThreshold: storeIngredient.LowStockThreshold,
		LowStockAlert:     storeIngredient.Quantity < storeIngredient.LowStockThreshold,
	}

	return dto
}
