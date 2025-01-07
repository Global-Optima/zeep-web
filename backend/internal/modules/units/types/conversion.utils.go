package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func ToUnitResponse(unit data.Unit) UnitResponse {
	return UnitResponse{
		ID:               unit.ID,
		Name:             unit.Name,
		ConversionFactor: unit.ConversionFactor,
	}
}

func ToUnitResponses(units []data.Unit) []UnitResponse {
	responses := make([]UnitResponse, len(units))
	for i, unit := range units {
		responses[i] = ToUnitResponse(unit)
	}
	return responses
}
