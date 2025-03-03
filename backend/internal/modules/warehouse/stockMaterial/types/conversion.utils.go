package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	stockMaterialCategoryTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory/types"
)

func ToGenerateBarcodeResponse(barcode string) GenerateBarcodeResponse {
	return GenerateBarcodeResponse{
		Barcode: barcode,
	}
}

func ConvertCreateStockMaterialRequestToStockMaterial(req *CreateStockMaterialDTO) *data.StockMaterial {
	return &data.StockMaterial{
		Name:                   req.Name,
		Description:            req.Description,
		SafetyStock:            req.SafetyStock,
		UnitID:                 req.UnitID,
		Size:                   req.Size,
		CategoryID:             req.CategoryID,
		IngredientID:           req.IngredientID,
		Barcode:                req.Barcode,
		ExpirationPeriodInDays: req.ExpirationPeriodInDays,
		IsActive:               req.IsActive,
	}
}

func ConvertStockMaterialToStockMaterialResponse(stockMaterial *data.StockMaterial) *StockMaterialsDTO {
	isActive := true
	if stockMaterial.IsActive != nil {
		isActive = *stockMaterial.IsActive
	}

	return &StockMaterialsDTO{
		ID:          stockMaterial.ID,
		Name:        stockMaterial.Name,
		Description: stockMaterial.Description,
		SafetyStock: stockMaterial.SafetyStock,
		Unit: unitTypes.UnitsDTO{
			ID:               stockMaterial.UnitID,
			Name:             stockMaterial.Unit.Name,
			ConversionFactor: stockMaterial.Unit.ConversionFactor,
		},
		Size: stockMaterial.Size,
		Category: stockMaterialCategoryTypes.StockMaterialCategoryResponse{
			ID:          stockMaterial.CategoryID,
			Name:        stockMaterial.StockMaterialCategory.Name,
			Description: stockMaterial.StockMaterialCategory.Description,
		},
		Ingredient:             *ingredientTypes.ConvertToIngredientResponseDTO(&stockMaterial.Ingredient),
		Barcode:                stockMaterial.Barcode,
		ExpirationPeriodInDays: stockMaterial.ExpirationPeriodInDays,
		IsActive:               isActive,
		CreatedAt:              stockMaterial.CreatedAt.Format(time.RFC3339),
		UpdatedAt:              stockMaterial.UpdatedAt.Format(time.RFC3339),
	}
}
