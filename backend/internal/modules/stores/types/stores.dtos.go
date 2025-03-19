package types

import (
	facilityAddressesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/facilityAddresses/types"
	franchiseesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type CreateStoreDTO struct {
	Name            string                                                  `json:"name" binding:"required"`
	FranchiseeID    *uint                                                   `json:"franchiseeId"`
	WarehouseID     uint                                                    `json:"warehouseId" binding:"required"`
	FacilityAddress facilityAddressesTypes.CreateOrUpdateFacilityAddressDTO `json:"facilityAddress" binding:"required"`
	IsActive        bool                                                    `json:"isActive" binding:"required"`
	ContactPhone    string                                                  `json:"contactPhone" binding:"required"`
	ContactEmail    string                                                  `json:"contactEmail" binding:"required"`
	StoreHours      string                                                  `json:"storeHours" binding:"required"`
}

type UpdateStoreDTO struct {
	Name            string                                                   `json:"name"`
	FranchiseeID    *uint                                                    `json:"franchiseeId"`
	WarehouseID     *uint                                                    `json:"warehouseId"`
	FacilityAddress *facilityAddressesTypes.CreateOrUpdateFacilityAddressDTO `json:"facilityAddress"`
	IsActive        *bool                                                    `json:"isActive"`
	ContactPhone    string                                                   `json:"contactPhone"`
	ContactEmail    string                                                   `json:"contactEmail"`
	StoreHours      string                                                   `json:"storeHours"`
}

type StoreDTO struct {
	ID              uint                                       `json:"id"`
	Name            string                                     `json:"name"`
	Franchisee      *franchiseesTypes.FranchiseeDTO            `json:"franchisee,omitempty"`
	Warehouse       warehouseTypes.WarehouseDTO                `json:"warehouse"`
	FacilityAddress *facilityAddressesTypes.FacilityAddressDTO `json:"facilityAddress"`
	IsActive        bool                                       `json:"isActive"`
	ContactPhone    string                                     `json:"contactPhone"`
	ContactEmail    string                                     `json:"contactEmail"`
	StoreHours      string                                     `json:"storeHours"`
}

type StoreFilter struct {
	utils.BaseFilter
	IsFranchisee *bool   `form:"isFranchise"`
	FranchiseeID *uint   `form:"franchiseeId"`
	WarehouseID  *uint   `json:"warehouseId"`
	Search       *string `form:"search"`
}
