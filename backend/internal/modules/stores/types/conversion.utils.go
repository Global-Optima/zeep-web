package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	franchiseesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
)

func MapToStoreDTO(store *data.Store) *StoreDTO {

	facilityAddressDTO := &FacilityAddressDTO{
		ID:        store.FacilityAddress.ID,
		Address:   store.FacilityAddress.Address,
		Longitude: safeFloat(store.FacilityAddress.Longitude),
		Latitude:  safeFloat(store.FacilityAddress.Latitude),
	}

	var franchisee *franchiseesTypes.FranchiseeDTO = nil
	if store.Franchisee != nil {
		franchisee = franchiseesTypes.ConvertFranchiseeToDTO(store.Franchisee)
	}

	return &StoreDTO{
		ID:              store.ID,
		Name:            store.Name,
		Franchisee:      franchisee,
		ContactPhone:    store.ContactPhone,
		ContactEmail:    store.ContactEmail,
		StoreHours:      store.StoreHours,
		FacilityAddress: facilityAddressDTO,
	}
}

func safeFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
