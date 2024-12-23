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

func ValidatePackage(stockMaterial data.StockMaterial) *data.StockMaterialPackage {
	if stockMaterial.Package.Size <= 0 || stockMaterial.Package.UnitID == 0 {
		return nil
	}
	return &data.StockMaterialPackage{
		StockMaterialID: stockMaterial.ID,
		Size:            stockMaterial.Package.Size,
		UnitID:          stockMaterial.Package.UnitID,
	}
}
