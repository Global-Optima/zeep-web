package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	franchiseesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
)

type UpdateFranchiseeEmployeeDTO struct {
	employeesTypes.UpdateEmployeeDTO
	Role         *data.EmployeeRole `json:"role,omitempty"`
	FranchiseeID *uint              `json:"franchiseeId,omitempty"`
}

type FranchiseeEmployeeDTO struct {
	employeesTypes.EmployeeDTO
	Franchisee franchiseesTypes.FranchiseeDTO
}
