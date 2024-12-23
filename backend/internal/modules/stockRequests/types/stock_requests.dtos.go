package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type CreateStockRequestDTO struct {
	StoreID uint                        `json:"storeId" binding:"required"`
	Items   []CreateStockRequestItemDTO `json:"items" binding:"required"`
}

type CreateStockRequestItemDTO struct {
	StockMaterialID uint    `json:"stockMaterialId" binding:"required"`
	Quantity        float64 `json:"quantity" binding:"required,gt=0"`
}

type UpdateStockRequestStatusDTO struct {
	Status data.StockRequestStatus `json:"status" binding:"required,oneof=CREATED IN_DELIVERY PROCESSED COMPLETED REJECTED"`
}

type UpdateIngredientDates struct {
	DeliveredDate  time.Time
	ExpirationDate time.Time
}

type StockRequestResponse struct {
	RequestID     uint                       `json:"requestId"`
	Store         StoreDTO                   `json:"store"`
	WarehouseID   uint                       `json:"warehouseId"`
	WarehouseName string                     `json:"warehouseName"`
	Status        data.StockRequestStatus    `json:"status"`
	Items         []StockRequestItemResponse `json:"items"`
	CreatedAt     time.Time                  `json:"createdAt"`
	UpdatedAt     time.Time                  `json:"updatedAt"`
}

type StoreDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
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

type GetStockRequestsFilter struct {
	Pagination  *utils.Pagination
	StoreID     *uint      `form:"storeId"`
	WarehouseID *uint      `form:"warehouseId"`
	StartDate   *time.Time `form:"startDate"`
	EndDate     *time.Time `form:"endDate"`

	Statuses []data.StockRequestStatus `form:"statuses[]"`
}

type StockMaterialDTO struct {
	StockMaterialID uint    `json:"stockMaterialId"`
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Unit            string  `json:"unit"`
	AvailableQty    float64 `json:"availableQuantity"`
}

type StockMaterialAvailabilityDTO struct {
	StockMaterialID uint    `json:"stockMaterialId"`
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	AvailableQty    float64 `json:"availableQty"`
	WarehouseID     uint    `json:"warehouseId"`
	WarehouseName   string  `json:"warehouseName"`
	Unit            string  `json:"unit"`
}

type StockMaterialFilter struct {
	Category *string `json:"category,omitempty"`
	Search   *string `json:"search,omitempty"`
}
