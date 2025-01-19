package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func ValidateExpirationDays(addDays int) error {
	if addDays <= 0 {
		return errors.New("the number of days to extend must be greater than zero")
	}
	return nil
}

func ValidatePackage(stockMaterialID uint, pkg data.StockMaterialPackage) *data.StockMaterialPackage {
	if pkg.Size <= 0 || pkg.UnitID == 0 {
		return nil
	}
	return &data.StockMaterialPackage{
		StockMaterialID: stockMaterialID,
		Size:            pkg.Size,
		UnitID:          pkg.UnitID,
	}
}

func ParseStockFilterParamsWithPagination(c *gin.Context) (*GetWarehouseStockFilterQuery, error) {
	var params GetWarehouseStockFilterQuery

	if err := c.ShouldBindQuery(&params); err != nil {
		return nil, err
	}

	params.Pagination = utils.ParsePagination(c)

	return &params, nil
}
