package types

import (
	"time"

	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type CreateSupplierDTO struct {
	Name         string `json:"name" validate:"required"`
	ContactEmail string `json:"contactEmail" validate:"email"`
	ContactPhone string `json:"contactPhone" binding:"required"`
	City         string `json:"city" binding:"required"`
	Address      string `json:"address,omitempty"`
}

type UpdateSupplierDTO struct {
	Name         *string `json:"name,omitempty"`
	ContactEmail *string `json:"contactEmail,omitempty"`
	ContactPhone *string `json:"contactPhone,omitempty"`
	City         *string `json:"city,omitempty"`
	Address      *string `json:"address,omitempty"`
}

type SupplierResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	ContactEmail string    `json:"contactEmail"`
	ContactPhone string    `json:"contactPhone"`
	City         string    `json:"city"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UpsertSupplierMaterialsDTO struct {
	Materials []UpdateSupplierMaterialDTO `json:"materials" binding:"required"`
}

type UpdateSupplierMaterialDTO struct {
	StockMaterialID uint    `json:"stockMaterialId" binding:"required"`
	BasePrice       float64 `json:"basePrice" binding:"required,gt=0"`
}

type SupplierMaterialResponse struct {
	StockMaterial SupplierStockMaterialDTO `json:"stockMaterial"`
	BasePrice     float64                  `json:"basePrice"`
}

type SupplierStockMaterialDTO struct {
	stockMaterialTypes.StockMaterialsDTO
	utils.PackageMeasure `json:"packageMeasures"`
}

type SuppliersFilter struct {
	Search *string `form:"search"`
	utils.BaseFilter
}
