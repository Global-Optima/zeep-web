package types

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

func ValidateEmployee(input CreateEmployeeDTO) error {
	if input.LastName == "" {
		return errors.New("employee name cannot be empty")
	}

	if !utils.IsValidEmail(input.Email) {
		return errors.New("invalid email format")
	}

	if err := utils.IsValidPassword(input.Password); err != nil {
		return fmt.Errorf("password validation failed: %v", err)
	}

	if data.IsValidEmployeeRole(input.Role) {
		return errors.New("invalid role specified")
	}

	return nil
}

func ValidateStoreEmployee(input *CreateStoreEmployeeDTO) error {
	if err := ValidateEmployee(input.CreateEmployeeDTO); err != nil {
		return err
	}

	if input.StoreID == 0 {
		return errors.New("store details are required for store employees")
	}

	return nil
}

func ValidateWarehouseEmployee(input *CreateWarehouseEmployeeDTO) error {
	if err := ValidateEmployee(input.CreateEmployeeDTO); err != nil {
		return err
	}

	if input.WarehouseID == 0 {
		return errors.New("warehouse details are required for warehouse employees")
	}

	return nil
}

func PrepareUpdateFields(input UpdateEmployeeDTO) (*data.Employee, error) {
	employee := &data.Employee{}

	if input.FirstName != nil {
		employee.FirstName = *input.FirstName
	}
	if input.Phone != nil {
		employee.Phone = *input.Phone
	}
	if input.Email != nil {
		if !utils.IsValidEmail(*input.Email) {
			return employee, errors.New("invalid email format")
		}
		employee.Email = *input.Email
	}
	if input.Role != nil {
		if !data.IsValidEmployeeRole(*input.Role) {
			return employee, fmt.Errorf("invalid role: %v", *input.Role)
		}
		employee.Role = *input.Role
	}

	return employee, nil
}

func StoreEmployeeUpdateFields(input *UpdateStoreEmployeeDTO) (*data.Employee, error) {
	employee, err := PrepareUpdateFields(input.UpdateEmployeeDTO)
	if err != nil {
		return employee, err
	}

	if input.StoreID != nil {
		employee.StoreEmployee = &data.StoreEmployee{
			StoreID: *input.StoreID,
		}
	}
	if input.IsFranchise != nil {
		if employee.StoreEmployee == nil {
			employee.StoreEmployee = &data.StoreEmployee{}
		}
		employee.StoreEmployee.IsFranchise = *input.IsFranchise
	}

	return employee, nil
}

func WarehouseEmployeeUpdateFields(input *UpdateWarehouseEmployeeDTO) (*data.Employee, error) {
	employee, err := PrepareUpdateFields(input.UpdateEmployeeDTO)
	if err != nil {
		return employee, err
	}

	if input.WarehouseID != nil {
		employee.WarehouseEmployee = &data.WarehouseEmployee{
			WarehouseID: *input.WarehouseID,
		}
	}

	return employee, nil
}
