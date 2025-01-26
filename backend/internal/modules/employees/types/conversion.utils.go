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
		IsActive:  employee.IsActive,
	}

	return dto
}

func MapToStoreEmployeeDTO(storeEmployee *data.StoreEmployee) *StoreEmployeeDTO {
	dto := &StoreEmployeeDTO{
		EmployeeDTO: *MapToEmployeeDTO(&storeEmployee.Employee),
		StoreID:     storeEmployee.StoreID,
		Role:        storeEmployee.Role.ToString(),
	}
	return dto
}

func MapToWarehouseEmployeeDTO(warehouseEmployee *data.WarehouseEmployee) *WarehouseEmployeeDTO {
	dto := &WarehouseEmployeeDTO{
		EmployeeDTO: *MapToEmployeeDTO(&warehouseEmployee.Employee),
		WarehouseID: warehouseEmployee.WarehouseID,
		Role:        warehouseEmployee.Role.ToString(),
	}
	return dto
}

func MapToFranchiseeEmployeeDTO(franchiseeEmployee *data.FranchiseeEmployee) *FranchiseeEmployeeDTO {
	dto := &FranchiseeEmployeeDTO{
		EmployeeDTO:  *MapToEmployeeDTO(&franchiseeEmployee.Employee),
		FranchiseeID: franchiseeEmployee.FranchiseeID,
		Role:         franchiseeEmployee.Role.ToString(),
	}
	return dto
}

func MapToRegionManagerDTO(regionManager *data.RegionManager) *RegionManagerDTO {
	dto := &RegionManagerDTO{
		EmployeeDTO: *MapToEmployeeDTO(&regionManager.Employee),
		RegionID:    regionManager.RegionID,
		Role:        regionManager.Role.ToString(),
	}

	return dto
}

func MapToAdminEmployeeDTO(admin *data.AdminEmployee) *AdminEmployeeDTO {
	dto := &AdminEmployeeDTO{
		EmployeeDTO: *MapToEmployeeDTO(&admin.Employee),
		Role:        admin.Role.ToString(),
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
	employee.FranchiseeEmployee = &data.FranchiseeEmployee{
		FranchiseeID: franchiseeID,
		Role:         dto.Role,
	}

	return employee, nil
}

func CreateToRegionManager(regionID uint, dto *CreateEmployeeDTO) (*data.Employee, error) {
	if !data.IsAllowableRole(data.WarehouseRegionManagerEmployeeType, dto.Role) {
		return nil, fmt.Errorf("%w: %s for type %s", ErrInvalidEmployeeRole, dto.Role, data.WarehouseRegionManagerEmployeeType)
	}

	employee, err := CreateToEmployee(dto)
	if err != nil {
		return nil, err
	}
	employee.RegionManager = &data.RegionManager{
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
	employee.AdminEmployee = &data.AdminEmployee{
		Role: dto.Role,
	}

	return employee, nil
}
