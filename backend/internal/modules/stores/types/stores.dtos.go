package types

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
	ContactPhone    string                   `json:"contactPhone"`
	ContactEmail    string                   `json:"contactEmail"`
	StoreHours      string                   `json:"storeHours"`
}

type UpdateStoreDTO struct {
	Name            string                   `json:"name"`
	FranchiseID     *uint                    `json:"franchiseId"`
	FacilityAddress CreateFacilityAddressDTO `json:"facilityAddress"`
	ContactPhone    string                   `json:"contactPhone"`
	ContactEmail    string                   `json:"contactEmail"`
	StoreHours      string                   `json:"storeHours"`
}

type StoreDTO struct {
	ID              uint                `json:"id"`
	Name            string              `json:"name"`
	FranchiseID     *uint               `json:"franchiseId"`
	FacilityAddress *FacilityAddressDTO `json:"facilityAddress"`
	ContactPhone    string              `json:"contactPhone"`
	ContactEmail    string              `json:"contactEmail"`
	StoreHours      string              `json:"storeHours"`
}
