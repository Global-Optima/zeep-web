package types

type CreateUnitDTO struct {
	Name             string  `json:"name" binding:"required"`
	ConversionFactor float64 `json:"conversionFactor" binding:"required,gt=0"`
}

type UpdateUnitDTO struct {
	Name             *string  `json:"name,omitempty"`             // Optional field
	ConversionFactor *float64 `json:"conversionFactor,omitempty"` // Optional field
}

type UnitResponse struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	ConversionFactor float64 `json:"conversionFactor"`
}
