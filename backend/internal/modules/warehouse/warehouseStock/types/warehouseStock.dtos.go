package types

import (
	"time"

	supplierTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/supplier/types"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type ReceiveWarehouseDelivery struct {
	SupplierID uint                            `json:"supplierId"`
	Materials  []ReceiveWarehouseStockMaterial `json:"materials"`
}

type ReceiveWarehouseStockMaterial struct {
	StockMaterialID uint    `json:"stockMaterialId"`
	Quantity        float64 `json:"quantity"`
}

type WarehouseDeliveryDTO struct {
	ID           uint                                `json:"id"`
	Supplier     supplierTypes.SupplierResponse      `json:"supplier"`
	Warehouse    warehouseTypes.WarehouseDTO         `json:"warehouse"`
	Materials    []WarehouseDeliveryStockMaterialDTO `json:"materials"`
	DeliveryDate time.Time                           `json:"deliveryDate"`
}

type WarehouseDeliveryStockMaterialDTO struct {
	StockMaterial  stockMaterialTypes.StockMaterialsDTO `json:"stockMaterial"`
	Quantity       float64                              `json:"quantity"`
	Barcode        string                               `json:"barcode"`
	ExpirationDate time.Time                            `json:"expirationDate"`
}

type WarehouseDeliveryFilter struct {
	utils.BaseFilter
	WarehouseID *uint      `form:"warehouseID"`
	StartDate   *time.Time `form:"startDate" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate     *time.Time `form:"endDate" time_format:"2006-01-02T15:04:05Z07:00"`
	Search      *string    `form:"search"`
}

// stocks
type GetWarehouseStockFilterQuery struct {
	utils.BaseFilter
	WarehouseID     *uint   `form:"warehouseId"`
	StockMaterialID *uint   `form:"stockMaterialId"`
	IngredientID    *uint   `form:"ingredientId"`
	LowStockOnly    *bool   `form:"lowStockOnly"`
	IsExpiring      *bool   `form:"isExpiring"`
	CategoryID      *uint   `form:"categoryId"`
	ExpirationDays  *int    `form:"daysToExpire"`
	Search          *string `form:"search"`
}

type UpdateWarehouseStockDTO struct {
	Quantity       *float64   `json:"quantity" binding:"omitempty,gte=0"`
	ExpirationDate *time.Time `json:"expirationDate" binding:"omitempty"`
}

type AddWarehouseStockMaterial struct {
	StockMaterialID uint    `json:"stockMaterialId" binding:"required"`
	Quantity        float64 `json:"quantity" binding:"gte=0"`
}

type AdjustWarehouseStock struct {
	WarehouseID     uint    `json:"warehouseId" binding:"required"`
	StockMaterialID uint    `json:"stockMaterialId" binding:"required"`
	Quantity        float64 `json:"quantity" binding:"required,gte=0"`
}

type WarehouseStockResponse struct {
	StockMaterial          StockMaterialResponse `json:"stockMaterial"`
	Quantity               float64               `json:"quantity"`
	EarliestExpirationDate *time.Time            `json:"earliestExpirationDate,omitempty"`
}

type StockMaterialResponse struct {
	stockMaterialTypes.StockMaterialsDTO
}

type AvailableStockMaterialFilter struct {
	utils.BaseFilter
	IngredientID    *uint   `form:"ingredientId"`
	StockMaterialID *uint   `form:"stockMaterialId"`
	Search          *string `form:"search"`
}
