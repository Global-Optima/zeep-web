package types

import "time"

type CreateSupplierDTO struct {
	Name         string `json:"name" validate:"required"`
	ContactEmail string `json:"contactEmail,omitempty" validate:"email"`
	ContactPhone string `json:"contactPhone,omitempty"`
	Address      string `json:"address,omitempty"`
}

type UpdateSupplierDTO struct {
	Name         *string `json:"name,omitempty"`
	ContactEmail *string `json:"contactEmail,omitempty"`
	ContactPhone *string `json:"contactPhone,omitempty"`
	Address      *string `json:"address,omitempty"`
}

type SupplierResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	ContactEmail string    `json:"contactEmail,omitempty"`
	ContactPhone string    `json:"contactPhone,omitempty"`
	Address      string    `json:"address,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
