package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ConvertCreateStockMaterialRequestToStockMaterial(req *CreateStockMaterialDTO) *data.StockMaterial {
	return &data.StockMaterial{
		Name:                   req.Name,
		Description:            req.Description,
		SafetyStock:            req.SafetyStock,
		ExpirationFlag:         req.ExpirationFlag,
		UnitID:                 req.UnitID,
		CategoryID:             req.CategoryID,
		IngredientID:           req.IngredientID,
		Barcode:                req.Barcode,
		ExpirationPeriodInDays: req.ExpirationPeriodInDays,
		IsActive:               true,
	}
}

func ConvertStockMaterialToStockMaterialResponse(stockMaterial *data.StockMaterial) *StockMaterialsDTO {
	return &StockMaterialsDTO{
		ID:                     stockMaterial.ID,
		Name:                   stockMaterial.Name,
		Description:            stockMaterial.Description,
		SafetyStock:            stockMaterial.SafetyStock,
		ExpirationFlag:         stockMaterial.ExpirationFlag,
		UnitID:                 stockMaterial.UnitID,
		UnitName:               stockMaterial.Unit.Name,
		Category:               stockMaterial.StockMaterialCategory.Name,
		Ingredient:             stockMaterial.Ingredient.Name,
		Barcode:                stockMaterial.Barcode,
		ExpirationPeriodInDays: stockMaterial.ExpirationPeriodInDays,
		IsActive:               stockMaterial.IsActive,
		CreatedAt:              stockMaterial.CreatedAt.Format(time.RFC3339),
		UpdatedAt:              stockMaterial.UpdatedAt.Format(time.RFC3339),
	}
}

func ConvertUpdateStockMaterialRequestToStockMaterial(stockMaterial *data.StockMaterial, req *UpdateStockMaterialDTO) error {
	return ValidateAndApplyUpdate(stockMaterial, req)
}
