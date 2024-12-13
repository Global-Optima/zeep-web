package types

import (
	"errors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ConvertCreateStockMaterialRequestToStockMaterial(req *CreateStockMaterialRequest) *data.StockMaterial {
	return &data.StockMaterial{
		Name:                   req.Name,
		Description:            req.Description,
		SafetyStock:            req.SafetyStock,
		ExpirationFlag:         req.ExpirationFlag,
		UnitID:                 req.UnitID,
		Category:               req.Category,
		Barcode:                req.Barcode,
		ExpirationPeriodInDays: req.ExpirationPeriodInDays,
		IsActive:               true,
	}
}

func ConvertStockMaterialToStockMaterialResponse(stockMaterial *data.StockMaterial) *StockMaterialResponse {
	return &StockMaterialResponse{
		ID:                     stockMaterial.ID,
		Name:                   stockMaterial.Name,
		Description:            stockMaterial.Description,
		SafetyStock:            stockMaterial.SafetyStock,
		ExpirationFlag:         stockMaterial.ExpirationFlag,
		UnitID:                 stockMaterial.UnitID,
		UnitName:               stockMaterial.Unit.Name,
		Category:               stockMaterial.Category,
		Barcode:                stockMaterial.Barcode,
		ExpirationPeriodInDays: stockMaterial.ExpirationPeriodInDays,
		IsActive:               stockMaterial.IsActive,
		CreatedAt:              stockMaterial.CreatedAt.Format(time.RFC3339),
		UpdatedAt:              stockMaterial.UpdatedAt.Format(time.RFC3339),
	}
}

func ConvertUpdateStockMaterialRequestToStockMaterial(stockMaterial *data.StockMaterial, req *UpdateStockMaterialRequest) error {
	return ValidateAndApplyUpdate(stockMaterial, req)
}

func ConvertUpdateStockMaterialRequestToMap(req *UpdateStockMaterialRequest, updateFields map[string]interface{}) error {
	if req.Name != nil {
		if *req.Name == "" {
			return errors.New("StockMaterial name cannot be empty")
		}
		updateFields["name"] = *req.Name
	}

	if req.Description != nil {
		updateFields["description"] = *req.Description
	}

	if req.SafetyStock != nil {
		if *req.SafetyStock <= 0 {
			return errors.New("StockMaterial safety stock must be greater than zero")
		}
		updateFields["safety_stock"] = *req.SafetyStock
	}

	if req.ExpirationFlag != nil {
		updateFields["expiration_flag"] = *req.ExpirationFlag
	}

	if req.Quantity != nil {
		if *req.Quantity < 0 {
			return errors.New("StockMaterial quantity cannot be negative")
		}
		updateFields["quantity"] = *req.Quantity
	}

	if req.UnitID != nil {
		updateFields["unit_id"] = *req.UnitID
	}

	if req.SupplierID != nil {
		updateFields["supplier_id"] = *req.SupplierID
	}

	if req.Category != nil {
		updateFields["category"] = *req.Category
	}

	if req.Barcode != nil {
		updateFields["barcode"] = *req.Barcode
	}

	if req.ExpirationPeriodInDays != nil {
		updateFields["expiration_period"] = *req.ExpirationPeriodInDays
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
