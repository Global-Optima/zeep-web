package types

import (
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

func PrepareUpdateFields(input UpdateStoreWarehouseIngredientDTO) (map[string]interface{}, error) {
	updateFields := map[string]interface{}{}

	if input.CurrentStock != nil {
		updateFields["store_warehouse_stock.quantity"] = *input.CurrentStock
	}
	if input.LowStockThreshold != nil {
		updateFields["store_warehouse_stock.low_stock_threshold"] = *input.LowStockThreshold
	}

	return updateFields, nil
}

func ParseStoreWarehouseIngredientParams(query url.Values) (*GetStoreWarehouseStockQuery, error) {
	params := &GetStoreWarehouseStockQuery{}

	if query.Get("storeId") != "" {
		storeID, err := strconv.ParseUint(query.Get("storeId"), 10, 64)
		if err != nil {
			return nil, errors.New("invalid store ID")
		}
		storeIDUint := uint(storeID)
		params.StoreID = storeIDUint
	} else {
		return nil, errors.New("storeId is required")
	}

	// Parse `searchTerm`
	if query.Get("searchTerm") != "" {
		searchTerm := query.Get("role")
		params.SearchTerm = &searchTerm
	}

	if query.Get("lowStockOnly") != "" {
		lowStockOnly, err := strconv.ParseBool(query.Get("lowStockOnly"))
		if err != nil {
			return nil, err
		}
		params.LowStockOnly = &lowStockOnly
	}

	if query.Get("limit") != "" {
		limit, err := strconv.Atoi(query.Get("limit"))
		if err != nil || limit < 0 {
			return nil, errors.New("invalid limit")
		}
		params.Limit = limit
	} else {
		params.Limit = 10 // Default value
	}

	if query.Get("offset") != "" {
		offset, err := strconv.Atoi(query.Get("offset"))
		if err != nil || offset < 0 {
			return nil, errors.New("invalid offset")
		}
		params.Offset = offset
	} else {
		params.Offset = 0 // Default value
	}

	return params, nil
}
