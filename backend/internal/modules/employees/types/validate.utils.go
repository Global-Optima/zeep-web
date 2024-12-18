package types

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func ValidateEmployee(input CreateEmployeeDTO) error {
	if input.Name == "" {
		return errors.New("employee name cannot be empty")
	}

	if !utils.IsValidEmail(input.Email) {
		return errors.New("invalid email format")
	}

	if err := utils.IsValidPassword(input.Password); err != nil {
		return fmt.Errorf("password validation failed: %v", err)
	}

	if data.IsValidEmployeeRole(input.Role) {
		return errors.New("invalid role specified")
	}

	if input.Type == data.StoreEmployeeType && input.StoreDetails == nil {
		return errors.New("store details are required for store employees")
	}
	if input.Type == data.WarehouseEmployeeType && input.WarehouseDetails == nil {
		return errors.New("warehouse details are required for warehouse employees")
	}

	return nil
}

func PrepareUpdateFields(input UpdateEmployeeDTO) (map[string]interface{}, error) {
	updateFields := map[string]interface{}{}

	if input.Name != nil {
		updateFields["name"] = *input.Name
	}
	if input.Phone != nil {
		updateFields["phone"] = *input.Phone
	}
	if input.Email != nil {
		if !utils.IsValidEmail(*input.Email) {
			return nil, errors.New("invalid email format")
		}
		updateFields["email"] = *input.Email
	}
	if input.Role != nil {
		if !data.IsValidEmployeeRole(*input.Role) {
			return nil, fmt.Errorf("invalid role: %v", *input.Role)
		}
		updateFields["role"] = *input.Role
	}

	if input.StoreDetails != nil {
		if input.StoreDetails.StoreID != nil {
			updateFields["store_employee.store_id"] = *input.StoreDetails.StoreID
		}
		if input.StoreDetails.IsFranchise != nil {
			updateFields["store_employee.is_franchise"] = *input.StoreDetails.IsFranchise
		}
	}

	if input.WarehouseDetails != nil {
		if input.WarehouseDetails.WarehouseID != nil {
			updateFields["warehouse_employee.warehouse_id"] = *input.WarehouseDetails.WarehouseID
		}
	}

	return updateFields, nil
}

func ParseEmployeeQuery(c *gin.Context) (*GetEmployeesFilter, error) {
	params := &GetEmployeesFilter{}

	err := c.ShouldBindQuery(params)
	if err != nil {
		return nil, err
	}

	params.Pagination = utils.ParsePagination(c)
	params.Sort, err = utils.ParseSortParamsForModel(c, &data.Employee{})
	if err != nil {
		return nil, err
	}

	return params, nil
}
