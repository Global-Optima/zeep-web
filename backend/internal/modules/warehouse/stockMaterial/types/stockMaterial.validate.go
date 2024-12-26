package types

import (
	"errors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ValidateAndApplyUpdate(stockMaterial *data.StockMaterial, req *UpdateStockMaterialDTO) error {
	if req.Name != nil {
		if *req.Name == "" {
			return errors.New("stockMaterial name cannot be empty")
		}
		stockMaterial.Name = *req.Name
	}

	if req.Description != nil {
		stockMaterial.Description = *req.Description
	}

	if req.SafetyStock != nil {
		if *req.SafetyStock <= 0 {
			return errors.New("stockMaterial safety stock must be greater than zero")
		}
		stockMaterial.SafetyStock = *req.SafetyStock
	}

	if req.ExpirationFlag != nil {
		stockMaterial.ExpirationFlag = *req.ExpirationFlag
	}

	if req.UnitID != nil {
		stockMaterial.UnitID = *req.UnitID
	}

	if req.CategoryID != nil {
		stockMaterial.CategoryID = *req.CategoryID
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

	stockMaterial.UpdatedAt = time.Now()

	return nil
}
