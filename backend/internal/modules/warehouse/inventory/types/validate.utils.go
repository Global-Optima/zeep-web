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

func ValidatePackage(item InventoryItem) *data.StockMaterialPackage {
	if item.Package == nil {
		return nil
	}
	if item.StockMaterialID == 0 {
		return nil
	}
	if item.Package.PackageSize == 0 {
		return nil
	}
	if item.Package.PackageUnitID == 0 {
		return nil
	}
	return &data.StockMaterialPackage{
		StockMaterialID: item.StockMaterialID,
		PackageSize:     item.Package.PackageSize,
		PackageUnitID:   item.Package.PackageUnitID,
	}
}
