package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	facilityAddressesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/facilityAddresses/types"
	franchiseesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
)

func MapToStoreDTO(store *data.Store) *StoreDTO {
	facilityAddressDTO := facilityAddressesTypes.MapToFacilityAddressDTO(&store.FacilityAddress)

	var franchisee *franchiseesTypes.FranchiseeDTO = nil
	if store.Franchisee != nil {
		franchisee = franchiseesTypes.ConvertFranchiseeToDTO(store.Franchisee)
	}

	warehouse := warehouseTypes.ToWarehouseDTO(store.Warehouse)

	return &StoreDTO{
		ID:              store.ID,
		Name:            store.Name,
		Franchisee:      franchisee,
		IsActive:        store.IsActive,
		Warehouse:       *warehouse,
		ContactPhone:    store.ContactPhone,
		ContactEmail:    store.ContactEmail,
		StoreHours:      store.StoreHours,
		FacilityAddress: facilityAddressDTO,
	}
}
