package types

type StoreWarehouseIngredientDTO struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	CurrentStock      float64 `json:"currentStock"`
	Unit              string  `json:"unit"`
	LowStockAlert     bool    `json:"lowStockAlert"`
	LowStockThreshold float64 `json:"minimumStockThreshold"`
}

type GetStoreWarehouseStockQuery struct {
	StoreID      uint    `json:"storeId" binding:"required"`
	SearchTerm   *string `json:"searchTerm,omitempty"`
	LowStockOnly *bool   `json:"lowStockAlert,omitempty"`
	Limit        int     `json:"limit,omitempty"`
	Offset       int     `json:"offset,omitempty"`
}

type UpdateStoreWarehouseIngredientDTO struct {
	CurrentStock      *float64 `json:"currentStock"`
	LowStockThreshold *float64 `json:"lowStockThreshold"`
}

type AddIngredientDTO struct {
	StoreID           uint     `json:"storeId" binding:"required"`
	IngredientID      uint     `json:"ingredientId" binding:"required"`
	CurrentStock      *float64 `json:"currentStock"`
	LowStockThreshold float64  `json:"lowStockAlert"`
}
