package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func MapToEmployeeDTO(employee *data.Employee) *EmployeeDTO {
	dto := &EmployeeDTO{
		ID:        employee.ID,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Phone:     employee.Phone,
		Email:     employee.Email,
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
		WarehouseID: employee.StoreEmployee.StoreID,
	}

	return dto
}
