package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func MapEmployeeToClaimsData(employee *data.Employee) (*EmployeeClaimsData, error) {
	var workplaceID uint = 0
	var role data.EmployeeRole = ""

	switch employee.Type {
	case data.StoreEmployeeType:
		workplaceID = employee.StoreEmployee.StoreID
		role = employee.StoreEmployee.Role
	case data.WarehouseEmployeeType:
		workplaceID = employee.WarehouseEmployee.WarehouseID
		role = employee.WarehouseEmployee.Role
	case data.RegionEmployeeType:
		workplaceID = employee.RegionEmployee.RegionID
		role = employee.RegionEmployee.Role
	case data.FranchiseeEmployeeType:
		workplaceID = employee.FranchiseeEmployee.FranchiseeID
		role = employee.FranchiseeEmployee.Role
	case data.AdminEmployeeType:
		workplaceID = 0
		role = employee.AdminEmployee.Role
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedEmployeeType, employee.Type)
	}

	employeeData := EmployeeClaimsData{
		ID:           employee.ID,
		Role:         role,
		WorkplaceID:  workplaceID,
		EmployeeType: employee.Type,
	}

	return &employeeData, nil
}

func MapCustomerToClaimsData(customer *data.Customer) *CustomerClaimsData {
	customerData := CustomerClaimsData{
		ID:         customer.ID,
		IsVerified: customer.IsVerified,
	}
	return &customerData
}
