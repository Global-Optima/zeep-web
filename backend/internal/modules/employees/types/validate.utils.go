package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"strings"
)

func ValidateEmployee(input *CreateEmployeeDTO) (*data.Employee, error) {
	if strings.TrimSpace(input.LastName) == "" || strings.TrimSpace(input.FirstName) == "" {
		return nil, fmt.Errorf("%w: employee name cannot contain empty values", ErrValidation)
	}

	if !utils.IsValidEmail(input.Email) {
		return nil, fmt.Errorf("%w: invalid email format", ErrValidation)
	}

	if input.Phone != "" {
		if !utils.IsValidPhone(input.Phone, "") {
			return nil, fmt.Errorf("%w: invalid phone number format", ErrValidation)
		}
	}

	if err := utils.IsValidPassword(input.Password); err != nil {
		return nil, fmt.Errorf("%w: password validation failed: %v", ErrValidation, err)
	}
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing password failed: %w", err)
	}

	if !data.IsValidEmployeeRole(input.Role) {
		return nil, fmt.Errorf("%w: invalid role specified: %s", ErrValidation, input.Role)
	}

	var workdays = make([]data.EmployeeWorkday, len(input.Workdays))

	for i, workday := range input.Workdays {
		validatedWorkday, err := ValidateWorkday(&workday)
		if err != nil {
			return nil, err
		}
		workdays[i] = *validatedWorkday
	}

	employee := &data.Employee{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Phone:          utils.FormatPhoneInput(input.Phone),
		HashedPassword: hashedPassword,
		Role:           input.Role,
		IsActive:       input.IsActive,
		Workdays:       workdays,
	}

	return employee, nil
}

func PrepareUpdateFields(input UpdateEmployeeDTO) (*data.Employee, error) {
	employee := &data.Employee{}

	if input.FirstName != nil {
		employee.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		employee.LastName = *input.LastName
	}
	if input.IsActive != nil {
		employee.IsActive = *input.IsActive
	}
	if input.Phone != nil {
		if utils.IsValidPhone(*input.Phone, "") {
			employee.Phone = utils.FormatPhoneInput(*input.Phone)
		} else {
			return nil, fmt.Errorf("%w: invalid phone number format: %s", ErrValidation, *input.Phone)
		}
	}
	if input.Email != nil {
		if !utils.IsValidEmail(*input.Email) {
			return employee, fmt.Errorf("%w: invalid email format", ErrValidation)
		}
		employee.Email = *input.Email
	}
	if input.Role != nil {
		if !data.IsValidEmployeeRole(*input.Role) {
			return employee, fmt.Errorf("%w: invalid role: %v", ErrValidation, *input.Role)
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

func ValidateWorkday(dto *CreateWorkdayDTO) (*data.EmployeeWorkday, error) {
	var workday data.EmployeeWorkday
	weekday, err := data.ToWeekday(dto.Day)
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %s", ErrValidation, ErrInvalidWeekdayFormat, dto.Day)
	}
	workday.Day = weekday

	if err := utils.ValidateTime(dto.StartAt); err != nil {
		return nil, fmt.Errorf("%w: start at validation failed: %v", ErrValidation, err)
	}
	workday.StartAt = dto.StartAt

	if err := utils.ValidateTime(dto.EndAt); err != nil {
		return nil, fmt.Errorf("%w: end at validation failed: %v", ErrValidation, err)
	}
	workday.EndAt = dto.EndAt

	return &workday, nil
}

func ValidateEmployeeWorkday(input *CreateEmployeeWorkdayDTO) (*data.EmployeeWorkday, error) {
	workday, err := ValidateWorkday(&CreateWorkdayDTO{
		Day:     input.Day,
		StartAt: input.StartAt,
		EndAt:   input.EndAt,
	})
	if err != nil {
		return nil, err
	}

	if input.EmployeeID == 0 {
		return nil, fmt.Errorf("%w: employee ID cannot be zero", ErrValidation)
	}
	workday.EmployeeID = input.EmployeeID

	return workday, nil
}

func WorkdaysUpdateFields(input *UpdateEmployeeWorkdayDTO) (*data.EmployeeWorkday, error) {
	workday := &data.EmployeeWorkday{}

	if input.StartAt != nil {
		if err := utils.ValidateTime(*input.StartAt); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrValidation, err)
		}
		workday.StartAt = *input.StartAt
	}

	if input.EndAt != nil {
		if err := utils.ValidateTime(*input.EndAt); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrValidation, err)
		}
		workday.EndAt = *input.EndAt
	}

	return workday, nil
}
