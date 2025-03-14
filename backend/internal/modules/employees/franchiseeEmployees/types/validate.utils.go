package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/employeeToken"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

type UpdateFranchiseeEmployeeModels struct {
	*employeesTypes.UpdateEmployeeModels
	FranchiseeEmployee *data.FranchiseeEmployee
}

func FranchiseeEmployeeUpdateFields(franchiseeEmployeeID uint, input *UpdateFranchiseeEmployeeDTO, role data.EmployeeRole, employeeTokenManager employeeToken.EmployeeTokenManager) (*UpdateFranchiseeEmployeeModels, error) {
	franchiseeEmployee := &data.FranchiseeEmployee{}
	if input.FranchiseeID != nil {
		franchiseeEmployee.FranchiseeID = *input.FranchiseeID
	}

	if input.Role != nil {
		if !data.IsAllowableRole(data.FranchiseeEmployeeType, *input.Role) {
			return nil, employeesTypes.ErrUnsupportedEmployeeType
		}
		if !data.CanManageRole(role, *input.Role) {
			return nil, fmt.Errorf("%s %w %s", role, employeesTypes.ErrNotAllowedToManageTheRole, *input.Role)
		}
		franchiseeEmployee.Role = *input.Role

		err := employeeTokenManager.DeleteTokenByFranchiseeEmployeeID(franchiseeEmployeeID)
		if err != nil {
			return nil, err
		}
	}

	updateEmployeeModels, err := employeesTypes.PrepareUpdateFields(&input.UpdateEmployeeDTO)
	if err != nil {
		return nil, err
	}

	return &UpdateFranchiseeEmployeeModels{
		UpdateEmployeeModels: updateEmployeeModels,
		FranchiseeEmployee:   franchiseeEmployee,
	}, nil
}
