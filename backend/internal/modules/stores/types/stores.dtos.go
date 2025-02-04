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

type CreateFacilityAddressDTO struct {
	Address   string   `json:"address"`
	Longitude *float64 `json:"longitude,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
}

type UpdateFacilityAddressDTO struct {
	Address   string   `json:"address"`
	Longitude *float64 `json:"longitude,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
}

type CreateStoreDTO struct {
	Name            string                   `json:"name"`
	FranchiseID     *uint                    `json:"franchiseId"`
	WarehouseID     uint                     `json:"warehouseId"`
	FacilityAddress UpdateFacilityAddressDTO `json:"facilityAddress"`
	IsActive        bool                     `json:"isActive"`
	ContactPhone    string                   `json:"contactPhone"`
	ContactEmail    string                   `json:"contactEmail"`
	StoreHours      string                   `json:"storeHours"`
}

type UpdateStoreDTO struct {
	Name            string                   `json:"name"`
	FranchiseID     *uint                    `json:"franchiseId"`
	WarehouseID     *uint                    `json:"warehouseId"`
	FacilityAddress CreateFacilityAddressDTO `json:"facilityAddress"`
	IsActive        bool                     `json:"isActive"`
	ContactPhone    string                   `json:"contactPhone"`
	ContactEmail    string                   `json:"contactEmail"`
	StoreHours      string                   `json:"storeHours"`
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
	WarehouseID  *uint   `json:"warehouseId"`
	Search       *string `form:"search"`
}
