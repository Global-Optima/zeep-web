package types

type FacilityAddressDTO struct {
	ID        uint   `json:"id"`
	Address   string `json:"address"`
	Longitude string `json:"longitude,omitempty"`
	Latitude  string `json:"latitude,omitempty"`
}

type StoreDTO struct {
	ID              uint                `json:"id"`
	Name            string              `json:"name"`
	IsFranchise     bool                `json:"isFranchise"`
	FacilityAddress *FacilityAddressDTO `json:"facilityAddress,omitempty"`
	Status          string              `json:"status"`
}

type EmployeeDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	IsActive bool   `json:"isActive"`
	Role     string `json:"role,omitempty"`
}
