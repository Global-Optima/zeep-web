package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func ParseStockParamsWithPagination(c *gin.Context) (*GetStockFilterQuery, error) {
	var params GetStockFilterQuery
	var err error

	if err = c.ShouldBindQuery(&params); err != nil {
		return nil, err
	}

	params.Pagination = utils.ParsePagination(c)

	params.Sort, err = utils.ParseSortParamsForModel(c, &data.StoreStock{})
	if err != nil {
		return nil, err
	}

	return &params, nil
}
