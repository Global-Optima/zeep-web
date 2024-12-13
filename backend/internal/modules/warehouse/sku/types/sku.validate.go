package types

import (
	"errors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ValidateAndApplyUpdate(sku *data.StockMaterial, req *UpdateSKURequest) error {
	if req.Name != nil {
		if *req.Name == "" {
			return errors.New("SKU name cannot be empty")
		}
		sku.Name = *req.Name
	}

	if req.Description != nil {
		sku.Description = *req.Description
	}

	if req.SafetyStock != nil {
		if *req.SafetyStock <= 0 {
			return errors.New("SKU safety stock must be greater than zero")
		}
		sku.SafetyStock = *req.SafetyStock
	}

	if req.ExpirationFlag != nil {
		sku.ExpirationFlag = *req.ExpirationFlag
	}

	if req.UnitID != nil {
		sku.UnitID = *req.UnitID
	}

	if req.Category != nil {
		sku.Category = *req.Category
	}

	if req.Barcode != nil {
		sku.Barcode = *req.Barcode
	}

	if req.ExpirationPeriod != nil {
		sku.ExpirationPeriodInDays = *req.ExpirationPeriod
	}

	if req.IsActive != nil {
		sku.IsActive = *req.IsActive
	}

	sku.UpdatedAt = time.Now()

	return nil
}
