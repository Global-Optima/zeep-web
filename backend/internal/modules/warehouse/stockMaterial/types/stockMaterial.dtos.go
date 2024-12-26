package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type CreateStockMaterialDTO struct {
	Name                   string  `json:"name" binding:"required"`
	Description            string  `json:"description"`
	SafetyStock            float64 `json:"safetyStock" binding:"required,gt=0"`
	ExpirationFlag         bool    `json:"expirationFlag"`
	UnitID                 uint    `json:"unitId" binding:"required"`
	SupplierID             uint    `json:"supplierId" binding:"required"`
	CategoryID             uint    `json:"categoryId" binding:"required"`
	Barcode                string  `json:"barcode"`
	ExpirationPeriodInDays int     `json:"expirationPeriodInDays"`
}

type UpdateStockMaterialDTO struct {
	Name                   *string  `json:"name"`
	Description            *string  `json:"description"`
	SafetyStock            *float64 `json:"safetyStock" binding:"omitempty,gt=0"`
	ExpirationFlag         *bool    `json:"expirationFlag"`
	UnitID                 *uint    `json:"unitId"`
	CategoryID             *uint    `json:"categoryId" binding:"required"`
	Barcode                *string  `json:"barcode"`
	ExpirationPeriodInDays *int     `json:"expirationPeriodInDays"`
	IsActive               *bool    `json:"isActive"`
}

type StockMaterialsDTO struct {
	ID                     uint    `json:"id"`
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	SafetyStock            float64 `json:"safetyStock"`
	ExpirationFlag         bool    `json:"expirationFlag"`
	UnitID                 uint    `json:"unitId"`
	UnitName               string  `json:"unitName,omitempty"`
	Category               string  `json:"category"`
	Barcode                string  `json:"barcode"`
	ExpirationPeriodInDays int     `json:"expirationPeriodInDays"`
	IsActive               bool    `json:"isActive"`
	CreatedAt              string  `json:"createdAt"`
	UpdatedAt              string  `json:"updatedAt"`
}

type StockMaterialFilter struct {
	Search         *string `form:"search"`
	LowStock       *bool   `form:"lowStock"`
	ExpirationFlag *bool   `form:"expirationFlag"`
	IsActive       *bool   `form:"isActive"`
	SupplierID     *uint   `form:"supplierId"`
	Pagination     *utils.Pagination
}
