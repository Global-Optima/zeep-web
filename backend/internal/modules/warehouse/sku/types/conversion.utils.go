package types

import (
	"errors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ConvertCreateSKURequestToSKU(req *CreateSKURequest) *data.SKU {
	return &data.SKU{
		Name:             req.Name,
		Description:      req.Description,
		SafetyStock:      req.SafetyStock,
		ExpirationFlag:   req.ExpirationFlag,
		Quantity:         req.Quantity,
		UnitID:           req.UnitID,
		Category:         req.Category,
		Barcode:          req.Barcode,
		ExpirationPeriod: req.ExpirationPeriod,
		IsActive:         true,
	}
}

func ConvertSKUToSKUResponse(sku *data.SKU) *SKUResponse {
	return &SKUResponse{
		ID:               sku.ID,
		Name:             sku.Name,
		Description:      sku.Description,
		SafetyStock:      sku.SafetyStock,
		ExpirationFlag:   sku.ExpirationFlag,
		Quantity:         sku.Quantity,
		UnitID:           sku.UnitID,
		UnitName:         sku.Unit.Name,
		Category:         sku.Category,
		Barcode:          sku.Barcode,
		ExpirationPeriod: sku.ExpirationPeriod,
		IsActive:         sku.IsActive,
		CreatedAt:        sku.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        sku.UpdatedAt.Format(time.RFC3339),
	}
}

func ConvertUpdateSKURequestToSKU(sku *data.SKU, req *UpdateSKURequest) error {
	return ValidateAndApplyUpdate(sku, req)
}

func ConvertUpdateSKURequestToMap(req *UpdateSKURequest, updateFields map[string]interface{}) error {
	if req.Name != nil {
		if *req.Name == "" {
			return errors.New("SKU name cannot be empty")
		}
		updateFields["name"] = *req.Name
	}

	if req.Description != nil {
		updateFields["description"] = *req.Description
	}

	if req.SafetyStock != nil {
		if *req.SafetyStock <= 0 {
			return errors.New("SKU safety stock must be greater than zero")
		}
		updateFields["safety_stock"] = *req.SafetyStock
	}

	if req.ExpirationFlag != nil {
		updateFields["expiration_flag"] = *req.ExpirationFlag
	}

	if req.Quantity != nil {
		if *req.Quantity < 0 {
			return errors.New("SKU quantity cannot be negative")
		}
		updateFields["quantity"] = *req.Quantity
	}

	if req.UnitID != nil {
		updateFields["unit_id"] = *req.UnitID
	}

	if req.Category != nil {
		updateFields["category"] = *req.Category
	}

	if req.Barcode != nil {
		updateFields["barcode"] = *req.Barcode
	}

	if req.ExpirationPeriod != nil {
		updateFields["expiration_period"] = *req.ExpirationPeriod
	}

	if req.IsActive != nil {
		updateFields["is_active"] = *req.IsActive
	}

	if len(updateFields) == 0 {
		return errors.New("no valid fields provided for update")
	}

	// Update the updated_at field
	updateFields["updated_at"] = time.Now()

	return nil
}
