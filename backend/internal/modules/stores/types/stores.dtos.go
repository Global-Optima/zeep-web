package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

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
	FacilityAddress UpdateFacilityAddressDTO `json:"facilityAddress"`
	IsActive        bool                     `json:"isActive"`
	ContactPhone    string                   `json:"contactPhone"`
	ContactEmail    string                   `json:"contactEmail"`
	StoreHours      string                   `json:"storeHours"`
}

type UpdateStoreDTO struct {
	Name            string                   `json:"name"`
	FranchiseID     *uint                    `json:"franchiseId"`
	FacilityAddress CreateFacilityAddressDTO `json:"facilityAddress"`
	IsActive        bool                     `json:"isActive"`
	ContactPhone    string                   `json:"contactPhone"`
	ContactEmail    string                   `json:"contactEmail"`
	StoreHours      string                   `json:"storeHours"`
}

type StoreDTO struct {
	ID              uint                `json:"id"`
	Name            string              `json:"name"`
	FranchiseID     *uint               `json:"franchiseId"`
	FacilityAddress *FacilityAddressDTO `json:"facilityAddress"`
	IsActive        bool                `json:"isActive"`
	ContactPhone    string              `json:"contactPhone"`
	ContactEmail    string              `json:"contactEmail"`
	StoreHours      string              `json:"storeHours"`
}

type StoreFilter struct {
	utils.BaseFilter
	IsFranchisee *bool   `form:"isFranchise"`
	Search       *string `form:"search"`
}
