package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func MapToCustomerDTO(customer *data.Customer) *CustomerAdminDTO {
	return &CustomerAdminDTO{
		CustomerDTO: CustomerDTO{
			ID:        customer.ID,
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Phone:     customer.Phone,
		},
		IsVerified: customer.IsVerified,
		IsBanned:   customer.IsBanned,
	}
}
