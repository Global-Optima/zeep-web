package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

func MapEmployeeToClaimsData(employee *data.Employee) (*EmployeeClaimsData, error) {
	var workplaceID uint

	//TODO validate case when no subgroup attached
	switch employee.Type {
	case data.StoreEmployeeType:
		workplaceID = employee.StoreEmployee.StoreID
	case data.WarehouseEmployeeType:
		workplaceID = employee.WarehouseEmployee.WarehouseID
	case data.RegionManagerEmployeeType:
		workplaceID = employee.RegionManager.RegionID
	case data.FranchiseeEmployeeType:
		workplaceID = employee.FranchiseeEmployee.FranchiseeID
	case data.AdminEmployeeType:
		workplaceID = 0
	default:
		return nil, types.ErrUnsupportedEmployeeType
	}

	employeeData := EmployeeClaimsData{
		ID:           employee.ID,
		Role:         employee.Role,
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
