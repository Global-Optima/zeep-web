package types

import "time"

type ReceiveInventoryRequest struct {
	SupplierID  uint            `json:"supplierId" binding:"required"`
	WarehouseID uint            `json:"warehouseId" binding:"required"`
	Items       []InventoryItem `json:"items" binding:"required"`
}

type InventoryItem struct {
	SKU_ID         uint     `json:"skuId"`                             // For existing SKUs
	Name           *string  `json:"name"`                              // Required for new SKUs
	Description    *string  `json:"description,omitempty"`             // Optional for new SKUs
	SafetyStock    *float64 `json:"safetyStock,omitempty"`             // Required for new SKUs
	ExpirationFlag *bool    `json:"expirationFlag,omitempty"`          // Required for new SKUs
	Quantity       float64  `json:"quantity" binding:"required,gte=0"` // Quantity to log
	UnitID         *uint    `json:"unitId,omitempty"`                  // Required for new SKUs
	Category       *string  `json:"category,omitempty"`                // Optional for new SKUs
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

type DeliveryResponse struct {
	ID             uint      `json:"id"`
	SKU_ID         uint      `json:"skuId"`
	Source         uint      `json:"source"`
	Target         uint      `json:"target"`
	Barcode        string    `json:"barcode"`
	Quantity       float64   `json:"quantity"`
	DeliveryDate   time.Time `json:"deliveryDate"`
	ExpirationDate time.Time `json:"expirationDate"`
}
