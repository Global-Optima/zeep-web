package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
)

type CreateStockMaterialPackageDTO struct {
	StockMaterialID uint    `json:"stockMaterialId" binding:"required"`
	Size            float64 `json:"size" binding:"required,gt=0"`
	UnitID          uint    `json:"unitId" binding:"required"`
}

type UpdateStockMaterialPackageDTO struct {
	StockMaterialID *uint    `json:"stockMaterialId,omitempty"`
	Size            *float64 `json:"size,omitempty"`
	UnitID          *uint    `json:"unitId,omitempty"`
}

type StockMaterialPackageResponse struct {
	ID            uint             `json:"id"`
	Size          float64          `json:"size"`
	Unit          types.UnitsDTO   `json:"unit"`
	StockMaterial StockMaterialDTO `json:"stockMaterial"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

type StockMaterialDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
