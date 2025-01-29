package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func ConvertToListStoresResponse(stores []data.Store) []ListStoresResponse {
	response := make([]ListStoresResponse, len(stores))
	for i, store := range stores {
		response[i] = ListStoresResponse{
			StoreID: store.ID,
			Name:    store.Name,
		}
	}
	return response
}

func ToWarehouseResponse(warehouse data.Warehouse) *WarehouseResponse {
	return &WarehouseResponse{
		ID:   warehouse.ID,
		Name: warehouse.Name,
		FacilityAddress: FacilityAddressDTO{
			Address:   warehouse.FacilityAddress.Address,
			Longitude: warehouse.FacilityAddress.Longitude,
			Latitude:  warehouse.FacilityAddress.Latitude,
		},
		CreatedAt: warehouse.CreatedAt.String(),
		UpdatedAt: warehouse.UpdatedAt.String(),
	}
}

func ToFacilityAddressModel(dto FacilityAddressDTO) data.FacilityAddress {
	return data.FacilityAddress{
		Address:   dto.Address,
		Longitude: dto.Longitude,
		Latitude:  dto.Latitude,
	}
}

func ToWarehouseModel(dto CreateWarehouseDTO, facilityAddressID uint) data.Warehouse {
	return data.Warehouse{
		FacilityAddressID: facilityAddressID,
		RegionID:          dto.RegionID,
		Name:              dto.Name,
	}
}
