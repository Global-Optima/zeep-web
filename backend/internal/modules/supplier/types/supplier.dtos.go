package types

import "time"

type CreateSupplierDTO struct {
	Name         string `json:"name" validate:"required"`
	ContactEmail string `json:"contactEmail" validate:"email"`
	ContactPhone string `json:"contactPhone" binding:"required"`
	City         string `json:"city" binding:"required"`
	Address      string `json:"address,omitempty"`
}

type UpdateSupplierDTO struct {
	Name         *string `json:"name,omitempty"`
	ContactEmail *string `json:"contactEmail,omitempty"`
	ContactPhone *string `json:"contactPhone,omitempty"`
	City         *string `json:"city,omitempty"`
	Address      *string `json:"address,omitempty"`
}

type SupplierResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	ContactEmail string    `json:"contactEmail"`
	ContactPhone string    `json:"contactPhone"`
	City         string    `json:"city"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
