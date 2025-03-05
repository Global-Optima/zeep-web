package types

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
