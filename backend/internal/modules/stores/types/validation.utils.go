package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func UpdateStoreFields(store *data.Store, dto UpdateStoreDTO) {
	if dto.Name != "" {
		store.Name = dto.Name
	}
	if dto.FranchiseID != nil {
		store.FranchiseeID = dto.FranchiseID
	}
	if dto.WarehouseID != nil {
		store.WarehouseID = *dto.WarehouseID
	}
	if dto.ContactPhone != "" {
		store.ContactPhone = dto.ContactPhone
	}
	if dto.ContactEmail != "" {
		store.ContactEmail = dto.ContactEmail
	}
	if dto.StoreHours != "" {
		store.StoreHours = dto.StoreHours
	}
	if dto.FacilityAddress.Address != "" {
		store.FacilityAddress.Address = dto.FacilityAddress.Address
	}
	if dto.FacilityAddress.Latitude != nil {
		store.FacilityAddress.Latitude = dto.FacilityAddress.Latitude
	}
	if dto.FacilityAddress.Longitude != nil {
		store.FacilityAddress.Longitude = dto.FacilityAddress.Longitude
	}
}
