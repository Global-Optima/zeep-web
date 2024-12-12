package types

import (
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

func PrepareUpdateFields(input UpdateStockDTO) (map[string]interface{}, error) {
	updateFields := map[string]interface{}{}

	if input.CurrentStock != nil {
		updateFields["quantity"] = *input.CurrentStock
	}
	if input.LowStockThreshold != nil {
		updateFields["low_stock_threshold"] = *input.LowStockThreshold
	}

	return updateFields, nil
}

func ParseStockParams(query url.Values) (*GetStockQuery, error) {
	params := &GetStockQuery{}

	if query.Get("searchTerm") != "" {
		searchTerm := query.Get("searchTerm")
		params.SearchTerm = &searchTerm
	}

	if query.Get("lowStockOnly") != "" {
		lowStockOnly, err := strconv.ParseBool(query.Get("lowStockOnly"))
		if err != nil {
			return nil, err
		}
		params.LowStockOnly = &lowStockOnly
	}

	return params, nil
}

func ParseAddStockParams(query url.Values) (*AddStockDTO, error) {
	params := &AddStockDTO{}

	if query.Get("ingredientId") != "" {
		ingredientID, err := strconv.ParseUint(query.Get("ingredientId"), 10, 64)
		if err != nil {
			return nil, errors.New("invalid ingredient ID")
		}
		params.IngredientID = uint(ingredientID)
	} else {
		return nil, errors.New("ingredientId is required")
	}

	if query.Get("currentStock") != "" {
		currentStock, err := strconv.ParseFloat(query.Get("currentStock"), 64)
		if err != nil {
			return nil, errors.New("invalid current stock value")
		}
		params.CurrentStock = currentStock
	}

	if query.Get("lowStockAlert") != "" {
		lowStockThreshold, err := strconv.ParseFloat(query.Get("lowStockAlert"), 64)
		if err != nil {
			return nil, errors.New("invalid low stock alert value")
		}
		params.LowStockThreshold = lowStockThreshold
	} else {
		return nil, errors.New("lowStockAlert is required")
	}

	return params, nil
}
