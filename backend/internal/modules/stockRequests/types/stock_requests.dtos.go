package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type CreateStockRequestDTO struct {
	StoreID uint                  `json:"storeId" binding:"required"`
	Items   []StockRequestItemDTO `json:"items" binding:"required"`
}

type StockRequestItemDTO struct {
	StockMaterialID uint    `json:"stockMaterialId" binding:"required"`
	Quantity        float64 `json:"quantity" binding:"required,gt=0"`
}

type UpdateStockRequestStatusDTO struct {
	Status data.StockRequestStatus `json:"status" binding:"required,oneof=CREATED IN_DELIVERY PROCESSED COMPLETED REJECTED"`
}

type StockRequestResponse struct {
	RequestID     uint                       `json:"requestId"`
	StoreID       uint                       `json:"storeId"`
	StoreName     string                     `json:"storeName"`
	WarehouseID   uint                       `json:"warehouseId"`
	WarehouseName string                     `json:"warehouseName"`
	Status        data.StockRequestStatus    `json:"status"`
	Items         []StockRequestItemResponse `json:"items"`
	CreatedAt     time.Time                  `json:"createdAt"`
	UpdatedAt     time.Time                  `json:"updatedAt"`
}

type StockRequestItemResponse struct {
	StockMaterialID uint    `json:"stockMaterialId"`
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Unit            string  `json:"unit"`
	Quantity        float64 `json:"quantity"`
}

type LowStockIngredientResponse struct {
	IngredientID      uint    `json:"ingredientId"`
	Name              string  `json:"name"`
	Unit              string  `json:"unit"`
	Quantity          float64 `json:"quantity"`
	LowStockThreshold float64 `json:"lowStockThreshold"`
}

type StockRequestFilter struct {
	StoreID     *uint      `json:"storeId,omitempty"`
	WarehouseID *uint      `json:"warehouseId,omitempty"`
	Status      *string    `json:"status,omitempty"`
	StartDate   *time.Time `json:"startDate,omitempty"`
	EndDate     *time.Time `json:"endDate,omitempty"`
}

type ProductMarketplaceDTO struct {
	StockMaterialID uint    `json:"stockMaterialId"`
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Unit            string  `json:"unit"`
	AvailableQty    float64 `json:"availableQuantity"`
	Price           float64 `json:"price,omitempty"`
}

type MarketplaceFilter struct {
	Category *string `json:"category,omitempty"`
	Search   *string `json:"search,omitempty"`
}
