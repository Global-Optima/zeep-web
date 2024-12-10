package types

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
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

func ParseEmployeeQueryParams(query url.Values) (*GetEmployeesQuery, error) {
	params := &GetEmployeesQuery{}

	// Parse `type`
	if query.Get("type") != "" {
		employeeType := query.Get("type")
		params.Type = &employeeType
	}

	// Parse `role`
	if query.Get("role") != "" {
		role := query.Get("role")
		params.Role = &role
	}

	// Parse `store_id`
	if query.Get("storeId") != "" {
		storeID, err := strconv.ParseUint(query.Get("storeId"), 10, 64)
		if err != nil {
			return nil, errors.New("invalid store ID")
		}
		storeIDUint := uint(storeID)
		params.StoreID = &storeIDUint
	}

	// Parse `warehouse_id`
	if query.Get("warehouseId") != "" {
		warehouseID, err := strconv.ParseUint(query.Get("warehouseId"), 10, 64)
		if err != nil {
			return nil, errors.New("invalid warehouse ID")
		}
		warehouseIDUint := uint(warehouseID)
		params.WarehouseID = &warehouseIDUint
	}

	// Parse `limit`
	if query.Get("limit") != "" {
		limit, err := strconv.Atoi(query.Get("limit"))
		if err != nil || limit < 0 {
			return nil, errors.New("invalid limit value")
		}
		params.Limit = limit
	} else {
		params.Limit = 10 // Default value
	}

	// Parse `offset`
	if query.Get("offset") != "" {
		offset, err := strconv.Atoi(query.Get("offset"))
		if err != nil || offset < 0 {
			return nil, errors.New("invalid offset value")
		}
		params.Offset = offset
	} else {
		params.Offset = 0 // Default value
	}

	return params, nil
}
