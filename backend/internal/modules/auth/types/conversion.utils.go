package types

import (
	"fmt"

	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type EmployeeSession struct {
	EmployeeID   uint
	WorkplaceID  uint
	Role         data.EmployeeRole
	EmployeeType data.EmployeeType
}

type CustomerSession struct {
	CustomerID uint
	IsVerified bool
}

func MapEmployeeToEmployeeSessionData(employee *data.Employee) (*EmployeeSession, error) {
	// Declare them without initial assignment; Go defaults to zero values anyway.
	var workplaceID uint
	var role data.EmployeeRole

	switch employee.GetType() {
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
		return nil, fmt.Errorf("%w: %s", employeesTypes.ErrUnsupportedEmployeeType, employee.GetType())
	}

	employeeData := EmployeeSession{
		EmployeeID:   employee.ID,
		Role:         role,
		WorkplaceID:  workplaceID,
		EmployeeType: employee.GetType(),
	}

	return &employeeData, nil
}

func MapCustomerToClaimsData(customer *data.Customer) *CustomerSession {
	isVerified := false
	if customer.IsVerified == nil {
		isVerified = *customer.IsVerified
	}

	return &CustomerSession{
		CustomerID: customer.ID,
		IsVerified: isVerified,
	}
}
