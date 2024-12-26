package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func MapToStockDTO(stock data.StoreWarehouseStock) StockDTO {
	return StockDTO{
		ID:                stock.ID,
		Name:              stock.Ingredient.Name,
		Quantity:          stock.Quantity,
		Unit:              stock.Ingredient.Unit.Name,
		LowStockThreshold: stock.LowStockThreshold,
		LowStockAlert:     stock.Quantity < stock.LowStockThreshold,
	}
}
