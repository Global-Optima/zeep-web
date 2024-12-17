package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func MapEmployeeToClaimsData(employee *data.Employee) *EmployeeClaimsData {
	var workplaceID uint
	var workplaceType data.EmployeeType

	if employee.StoreEmployee != nil {
		workplaceID = employee.StoreEmployee.StoreID
		workplaceType = data.StoreEmployeeType
	} else if employee.WarehouseEmployee != nil {
		workplaceID = employee.WarehouseEmployee.WarehouseID
		workplaceType = data.WarehouseEmployeeType
	}

	employeeData := EmployeeClaimsData{
		ID:           employee.ID,
		Role:         employee.Role,
		WorkplaceID:  workplaceID,
		EmployeeType: workplaceType,
	}

	return &employeeData
}

func MapCustomerToClaimsData(customer *data.Customer) *CustomerClaimsData {
	customerData := CustomerClaimsData{
		ID:         customer.ID,
		IsVerified: customer.IsVerified,
	}
	return &customerData
}
