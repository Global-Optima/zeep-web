package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	facilityAddressesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/facilityAddresses/types"
	regionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
)

type WarehouseUpdateModels struct {
	Warehouse       *data.Warehouse
	FacilityAddress *data.FacilityAddress
}

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
		ID:              warehouse.ID,
		Name:            warehouse.Name,
		Region:          *regionsTypes.MapRegionToDTO(&warehouse.Region),
		FacilityAddress: *facilityAddressesTypes.MapToFacilityAddressDTO(&warehouse.FacilityAddress),
		CreatedAt:       warehouse.CreatedAt.String(),
		UpdatedAt:       warehouse.UpdatedAt.String(),
	}
}

func ToWarehouseModel(dto CreateWarehouseDTO, facilityAddressID uint) data.Warehouse {
	return data.Warehouse{
		FacilityAddressID: facilityAddressID,
		RegionID:          dto.RegionID,
		Name:              dto.Name,
	}
}

func UpdateWarehouseToModels(dto *UpdateWarehouseDTO) *WarehouseUpdateModels {
	warehouse := &data.Warehouse{}
	var facilityAddress *data.FacilityAddress

	if dto.FacilityAddress.Address != "" {
		warehouse.FacilityAddress.Address = dto.FacilityAddress.Address
	}

	if dto.FacilityAddress.Latitude != nil {
		warehouse.FacilityAddress.Latitude = dto.FacilityAddress.Latitude
	}

	if dto.FacilityAddress.Longitude != nil {
		warehouse.FacilityAddress.Longitude = dto.FacilityAddress.Longitude
	}

	if dto.RegionID != nil {
		warehouse.RegionID = *dto.RegionID
	}

	if dto.Name != nil {
		warehouse.Name = *dto.Name
	}

	if dto.FacilityAddress != nil {
		facilityAddress = facilityAddressesTypes.MapToFacilityAddressModel(dto.FacilityAddress)
	}

	return &WarehouseUpdateModels{
		Warehouse:       warehouse,
		FacilityAddress: facilityAddress,
	}
}
