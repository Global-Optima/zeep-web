package types

import (
	"time"

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
	Description      *string `json:"description,omitempty"` // Optional
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

type PickupRequest struct {
	StoreWarehouseID uint                    `json:"storeWarehouseId" binding:"required"`
	Items            []ExistingInventoryItem `json:"items" binding:"required"`
}

type InventoryLevel struct {
	StockMaterialID uint    `json:"stockMaterialId"`
	Name            string  `json:"name"`
	Quantity        float64 `json:"quantity"`
}

type InventoryLevelsResponse struct {
	Levels []InventoryLevel `json:"levels"`
}

type UpcomingExpirationResponse struct {
	DeliveryID      uint      `json:"deliveryId"`
	StockMaterialID uint      `json:"stockMaterialId"`
	Name            string    `json:"name"`
	ExpirationDate  time.Time `json:"expirationDate"`
	Quantity        float64   `json:"quantity"`
}

type ExtendExpirationRequest struct {
	DeliveryID uint `json:"deliveryId" binding:"required"`
	AddDays    int  `json:"addDays" binding:"required,gte=0"`
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

type GetInventoryLevelsFilterQuery struct {
	Search      *string `form:"search"`
	WarehouseID *uint   `form:"warehouseId"`
	Pagination  *utils.Pagination
}
