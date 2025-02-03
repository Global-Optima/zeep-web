package types

import employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"

type AdminEmployeeDTO struct {
	ID uint `json:"id"`
	employeesTypes.BaseEmployeeDTO
	EmployeeID uint `json:"employeeId"`
}

type AdminEmployeeDetailsDTO struct {
	ID uint `json:"id"`
	employeesTypes.BaseEmployeeDetailsDTO
	EmployeeID uint `json:"employeeId"`
}
