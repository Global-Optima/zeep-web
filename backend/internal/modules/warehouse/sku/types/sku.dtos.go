package types

type CreateSKURequest struct {
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description"`
	SafetyStock      float64 `json:"safetyStock" binding:"required,gt=0"`
	ExpirationFlag   bool    `json:"expirationFlag"`
	Quantity         float64 `json:"quantity" binding:"required,gte=0"`
	UnitID           uint    `json:"unitId" binding:"required"`
	Category         string  `json:"category"`
	Barcode          string  `json:"barcode"`
	ExpirationPeriod int     `json:"expirationPeriod"` // in days, default is 1095 (3 years)
}

type UpdateSKURequest struct {
	Name             *string  `json:"name"`
	Description      *string  `json:"description"`
	SafetyStock      *float64 `json:"safetyStock" binding:"omitempty,gt=0"`
	ExpirationFlag   *bool    `json:"expirationFlag"`
	Quantity         *float64 `json:"quantity" binding:"omitempty,gte=0"`
	UnitID           *uint    `json:"unitId"`
	Category         *string  `json:"category"`
	Barcode          *string  `json:"barcode"`
	ExpirationPeriod *int     `json:"expirationPeriod"` // in days
	IsActive         *bool    `json:"isActive"`
}

type SKUResponse struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	SafetyStock      float64 `json:"safetyStock"`
	ExpirationFlag   bool    `json:"expirationFlag"`
	Quantity         float64 `json:"quantity"`
	UnitID           uint    `json:"unitId"`
	UnitName         string  `json:"unitName"`
	Category         string  `json:"category"`
	Barcode          string  `json:"barcode"`
	ExpirationPeriod int     `json:"expirationPeriod"`
	IsActive         bool    `json:"isActive"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

type SKUFilter struct {
	Name           *string
	Category       *string
	LowStock       *bool
	ExpirationFlag *bool
	IsActive       *bool
}
