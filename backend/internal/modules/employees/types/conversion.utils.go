package types

import (
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
		Role:      employee.Role,
		IsActive:  employee.IsActive,
		Type:      employee.Type,
	}

	return dto
}

func MapToStoreEmployeeDTO(employee *data.Employee) *StoreEmployeeDTO {
	dto := &StoreEmployeeDTO{
		EmployeeDTO: *MapToEmployeeDTO(employee),
		StoreID:     employee.StoreEmployee.StoreID,
		IsFranchise: employee.StoreEmployee.IsFranchise,
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

func CreateToEmployee(dto *CreateEmployeeDTO, hashedPassword string) *data.Employee {
	employee := &data.Employee{
		FirstName:      dto.FirstName,
		LastName:       dto.LastName,
		Phone:          utils.FormatPhoneInput(dto.Phone),
		Email:          dto.Email,
		Role:           dto.Role,
		HashedPassword: hashedPassword,
		IsActive:       dto.IsActive,
	}
	return employee
}

func CreateToStoreEmployee(dto *CreateStoreEmployeeDTO, hashedPassword string) *data.Employee {
	employee := CreateToEmployee(&dto.CreateEmployeeDTO, hashedPassword)

	employee.Type = data.StoreEmployeeType
	employee.StoreEmployee = &data.StoreEmployee{
		StoreID:     dto.StoreID,
		IsFranchise: dto.IsFranchise,
	}

	return employee
}

func CreateToWarehouseEmployee(dto *CreateWarehouseEmployeeDTO, hashedPassword string) *data.Employee {
	employee := CreateToEmployee(&dto.CreateEmployeeDTO, hashedPassword)

	employee.Type = data.WarehouseEmployeeType
	employee.WarehouseEmployee = &data.WarehouseEmployee{
		WarehouseID: dto.WarehouseID,
	}
	return employee
}
