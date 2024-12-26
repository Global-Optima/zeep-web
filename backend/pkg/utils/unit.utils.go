package utils

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ConvertPackagesToUnits(ingredient data.StockRequestIngredient, quantityInPackages float64) float64 {
	return ingredient.StockMaterial.Package.Size * quantityInPackages
}

type PackageMeasure struct {
	Quantity          float64 `json:"quantity"`          // Quantity of the packages
	UnitsPerPackage   string  `json:"unitsPerPackage"`   // Number of units contained in a package
	PackageUnit       string  `json:"packageUnit"`       // Name of the unit of the package
	TotalUnitsInStock float64 `json:"totalUnitsInStock"` // Total number of units
}

func ReturnPackageMeasure(ingredient data.StockRequestIngredient, quantityInPackages float64) PackageMeasure {
	pkgSize := fmt.Sprintf("%v %s", ingredient.StockMaterial.Package.Size, ingredient.StockMaterial.Package.Unit.Name)

	return PackageMeasure{
		Quantity:          quantityInPackages,
		UnitsPerPackage:   pkgSize,
		PackageUnit:       ingredient.StockMaterial.Package.Unit.Name,
		TotalUnitsInStock: ConvertPackagesToUnits(ingredient, quantityInPackages),
	}
}
