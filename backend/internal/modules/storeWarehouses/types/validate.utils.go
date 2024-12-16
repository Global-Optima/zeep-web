package types

import (
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func ParseStockParamsWithPagination(c *gin.Context) (*GetStockFilterQuery, error) {
	var params GetStockFilterQuery

	if err := c.ShouldBindQuery(&params); err != nil {
		return nil, err
	}

	params.Pagination = utils.ParsePagination(c)

	return &params, nil
}
