package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ValidateAndApplyUpdate(stockMaterial *data.StockMaterial, req *UpdateStockMaterialDTO) (*data.StockMaterial, error) {
	if stockMaterial == nil {
		return nil, errors.New("stockMaterial cannot be nil")
	}

	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("stockMaterial name cannot be empty")
		}
		stockMaterial.Name = *req.Name
	}

	if req.Description != nil {
		stockMaterial.Description = *req.Description
	}

	if req.SafetyStock != nil {
		if *req.SafetyStock <= 0 {
			return nil, errors.New("stockMaterial safety stock must be greater than zero")
		}
		stockMaterial.SafetyStock = *req.SafetyStock
	}

	if req.UnitID != nil {
		stockMaterial.UnitID = *req.UnitID
	}

	if req.Size != nil {
		stockMaterial.Size = *req.Size
	}

	if req.CategoryID != nil {
		stockMaterial.CategoryID = *req.CategoryID
	}

	if req.IngredientID != nil {
		stockMaterial.IngredientID = *req.IngredientID
	}

	if req.Barcode != nil {
		stockMaterial.Barcode = *req.Barcode
	}

	if req.ExpirationPeriodInDays != nil {
		stockMaterial.ExpirationPeriodInDays = *req.ExpirationPeriodInDays
	}

	if req.IsActive != nil {
		stockMaterial.IsActive = *req.IsActive
	}

	return stockMaterial, nil
}
