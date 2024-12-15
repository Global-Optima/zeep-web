package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func MapToCustomerDTO(customer *data.Customer) *CustomerDTO {
	return &CustomerDTO{
		ID:         customer.ID,
		Name:       customer.Name,
		Phone:      customer.Phone,
		IsVerified: customer.IsVerified,
		IsBanned:   customer.IsBanned,
	}
}
