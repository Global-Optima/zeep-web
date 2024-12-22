package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ToStockRequestResponse(request *data.StockRequest) StockRequestResponse {
	items := make([]StockRequestItemResponse, len(request.Ingredients))
	for i, ingredient := range request.Ingredients {
		items[i] = StockRequestItemResponse{
			StockMaterialID: ingredient.IngredientID,
			Name:            ingredient.Ingredient.Name,
			Category:        "todo: add ingredients",
			Unit:            ingredient.Ingredient.Unit.Name,
			Quantity:        ingredient.Quantity,
		}
	}

	return StockRequestResponse{
		RequestID:     request.ID,
		StoreID:       request.StoreID,
		StoreName:     request.Store.Name,
		WarehouseID:   request.WarehouseID,
		WarehouseName: request.Warehouse.Name,
		Status:        request.Status,
		Items:         items,
		CreatedAt:     request.CreatedAt,
		UpdatedAt:     request.UpdatedAt,
	}
}

func ToLowStockIngredientResponse(stock data.StoreWarehouseStock) LowStockIngredientResponse {
	return LowStockIngredientResponse{
		IngredientID:      stock.IngredientID,
		Name:              stock.Ingredient.Name,
		Unit:              stock.Ingredient.Unit.Name,
		Quantity:          stock.Quantity,
		LowStockThreshold: stock.LowStockThreshold,
	}
}
