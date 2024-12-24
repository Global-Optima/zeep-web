package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type AssignStoreToWarehouseRequest struct {
	StoreID     uint `json:"storeId" binding:"required"`
	WarehouseID uint `json:"warehouseId" binding:"required"`
}

type ReassignStoreRequest struct {
	WarehouseID uint `json:"warehouseId" binding:"required"`
}

type ListStoresResponse struct {
	StoreID uint   `json:"storeId"`
	Name    string `json:"name"`
}

type CreateWarehouseDTO struct {
	FacilityAddress FacilityAddressDTO `json:"facilityAddress" binding:"required"`
	Name            string             `json:"name" binding:"required"`
}

type UpdateWarehouseDTO struct {
	Name string `json:"name" binding:"required"`
}

type FacilityAddressDTO struct {
	Address   string   `json:"address" binding:"required"`
	Longitude *float64 `json:"longitude,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
}

type WarehouseResponse struct {
	ID              uint               `json:"id"`
	Name            string             `json:"name"`
	FacilityAddress FacilityAddressDTO `json:"facilityAddress"`
	CreatedAt       string             `json:"createdAt"`
	UpdatedAt       string             `json:"updatedAt"`
}

// stock related dtos
type AdjustWarehouseStockRequest struct {
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
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	SafetyStock float64 `json:"safetyStock"`
	Unit        string  `json:"unit"`
	Barcode     string  `json:"barcode"`
}

type ResetWarehouseStockRequest struct {
	WarehouseID uint    `json:"warehouseId" binding:"required"`
	Stocks      []Stock `json:"stocks" binding:"required"`
}

type Stock struct {
	StockMaterialID uint    `json:"stockMaterialId" binding:"required"`
	Quantity        float64 `json:"quantity" binding:"required,gte=0"`
}

type GetWarehouseStockFilterQuery struct {
	WarehouseID     *uint   `form:"warehouseId"`
	StockMaterialID *uint   `form:"stockMaterialId"`
	LowStockOnly    *bool   `form:"lowStockOnly"`
	Category        *string `form:"category"`
	ExpirationDays  *int    `form:"expirationDays"` // Number of days to expiration
	Search          *string `form:"search"`         // Search by stock material name
	Pagination      *utils.Pagination
}
