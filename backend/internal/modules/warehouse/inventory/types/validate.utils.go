package types

import (
	"errors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ValidateExpirationDate(newExpirationDate, oldExpirationDate time.Time) error {
	if newExpirationDate.Before(oldExpirationDate) {
		return errors.New("new expiration date cannot be earlier than the current expiration date")
	}
	return nil
}

func ValidatePackage(item InventoryItem) *data.Package {
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
	return &data.Package{
		StockMaterialID: item.StockMaterialID,
		PackageSize:     item.Package.PackageSize,
		PackageUnitID:   item.Package.PackageUnitID,
	}
}
