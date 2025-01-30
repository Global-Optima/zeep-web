package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func MapToStoreDTO(store *data.Store) *StoreDTO {

	facilityAddressDTO := &FacilityAddressDTO{
		ID:        store.FacilityAddress.ID,
		Address:   store.FacilityAddress.Address,
		Longitude: safeFloat(store.FacilityAddress.Longitude),
		Latitude:  safeFloat(store.FacilityAddress.Latitude),
	}
	
	return &StoreDTO{
		ID:              store.ID,
		Name:            store.Name,
		FranchiseID:     store.FranchiseeID,
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
