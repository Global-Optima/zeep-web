package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	regionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
)

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

func ToWarehouseDTO(warehouse data.Warehouse) *WarehouseDTO {
	return &WarehouseDTO{
		ID:     warehouse.ID,
		Name:   warehouse.Name,
		Region: *regionsTypes.MapRegionToDTO(&warehouse.Region),
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

func UpdateWarehouseToModel(dto *UpdateWarehouseDTO) *data.Warehouse {
	warehouse := &data.Warehouse{}

	if dto.FacilityAddressID != nil {
		warehouse.FacilityAddressID = *dto.FacilityAddressID
	}
	if dto.RegionID != nil {
		warehouse.RegionID = *dto.RegionID
	}
	if dto.Name != nil {
		warehouse.Name = *dto.Name
	}

	return warehouse
}
