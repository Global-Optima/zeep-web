package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
)

const DEFAULT_LOW_STOCK_THRESHOLD = 2 // for conversion factor equal to kg or l

func MapToStockDTO(stock data.StoreStock) StoreStockDTO {
	return StoreStockDTO{
		ID:                stock.ID,
		Name:              stock.Ingredient.Name,
		Quantity:          stock.Quantity,
		LowStockThreshold: stock.LowStockThreshold,
		LowStockAlert:     stock.Quantity < stock.LowStockThreshold,
		Ingredient:        *ingredientTypes.ConvertToIngredientResponseDTO(&stock.Ingredient),
	}
}

func AddToStock(dto AddStoreStockDTO, storeID uint) *data.StoreStock {
	return &data.StoreStock{
		StoreID:           storeID,
		Quantity:          dto.Quantity,
		LowStockThreshold: dto.LowStockThreshold,
		IngredientID:      dto.IngredientID,
	}
}

func DefaultStockFromIngredient(storeID uint, ingredient *data.Ingredient) (*data.StoreStock, error) {
	if ingredient == nil || ingredient.ID == 0 || ingredient.Unit.ConversionFactor == 0 {
		return nil, fmt.Errorf("not enough ingredient info to create default stock: no ingredient or unit preload fetched")
	}

	return &data.StoreStock{
		StoreID:           storeID,
		IngredientID:      ingredient.ID,
		Quantity:          0,
		LowStockThreshold: DEFAULT_LOW_STOCK_THRESHOLD / ingredient.Unit.ConversionFactor,
	}, nil
}
