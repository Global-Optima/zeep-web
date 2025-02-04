package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
)

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
