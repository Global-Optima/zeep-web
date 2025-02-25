package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func MapToCustomerDTO(customer *data.Customer) *CustomerAdminDTO {
	isVerified, isBanned := false, false

	if customer.IsVerified != nil {
		isVerified = *customer.IsVerified
	}

	if customer.IsBanned != nil {
		isBanned = *customer.IsBanned
	}

	return &CustomerAdminDTO{
		CustomerDTO: CustomerDTO{
			ID:        customer.ID,
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Phone:     customer.Phone,
		},
		IsVerified: isVerified,
		IsBanned:   isBanned,
	}
}
