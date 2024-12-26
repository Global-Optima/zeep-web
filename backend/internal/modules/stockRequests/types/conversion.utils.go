package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

func ToStockRequestResponse(request *data.StockRequest) StockRequestResponse {
	items := make([]StockRequestStockMaterialResponse, len(request.Ingredients))
	for i, ingredient := range request.Ingredients {
		var packageMeasures utils.PackageMeasure

		if ingredient.StockMaterial.Package != nil {
			packageMeasures = utils.ReturnPackageMeasure(ingredient, ingredient.Quantity)
		}

		items[i] = StockRequestStockMaterialResponse{
			StockMaterialID: ingredient.StockMaterialID,
			Name:            ingredient.StockMaterial.Name,
			Category:        ingredient.StockMaterial.StockMaterialCategory.Name,
			PackageMeasure:  packageMeasures,
		}
	}

	return StockRequestResponse{
		RequestID: request.ID,
		Store: StoreDTO{
			ID:      request.StoreID,
			Name:    request.Store.Name,
			Address: request.Store.FacilityAddress.Address,
		},
		Warehouse: WarehouseDTO{
			ID:   request.WarehouseID,
			Name: request.Warehouse.Name,
		},
		Status:         request.Status,
		StockMaterials: items,
		CreatedAt:      request.CreatedAt,
		UpdatedAt:      request.UpdatedAt,
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
