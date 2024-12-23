package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
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
