package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	stockMaterialCategoryTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory/types"
	stockMaterialPackageTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialPackage/types"
)

func ConvertCreateStockMaterialRequestToStockMaterial(req *CreateStockMaterialDTO) *data.StockMaterial {
	return &data.StockMaterial{
		Name:                   req.Name,
		Description:            req.Description,
		SafetyStock:            req.SafetyStock,
		UnitID:                 req.UnitID,
		CategoryID:             req.CategoryID,
		IngredientID:           req.IngredientID,
		Barcode:                req.Barcode,
		ExpirationPeriodInDays: req.ExpirationPeriodInDays,
		IsActive:               true,
	}
}

func ConvertPackageDTOToModel(stockMaterialID uint, req *CreateStockMaterialPackagesDTO) *data.StockMaterialPackage {
	return &data.StockMaterialPackage{
		StockMaterialID: stockMaterialID,
		Size:            req.Size,
		UnitID:          req.UnitID,
	}
}

func ConvertStockMaterialToStockMaterialResponse(stockMaterial *data.StockMaterial) *StockMaterialsDTO {
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
		Category: stockMaterialCategoryTypes.StockMaterialCategoryResponse{
			ID:          stockMaterial.CategoryID,
			Name:        stockMaterial.StockMaterialCategory.Name,
			Description: stockMaterial.StockMaterialCategory.Description,
		},
		Ingredient:             *ingredientTypes.ConvertToIngredientResponseDTO(&stockMaterial.Ingredient),
		Barcode:                stockMaterial.Barcode,
		ExpirationPeriodInDays: stockMaterial.ExpirationPeriodInDays,
		IsActive:               stockMaterial.IsActive,
		Packages:               stockMaterialPackageTypes.ToStockMaterialPackageResponses(stockMaterial.Packages),
		CreatedAt:              stockMaterial.CreatedAt.Format(time.RFC3339),
		UpdatedAt:              stockMaterial.UpdatedAt.Format(time.RFC3339),
	}
}

func ConvertUpdateStockMaterialRequestToStockMaterial(stockMaterial *data.StockMaterial, req *UpdateStockMaterialDTO) error {
	return ValidateAndApplyUpdate(stockMaterial, req)
}
