package types

type FacilityAddressDTO struct {
	ID        uint    `json:"id"`
	Address   string  `json:"address"`
	Longitude float64 `json:"longitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
}

type StoreDTO struct {
	ID              uint                `json:"id"`
	Name            string              `json:"name"`
	IsFranchise     bool                `json:"isFranchise"`
	FacilityAddress *FacilityAddressDTO `json:"facilityAddress,omitempty"`
	ContactPhone    string              `json:"contactPhone"`
	ContactEmail    string              `json:"contactEmail"`
	StoreHours      string              `json:"storeHours"`
}
