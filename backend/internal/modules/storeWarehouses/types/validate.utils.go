package types

import (
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ParseStockParamsWithPagination(c *gin.Context) (*GetStockQuery, error) {
	params := &GetStockQuery{}

	query := c.Request.URL.Query()

	if query.Get("search") != "" {
		search := query.Get("search")
		params.Search = &search
	}

	if query.Get("lowStockOnly") != "" {
		lowStockOnly, err := strconv.ParseBool(query.Get("lowStockOnly"))
		if err != nil {
			return nil, err
		}
		params.LowStockOnly = &lowStockOnly
	}

	params.Pagination = utils.ParsePagination(c)

	return params, nil
}
