package types

type CreateStockMaterialRequest struct {
	Name                   string  `json:"name" binding:"required"`
	Description            string  `json:"description"`
	SafetyStock            float64 `json:"safetyStock" binding:"required,gt=0"`
	ExpirationFlag         bool    `json:"expirationFlag"`
	UnitID                 uint    `json:"unitId" binding:"required"`
	SupplierID             uint    `json:"supplierId" binding:"required"`
	CategoryID             uint    `json:"categoryId" binding:"required"`
	Barcode                string  `json:"barcode"`
	ExpirationPeriodInDays int     `json:"expirationPeriodInDays"` // in days, default is 1095 (3 years)
}

type UpdateStockMaterialRequest struct {
	Name                   *string  `json:"name"`
	Description            *string  `json:"description"`
	SafetyStock            *float64 `json:"safetyStock" binding:"omitempty,gt=0"`
	ExpirationFlag         *bool    `json:"expirationFlag"`
	UnitID                 *uint    `json:"unitId"`
	CategoryID             *uint    `json:"categoryId" binding:"required"`
	Barcode                *string  `json:"barcode"`
	ExpirationPeriodInDays *int     `json:"expirationPeriodInDays"` // in days
	IsActive               *bool    `json:"isActive"`
}

type StockMaterialResponse struct {
	ID                     uint    `json:"id"`
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	SafetyStock            float64 `json:"safetyStock"`
	ExpirationFlag         bool    `json:"expirationFlag"`
	UnitID                 uint    `json:"unitId"`
	UnitName               string  `json:"unitName,omitempty"`
	Category               string  `json:"category"`
	Barcode                string  `json:"barcode"`
	ExpirationPeriodInDays int     `json:"expirationPeriodInDays"` // in days
	IsActive               bool    `json:"isActive"`
	CreatedAt              string  `json:"createdAt"`
	UpdatedAt              string  `json:"updatedAt"`
}

type StockMaterialFilter struct {
	Name           *string
	Category       *string
	LowStock       *bool
	ExpirationFlag *bool
	IsActive       *bool
	SupplierID     *uint
}
