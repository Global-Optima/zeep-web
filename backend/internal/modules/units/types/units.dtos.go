package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

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

type UnitFilter struct {
	utils.BaseFilter
	Search *string `form:"search"` // Search by name
}
