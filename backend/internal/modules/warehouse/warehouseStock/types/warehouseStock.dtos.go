package types

import (
	"time"

	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type ReceiveInventoryRequest struct {
	SupplierID    uint                    `json:"supplierId" binding:"required"`
	WarehouseID   uint                    `json:"warehouseId" binding:"required"`
	NewItems      []NewInventoryItem      `json:"newItems,omitempty"`      // For new SKUs
	ExistingItems []ExistingInventoryItem `json:"existingItems,omitempty"` // For existing SKUs
}

type NewInventoryItem struct {
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description,omitempty"` // Optional
	SafetyStock      float64 `json:"safetyStock" binding:"required"`
	ExpirationFlag   bool    `json:"expirationFlag" binding:"required"`
	Quantity         float64 `json:"quantity" binding:"required,gte=0"`
	UnitID           uint    `json:"unitId" binding:"required"`
	CategoryID       uint    `json:"categoryId" binding:"required"`
	ExpirationInDays *int    `json:"expirationInDays,omitempty"` // Optional override
	Package          Package `json:"package" binding:"required"`
	IngredientID     uint    `json:"ingredientId" binding:"required"` // New field for ingredient linkage
}

type ExistingInventoryItem struct {
	StockMaterialID uint    `json:"stockMaterialId"`                   // For existing SKUs
	Quantity        float64 `json:"quantity" binding:"required,gte=0"` // Quantity to log
	IngredientID    uint    `json:"ingredientId" binding:"required"`   // Link to ingredient
}

type Package struct {
	Size   float64 `json:"size" binding:"required,gte=0"`
	UnitID uint    `json:"unitId" binding:"required"`
}

type TransferInventoryRequest struct {
	SourceWarehouseID uint                    `json:"sourceWarehouseId" binding:"required"`
	TargetWarehouseID uint                    `json:"targetWarehouseId" binding:"required"`
	Items             []ExistingInventoryItem `json:"items" binding:"required"`
}

type DeliveryResponse struct {
	ID              uint      `json:"id"`
	StockMaterialID uint      `json:"stockMaterialId"`
	SupplierID      uint      `json:"supplierId"`
	WarehouseID     uint      `json:"warehouseId"`
	Barcode         string    `json:"barcode"`
	Quantity        float64   `json:"quantity"`
	DeliveryDate    time.Time `json:"deliveryDate"`
	ExpirationDate  time.Time `json:"expirationDate"`
}

type DeliveryFilter struct {
	WarehouseID *uint      `form:"warehouseID"`
	StartDate   *time.Time `form:"startDate" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate     *time.Time `form:"endDate" time_format:"2006-01-02T15:04:05Z07:00"`
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

type ResetWarehouseStock struct {
	WarehouseID uint    `json:"warehouseId" binding:"required"`
	Stocks      []Stock `json:"stocks" binding:"required"`
}

type UpdateWarehouseStockDTO struct {
	Quantity       float64   `json:"quantity" binding:"required,gt=0"`
	ExpirationDate time.Time `json:"expirationDate" binding:"required"`
}

type Stock struct {
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
	TotalQuantity          float64               `json:"totalQuantity"`
	EarliestExpirationDate *time.Time            `json:"earliestExpirationDate,omitempty"`
}

type StockMaterialResponse struct {
	stockMaterialTypes.StockMaterialsDTO
	utils.PackageMeasure `json:"packageMeasures"`
}

type WarehouseStockMaterialDetailsDTO struct {
	StockMaterial          stockMaterialTypes.StockMaterialsDTO `json:"stockMaterial"`
	PackageMeasure         utils.PackageMeasure                 `json:"packageMeasure"`
	TotalQuantity          float64                              `json:"totalQuantity"`
	EarliestExpirationDate *time.Time                           `json:"earliestExpirationDate,omitempty"`
	Deliveries             []StockMaterialDeliveryDTO           `json:"deliveries"`
}

type StockMaterialDeliveryDTO struct {
	Supplier       string    `json:"supplierName"`
	Quantity       float64   `json:"quantity"`
	DeliveryDate   time.Time `json:"deliveryDate"`
	ExpirationDate time.Time `json:"expirationDate"`
}
