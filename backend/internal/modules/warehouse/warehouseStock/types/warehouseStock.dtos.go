package types

import (
	"time"

	supplierTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/supplier/types"
	stockMaterialPackageTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialPackage/types"
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
	PackageID       uint    `json:"packageId"`
}

type TransferInventoryRequest struct {
	SourceWarehouseID uint                            `json:"sourceWarehouseId" binding:"required"`
	TargetWarehouseID uint                            `json:"targetWarehouseId" binding:"required"`
	Items             []ReceiveWarehouseStockMaterial `json:"items" binding:"required"`
}

type DeliveryResponse struct {
	ID             uint                              `json:"id"`
	Barcode        string                            `json:"barcode"`
	Material       WarehouseDeliveryStockMaterialDTO `json:"materials"`
	Supplier       supplierTypes.SupplierResponse    `json:"supplier"`
	Warehouse      warehouseTypes.WarehouseResponse  `json:"warehouse"`
	DeliveryDate   time.Time                         `json:"deliveryDate"`
	ExpirationDate time.Time                         `json:"expirationDate"`
}

type WarehouseDeliveryStockMaterialDTO struct {
	StockMaterial stockMaterialTypes.StockMaterialsDTO                   `json:"stockMaterial"`
	Package       stockMaterialPackageTypes.StockMaterialPackageResponse `json:"package"`
	Quantity      float64                                                `json:"quantity"`
}

type DeliveryFilter struct {
	WarehouseID      *uint      `form:"warehouseID"`
	StartDate        *time.Time `form:"startDate" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate          *time.Time `form:"endDate" time_format:"2006-01-02T15:04:05Z07:00"`
	SearchBySupplier *string    `form:"searchBySupplier"`
	utils.BaseFilter
}

// stocks
type GetWarehouseStockFilterQuery struct {
	WarehouseID     *uint   `form:"warehouseId"`
	StockMaterialID *uint   `form:"stockMaterialId"`
	IngredientID    *uint   `form:"ingredientId"`
	LowStockOnly    *bool   `form:"lowStockOnly"`
	IsExpiring      *bool   `form:"isExpiring"`
	CategoryID      *uint   `form:"categoryId"`
	ExpirationDays  *int    `form:"daysToExpire"`
	Search          *string `form:"search"`
	utils.BaseFilter
}

type UpdateWarehouseStockDTO struct {
	Quantity       *float64   `json:"quantity" binding:"omitempty,gt=0"`
	ExpirationDate *time.Time `json:"expirationDate" binding:"omitempty"`
}

type AddWarehouseStockMaterial struct {
	StockMaterialID uint    `json:"stockMaterialId" binding:"required"`
	Quantity        float64 `json:"quantity" binding:"required,gte=0"`
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

type WarehouseStockMaterialDetailsDTO struct {
	StockMaterial          stockMaterialTypes.StockMaterialsDTO `json:"stockMaterial"`
	Quantity               float64                              `json:"quantity"`
	EarliestExpirationDate *time.Time                           `json:"earliestExpirationDate,omitempty"`
	Deliveries             []StockMaterialDeliveryDTO           `json:"deliveries"`
}

type StockMaterialDeliveryDTO struct {
	Supplier       supplierTypes.SupplierResponse `json:"supplier"`
	Quantity       float64                        `json:"quantity"`
	DeliveryDate   time.Time                      `json:"deliveryDate"`
	ExpirationDate time.Time                      `json:"expirationDate"`
}
