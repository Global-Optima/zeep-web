package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	franchiseesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
)

func MapToStoreDTO(store *data.Store) *StoreDTO {

	facilityAddressDTO := MapToFacilityAddressDTO(&store.FacilityAddress)

	var isActive = false
	if store.IsActive != nil && *store.IsActive {
		isActive = true
	}
	var franchisee *franchiseesTypes.FranchiseeDTO = nil
	if store.Franchisee != nil {
		franchisee = franchiseesTypes.ConvertFranchiseeToDTO(store.Franchisee)
	}

	warehouse := warehouseTypes.ToWarehouseDTO(store.Warehouse)

	return &StoreDTO{
		ID:              store.ID,
		Name:            store.Name,
		Franchisee:      franchisee,
		IsActive:        isActive,
		Warehouse:       *warehouse,
		ContactPhone:    store.ContactPhone,
		ContactEmail:    store.ContactEmail,
		StoreHours:      store.StoreHours,
		FacilityAddress: facilityAddressDTO,
	}
}

func MapToFacilityAddressDTO(facilityAddress *data.FacilityAddress) *FacilityAddressDTO {
	return &FacilityAddressDTO{
		ID:        facilityAddress.ID,
		Address:   facilityAddress.Address,
		Longitude: safeFloat(facilityAddress.Longitude),
		Latitude:  safeFloat(facilityAddress.Latitude),
	}
}

func MapToFacilityAddressModel(dto *CreateOrUpdateFacilityAddressDTO) *data.FacilityAddress {
	return &data.FacilityAddress{
		Address:   dto.Address,
		Longitude: dto.Longitude,
		Latitude:  dto.Latitude,
	}
}

func safeFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
