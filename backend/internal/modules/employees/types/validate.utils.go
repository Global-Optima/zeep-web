package types

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

func ValidateEmployee(input CreateEmployeeDTO) error {
	if input.LastName == "" || input.FirstName == "" {
		return errors.New("employee name cannot contain empty values")
	}

	if !utils.IsValidEmail(input.Email) {
		return errors.New("invalid email format")
	}

	if input.Phone != "" {
		if !utils.IsValidPhone(input.Phone, "") {
			return errors.New("invalid phone number format")
		}
	}

	if err := utils.IsValidPassword(input.Password); err != nil {
		return fmt.Errorf("password validation failed: %v", err)
	}

	if !data.IsValidEmployeeRole(input.Role) {
		return fmt.Errorf("invalid role specified: %s", input.Role)
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
	if input.LastName != nil {
		employee.LastName = *input.LastName
	}
	if input.Phone != nil {
		if utils.IsValidPhone(*input.Phone, "") {
			employee.Phone = utils.FormatPhoneInput(*input.Phone)
		}
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

func ValidateEmployeeWorkday(input *CreateEmployeeWorkdayDTO) error {
	if !data.IsValidWeekday(input.Day) {
		return fmt.Errorf("not a valid weekday: %s", input.Day)
	}

	if err := utils.ValidateTime(input.StartAt); err != nil {
		return fmt.Errorf("start at validation failed: %v", err)
	}

	if err := utils.ValidateTime(input.EndAt); err != nil {
		return fmt.Errorf("end at validation failed: %v", err)
	}

	if input.EmployeeID == 0 {
		return fmt.Errorf("employee ID cannot be zero")
	}

	return nil
}

func WorkdaysUpdateFields(input *UpdateEmployeeWorkdayDTO) (*data.EmployeeWorkday, error) {
	workday := &data.EmployeeWorkday{}

	if input.StartAt != nil {
		if err := utils.ValidateTime(*input.StartAt); err != nil {
			return nil, err
		}
		workday.StartAt = *input.StartAt
	}

	if input.EndAt != nil {
		if err := utils.ValidateTime(*input.EndAt); err != nil {
			return nil, err
		}
		workday.EndAt = *input.EndAt
	}

	return workday, nil
}
