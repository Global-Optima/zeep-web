package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

type StoreUpdateModels struct {
	Store           *data.Store
	FacilityAddress *data.FacilityAddress
}

func MapToFacilityAddressDTO(facilityAddress *data.FacilityAddress) *FacilityAddressDTO {
	return &FacilityAddressDTO{
		ID:        facilityAddress.ID,
		Address:   facilityAddress.Address,
		Longitude: safeFloat(facilityAddress.Longitude),
		Latitude:  safeFloat(facilityAddress.Latitude),
	}
}

func MapToFacilityAddressModel(dto *CreateOrUpdateFacilityAddressDTO, address *data.FacilityAddress) *data.FacilityAddress {
	if address == nil {
		return &data.FacilityAddress{
			Address:   dto.Address,
			Longitude: dto.Longitude,
			Latitude:  dto.Latitude,
		}
	}

	address.Address = dto.Address
	address.Longitude = dto.Longitude
	address.Latitude = dto.Latitude

	return address
}

func safeFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
