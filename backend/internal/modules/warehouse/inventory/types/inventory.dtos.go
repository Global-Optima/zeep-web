package types

import "time"

type ReceiveInventoryRequest struct {
	WarehouseID uint            `json:"warehouseId" binding:"required"`
	Items       []InventoryItem `json:"items" binding:"required"`
}

type InventoryItem struct {
	SKU_ID   uint    `json:"skuId" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required,gte=0"`
}

type TransferInventoryRequest struct {
	SourceWarehouseID uint            `json:"sourceWarehouseId" binding:"required"`
	TargetWarehouseID uint            `json:"targetWarehouseId" binding:"required"`
	Items             []InventoryItem `json:"items" binding:"required"`
}

type PickupRequest struct {
	StoreWarehouseID uint            `json:"storeWarehouseId" binding:"required"`
	Items            []InventoryItem `json:"items" binding:"required"`
}

type InventoryLevelsResponse struct {
	WarehouseID uint            `json:"warehouseId"`
	Levels      []InventoryItem `json:"levels"`
}

type UpcomingExpirationResponse struct {
	SKU_ID         uint      `json:"skuId"`
	Name           string    `json:"name"`
	ExpirationDate time.Time `json:"expirationDate"`
	Quantity       float64   `json:"quantity"`
}

type ExtendExpirationRequest struct {
	DeliveryID        uint      `json:"deliveryId" binding:"required"`
	NewExpirationDate time.Time `json:"newExpirationDate" binding:"required"`
}
