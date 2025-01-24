package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

func MapToEmployeeDTO(employee *data.Employee) *EmployeeDTO {
	dto := &EmployeeDTO{
		ID:        employee.ID,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Phone:     utils.FormatPhoneOutput(employee.Phone),
		Email:     employee.Email,
		Type:      employee.Type,
		Role:      employee.Role,
		IsActive:  employee.IsActive,
	}

	return dto
}

func MapToStoreEmployeeDTO(employee *data.Employee) *StoreEmployeeDTO {
	dto := &StoreEmployeeDTO{
		EmployeeDTO: *MapToEmployeeDTO(employee),
		StoreID:     employee.StoreEmployee.StoreID,
	}

	return dto
}

func MapToWarehouseEmployeeDTO(employee *data.Employee) *WarehouseEmployeeDTO {
	dto := &WarehouseEmployeeDTO{
		EmployeeDTO: *MapToEmployeeDTO(employee),
		WarehouseID: employee.WarehouseEmployee.WarehouseID,
	}

	return dto
}

func MapToFranchiseeEmployeeDTO(employee *data.Employee) *FranchiseeEmployeeDTO {
	dto := &FranchiseeEmployeeDTO{
		EmployeeDTO:  *MapToEmployeeDTO(employee),
		FranchiseeID: employee.FranchiseeEmployee.FranchiseeID,
	}

	return dto
}

func MapToRegionManagerDTO(employee *data.Employee) *RegionManagerDTO {
	dto := &RegionManagerDTO{
		EmployeeDTO: *MapToEmployeeDTO(employee),
		RegionID:    employee.RegionManager.RegionID,
	}

	return dto
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

func CreateToStoreEmployee(dto *CreateStoreEmployeeDTO) (*data.Employee, error) {
	employee, err := CreateToEmployee(&dto.CreateEmployeeDTO)
	if err != nil {
		return nil, err
	}

	employee.Type = data.StoreEmployeeType
	if !data.IsAllowableRole(employee.Type, dto.Role) {
		return nil, fmt.Errorf("%w: role %s is not allowed to create store employee", ErrEmployeeTypeAndRoleMismatch, dto.Role)
	}

	employee.StoreEmployee = &data.StoreEmployee{
		StoreID: dto.StoreID,
	}

	return employee, nil
}

func CreateToWarehouseEmployee(dto *CreateWarehouseEmployeeDTO) (*data.Employee, error) {
	employee, err := CreateToEmployee(&dto.CreateEmployeeDTO)
	if err != nil {
		return nil, err
	}

	employee.Type = data.WarehouseEmployeeType
	if !data.IsAllowableRole(employee.Type, dto.Role) {
		return nil, fmt.Errorf("%w: role %s is not allowed to create store employee", ErrEmployeeTypeAndRoleMismatch, dto.Role)
	}

	employee.WarehouseEmployee = &data.WarehouseEmployee{
		WarehouseID: dto.WarehouseID,
	}
	return employee, nil
}

func CreateToFranchiseeEmployee(dto *CreateFranchiseeEmployeeDTO) (*data.Employee, error) {
	employee, err := CreateToEmployee(&dto.CreateEmployeeDTO)
	if err != nil {
		return nil, err
	}

	employee.Type = data.FranchiseeEmployeeType
	if !data.IsAllowableRole(employee.Type, dto.Role) {
		return nil, fmt.Errorf("%w: role %s is not allowed to create franchisee employee", ErrEmployeeTypeAndRoleMismatch, dto.Role)
	}

	employee.FranchiseeEmployee = &data.FranchiseeEmployee{
		FranchiseeID: dto.FranchiseeID,
	}

	return employee, nil
}

func CreateToRegionManager(dto *CreateRegionManagerDTO) (*data.Employee, error) {
	employee, err := CreateToEmployee(&dto.CreateEmployeeDTO)
	if err != nil {
		return nil, err
	}

	employee.Type = data.WarehouseRegionManagerEmployeeType
	if !data.IsAllowableRole(employee.Type, dto.Role) {
		return nil, fmt.Errorf("%w: role %s is not allowed to create region manager", ErrEmployeeTypeAndRoleMismatch, dto.Role)
	}

	employee.RegionManager = &data.RegionManager{
		RegionID: dto.RegionID,
	}

	return employee, nil
}
