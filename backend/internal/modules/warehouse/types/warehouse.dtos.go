package types

import (
	facilityAddressesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/facilityAddresses/types"
	regionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type AssignStoreToWarehouseRequest struct {
	StoreID     uint `json:"storeId" binding:"required"`
	WarehouseID uint `json:"warehouseId" binding:"required"`
}

type ListStoresResponse struct {
	StoreID uint   `json:"storeId"`
	Name    string `json:"name"`
}

type CreateWarehouseDTO struct {
	FacilityAddress facilityAddressesTypes.CreateOrUpdateFacilityAddressDTO `json:"facilityAddress" binding:"required"`
	RegionID        uint                                                    `json:"regionId" binding:"required"`
	Name            string                                                  `json:"name" binding:"required"`
}

type UpdateWarehouseDTO struct {
	Name            *string                                                  `json:"name"`
	RegionID        *uint                                                    `json:"regionId"`
	FacilityAddress *facilityAddressesTypes.CreateOrUpdateFacilityAddressDTO `json:"facilityAddress"`
}

type WarehouseDTO struct {
	ID              uint                                      `json:"id"`
	Name            string                                    `json:"name"`
	Region          regionsTypes.RegionDTO                    `json:"region"`
	FacilityAddress facilityAddressesTypes.FacilityAddressDTO `json:"facilityAddress"`
	CreatedAt       string                                    `json:"createdAt"`
	UpdatedAt       string                                    `json:"updatedAt"`
}

type WarehouseFilter struct {
	utils.BaseFilter
	RegionID *uint   `form:"regionId"`
	Search   *string `form:"search"`
}
