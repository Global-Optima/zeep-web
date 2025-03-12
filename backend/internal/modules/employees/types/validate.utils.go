package types

import (
	"fmt"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type UpdateEmployeeModels struct {
	Employee *data.Employee
	Workdays []data.EmployeeWorkday
}

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

	workdays := make([]data.EmployeeWorkday, len(input.Workdays))

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
		IsActive:       &input.IsActive,
		Workdays:       workdays,
	}

	return employee, nil
}

func PrepareUpdateFields(input *UpdateEmployeeDTO) (*UpdateEmployeeModels, error) {
	employee := &data.Employee{}
	if input.FirstName != nil {
		employee.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		employee.LastName = *input.LastName
	}
	if input.IsActive != nil {
		employee.IsActive = input.IsActive
	}

	workdays := make([]data.EmployeeWorkday, len(input.Workdays))
	if len(input.Workdays) > 0 {
		for i, workday := range input.Workdays {
			employeeWorkday, err := ValidateWorkday(&workday)
			if err != nil {
				return nil, err
			}
			workdays[i] = *employeeWorkday
		}
	}

	return &UpdateEmployeeModels{
		Employee: employee,
		Workdays: workdays,
	}, nil
}

func ValidateWorkday(dto *CreateOrReplaceWorkdayDTO) (*data.EmployeeWorkday, error) {
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
