package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

func MapToEmployeeDTO(employee *data.Employee) *EmployeeDTO {
	var employeeType data.EmployeeType = ""
	var role data.EmployeeRole = ""
	switch {
	case employee.StoreEmployee != nil:
		employeeType = data.StoreEmployeeType
		role = employee.StoreEmployee.Role
	case employee.WarehouseEmployee != nil:
		employeeType = data.WarehouseEmployeeType
		role = employee.WarehouseEmployee.Role
	case employee.RegionEmployee != nil:
		employeeType = data.RegionEmployeeType
		role = employee.RegionEmployee.Role
	case employee.FranchiseeEmployee != nil:
		employeeType = data.FranchiseeEmployeeType
		role = employee.FranchiseeEmployee.Role
	case employee.AdminEmployee != nil:
		employeeType = data.AdminEmployeeType
		role = employee.AdminEmployee.Role
	}

	dto := &EmployeeDTO{
		ID:        employee.ID,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Phone:     utils.FormatPhoneOutput(employee.Phone),
		Email:     employee.Email,
		Type:      employeeType,
		Role:      role,
		IsActive:  employee.IsActive,
	}

	return dto
}

func CreateToEmployee(dto *CreateEmployeeDTO) (*data.Employee, error) {
	employee, err := ValidateEmployee(dto)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	employee.HashedPassword = hashedPassword

	return employee, nil
}

func MapToEmployeeAccountDTO(employee *data.Employee) *EmployeeAccountDTO {
	return &EmployeeAccountDTO{
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Email:     employee.Email,
	}
}

func MapToEmployeeWorkdayDTO(workday *data.EmployeeWorkday) *EmployeeWorkdayDTO {
	dto := &EmployeeWorkdayDTO{
		ID:         workday.ID,
		Day:        workday.Day.ToString(),
		StartAt:    workday.StartAt,
		EndAt:      workday.EndAt,
		EmployeeID: workday.EmployeeID,
	}
	return dto
}
