package types

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
	Status  string `json:"status"`
}
