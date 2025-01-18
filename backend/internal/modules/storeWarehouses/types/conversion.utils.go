package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func MapToStockDTO(stock data.StoreWarehouseStock) StoreStockDTO {
	return StoreStockDTO{
		ID:                stock.ID,
		Name:              stock.Ingredient.Name,
		Quantity:          stock.Quantity,
		Unit:              stock.Ingredient.Unit.Name,
		LowStockThreshold: stock.LowStockThreshold,
		LowStockAlert:     stock.Quantity < stock.LowStockThreshold,
		IngredientID:      stock.IngredientID,
	}
}

func AddToStock(dto AddStoreStockDTO, storeWarehouseID uint) *data.StoreWarehouseStock {
	return &data.StoreWarehouseStock{
		StoreWarehouseID:  storeWarehouseID,
		Quantity:          dto.Quantity,
		LowStockThreshold: dto.LowStockThreshold,
		IngredientID:      dto.IngredientID,
	}
}
