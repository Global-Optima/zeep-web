package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

func MapToBaseEmployeeDTO(employee *data.Employee) *BaseEmployeeDTO {
	dto := &BaseEmployeeDTO{
		ID:        employee.ID,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Phone:     utils.FormatPhoneOutput(employee.Phone),
		Email:     employee.Email,
		Type:      employee.Type,
		IsActive:  employee.IsActive,
	}

	return dto
}

func MapToEmployeeDTO(employee *data.Employee) *EmployeeDTO {
	var role data.EmployeeRole = ""
	dto := MapToBaseEmployeeDTO(employee)

	switch {
	case employee.StoreEmployee != nil:
		role = employee.StoreEmployee.Role
	case employee.WarehouseEmployee != nil:
		role = employee.WarehouseEmployee.Role
	case employee.FranchiseeEmployee != nil:
		role = employee.FranchiseeEmployee.Role
	case employee.RegionEmployee != nil:
		role = employee.RegionEmployee.Role
	case employee.AdminEmployee != nil:
		role = employee.AdminEmployee.Role
	}

	return &EmployeeDTO{
		BaseEmployeeDTO: *dto,
		Role:            role,
	}
}

func MapToStoreEmployeeDTO(storeEmployee *data.StoreEmployee) *StoreEmployeeDTO {
	dto := &StoreEmployeeDTO{
		BaseEmployeeDTO: *MapToBaseEmployeeDTO(&storeEmployee.Employee),
		StoreID:         storeEmployee.StoreID,
		Role:            storeEmployee.Role,
	}
	return dto
}

func MapToWarehouseEmployeeDTO(warehouseEmployee *data.WarehouseEmployee) *WarehouseEmployeeDTO {
	dto := &WarehouseEmployeeDTO{
		BaseEmployeeDTO: *MapToBaseEmployeeDTO(&warehouseEmployee.Employee),
		WarehouseID:     warehouseEmployee.WarehouseID,
		Role:            warehouseEmployee.Role,
	}
	return dto
}

func MapToFranchiseeEmployeeDTO(franchiseeEmployee *data.FranchiseeEmployee) *FranchiseeEmployeeDTO {
	dto := &FranchiseeEmployeeDTO{
		BaseEmployeeDTO: *MapToBaseEmployeeDTO(&franchiseeEmployee.Employee),
		FranchiseeID:    franchiseeEmployee.FranchiseeID,
		Role:            franchiseeEmployee.Role,
	}
	return dto
}

func MapToRegionEmployeeDTO(regionManager *data.RegionEmployee) *RegionEmployeeDTO {
	dto := &RegionEmployeeDTO{
		BaseEmployeeDTO: *MapToBaseEmployeeDTO(&regionManager.Employee),
		RegionID:        regionManager.RegionID,
		Role:            regionManager.Role,
	}

	return dto
}

func MapToAdminEmployeeDTO(admin *data.AdminEmployee) *AdminEmployeeDTO {
	dto := &AdminEmployeeDTO{
		BaseEmployeeDTO: *MapToBaseEmployeeDTO(&admin.Employee),
		Role:            admin.Role,
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

func CreateToStoreEmployee(storeID uint, dto *CreateEmployeeDTO) (*data.Employee, error) {
	if !data.IsAllowableRole(data.StoreEmployeeType, dto.Role) {
		return nil, fmt.Errorf("%w: %s for type %s", ErrInvalidEmployeeRole, dto.Role, data.StoreEmployeeType)
	}

	employee, err := CreateToEmployee(dto)
	if err != nil {
		return nil, err
	}
	employee.Type = data.StoreEmployeeType
	employee.StoreEmployee = &data.StoreEmployee{
		StoreID: storeID,
		Role:    dto.Role,
	}

	return employee, nil
}

func CreateToWarehouseEmployee(warehouseID uint, dto *CreateEmployeeDTO) (*data.Employee, error) {
	if !data.IsAllowableRole(data.WarehouseEmployeeType, dto.Role) {
		return nil, fmt.Errorf("%w: %s for type %s", ErrInvalidEmployeeRole, dto.Role, data.WarehouseEmployeeType)
	}

	employee, err := CreateToEmployee(dto)
	if err != nil {
		return nil, err
	}
	employee.Type = data.WarehouseEmployeeType
	employee.WarehouseEmployee = &data.WarehouseEmployee{
		WarehouseID: warehouseID,
		Role:        dto.Role,
	}

	return employee, nil
}

func CreateToFranchiseeEmployee(franchiseeID uint, dto *CreateEmployeeDTO) (*data.Employee, error) {
	if !data.IsAllowableRole(data.FranchiseeEmployeeType, dto.Role) {
		return nil, fmt.Errorf("%w: %s for type %s", ErrInvalidEmployeeRole, dto.Role, data.FranchiseeEmployeeType)
	}

	employee, err := CreateToEmployee(dto)
	if err != nil {
		return nil, err
	}
	employee.Type = data.FranchiseeEmployeeType
	employee.FranchiseeEmployee = &data.FranchiseeEmployee{
		FranchiseeID: franchiseeID,
		Role:         dto.Role,
	}

	return employee, nil
}

func CreateToRegionEmployee(regionID uint, dto *CreateEmployeeDTO) (*data.Employee, error) {
	if !data.IsAllowableRole(data.RegionEmployeeType, dto.Role) {
		return nil, fmt.Errorf("%w: %s for type %s", ErrInvalidEmployeeRole, dto.Role, data.RegionEmployeeType)
	}

	employee, err := CreateToEmployee(dto)
	if err != nil {
		return nil, err
	}
	employee.Type = data.RegionEmployeeType
	employee.RegionEmployee = &data.RegionEmployee{
		RegionID: regionID,
		Role:     dto.Role,
	}

	return employee, nil
}

func CreateToAdminEmployee(dto *CreateEmployeeDTO) (*data.Employee, error) {
	if !data.IsAllowableRole(data.AdminEmployeeType, dto.Role) {
		return nil, fmt.Errorf("%w: %s for type %s", ErrInvalidEmployeeRole, dto.Role, data.AdminEmployeeType)
	}

	employee, err := CreateToEmployee(dto)
	if err != nil {
		return nil, err
	}
	employee.Type = data.AdminEmployeeType
	employee.AdminEmployee = &data.AdminEmployee{
		Role: dto.Role,
	}

	return employee, nil
}
