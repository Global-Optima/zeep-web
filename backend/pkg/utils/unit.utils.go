package utils

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type PackageMeasure struct {
	UnitsPerPackage float64 `json:"unitsPerPackage"` // Number of units contained in a package
	PackageUnit     string  `json:"packageUnit"`     // Name of the unit of the package
}

// FindPackageByUnit finds the correct package for a given stock material by unit ID.
func FindPackageByUnit(stockMaterial data.StockMaterial, unitID uint) (*data.StockMaterialPackage, error) {
	for _, pkg := range stockMaterial.Packages {
		if pkg.UnitID == unitID {
			return &pkg, nil
		}
	}
	return nil, fmt.Errorf("no package found for unit ID: %d in stock material ID: %d", unitID, stockMaterial.ID)
}

// ConvertPackagesToUnits converts a quantity of packages to units, based on the correct package.
func ConvertPackagesToUnits(stockMaterial data.StockMaterial, unitID uint, quantityInPackages float64) (float64, error) {
	packageDetails, err := FindPackageByUnit(stockMaterial, unitID)
	if err != nil {
		return 0, err
	}

	if packageDetails.Unit.ID != stockMaterial.Ingredient.Unit.ID {
		return packageDetails.Size * quantityInPackages * stockMaterial.Ingredient.Unit.ConversionFactor, nil
	}

	return packageDetails.Size * quantityInPackages, nil
}

// ReturnPackageMeasures returns details of all available packages for a stock material.
func ReturnPackageMeasures(stockMaterial data.StockMaterial) ([]PackageMeasure, error) {
	if len(stockMaterial.Packages) == 0 {
		return nil, errors.New("no packages found for stock material")
	}

	var packageMeasures []PackageMeasure
	for _, pkg := range stockMaterial.Packages {
		packageMeasures = append(packageMeasures, PackageMeasure{
			UnitsPerPackage: pkg.Size,
			PackageUnit:     pkg.Unit.Name,
		})
	}

	return packageMeasures, nil
}

// ReturnAllPackageMeasures returns details for all available packages in a stock material.
func ReturnAllPackageMeasures(stockMaterial data.StockMaterial) ([]PackageMeasure, error) {
	if len(stockMaterial.Packages) == 0 {
		return nil, errors.New("no packages available for stock material")
	}

	var packageMeasures []PackageMeasure
	for _, pkg := range stockMaterial.Packages {
		packageMeasures = append(packageMeasures, PackageMeasure{
			UnitsPerPackage: pkg.Size,
			PackageUnit:     pkg.Unit.Name,
		})
	}

	return packageMeasures, nil
}

// ReturnTotalUnitsAcrossAllPackages calculates the total units across all packages for the given quantities.
func ReturnTotalUnitsAcrossAllPackages(stockMaterial data.StockMaterial, quantitiesInPackages map[uint]float64) (float64, error) {
	if len(stockMaterial.Packages) == 0 {
		return 0, errors.New("no packages available for stock material")
	}

	var totalUnits float64
	for _, pkg := range stockMaterial.Packages {
		quantityInPackages, ok := quantitiesInPackages[pkg.UnitID]
		if !ok {
			continue
		}

		totalUnits += pkg.Size * quantityInPackages
	}

	return totalUnits, nil
}
