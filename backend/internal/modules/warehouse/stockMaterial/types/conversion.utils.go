package types

import (
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
