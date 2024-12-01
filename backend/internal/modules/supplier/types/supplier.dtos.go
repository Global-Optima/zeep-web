package types

import "time"

type CreateSupplierDTO struct {
	Name         string `json:"name" validate:"required"`
	ContactEmail string `json:"contact_email,omitempty" validate:"email"`
	ContactPhone string `json:"contact_phone,omitempty"`
	Address      string `json:"address,omitempty"`
}

type UpdateSupplierDTO struct {
	Name         *string `json:"name,omitempty"`
	ContactEmail *string `json:"contact_email,omitempty"`
	ContactPhone *string `json:"contact_phone,omitempty"`
	Address      *string `json:"address,omitempty"`
}

type SupplierResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	ContactEmail string    `json:"contact_email,omitempty"`
	ContactPhone string    `json:"contact_phone,omitempty"`
	Address      string    `json:"address,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
