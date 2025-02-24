package types

import (
	franchiseesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type FacilityAddressDTO struct {
	ID        uint    `json:"id"`
	Address   string  `json:"address"`
	Longitude float64 `json:"longitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
}

type CreateOrUpdateFacilityAddressDTO struct {
	Address   string   `json:"address"`
	Longitude *float64 `json:"longitude,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
}

type CreateStoreDTO struct {
	Name            string                           `json:"name" binding:"required"`
	FranchiseeID    *uint                            `json:"franchiseeId"`
	WarehouseID     uint                             `json:"warehouseId" binding:"required"`
	FacilityAddress CreateOrUpdateFacilityAddressDTO `json:"facilityAddress" binding:"required"`
	IsActive        bool                             `json:"isActive" binding:"required"`
	ContactPhone    string                           `json:"contactPhone" binding:"required"`
	ContactEmail    string                           `json:"contactEmail" binding:"required"`
	StoreHours      string                           `json:"storeHours" binding:"required"`
}

type UpdateStoreDTO struct {
	Name            string                            `json:"name"`
	FranchiseeID    *uint                             `json:"franchiseeId"`
	WarehouseID     *uint                             `json:"warehouseId"`
	FacilityAddress *CreateOrUpdateFacilityAddressDTO `json:"facilityAddress"`
	IsActive        *bool                             `json:"isActive"`
	ContactPhone    string                            `json:"contactPhone"`
	ContactEmail    string                            `json:"contactEmail"`
	StoreHours      string                            `json:"storeHours"`
}

type StoreDTO struct {
	ID              uint                            `json:"id"`
	Name            string                          `json:"name"`
	Franchisee      *franchiseesTypes.FranchiseeDTO `json:"franchisee,omitempty"`
	Warehouse       warehouseTypes.WarehouseDTO     `json:"warehouse"`
	FacilityAddress *FacilityAddressDTO             `json:"facilityAddress"`
	IsActive        bool                            `json:"isActive"`
	ContactPhone    string                          `json:"contactPhone"`
	ContactEmail    string                          `json:"contactEmail"`
	StoreHours      string                          `json:"storeHours"`
}

type StoreFilter struct {
	utils.BaseFilter
	IsFranchisee *bool   `form:"isFranchise"`
	FranchiseeID *uint   `form:"franchiseeId"`
	WarehouseID  *uint   `json:"warehouseId"`
	Search       *string `form:"search"`
}
