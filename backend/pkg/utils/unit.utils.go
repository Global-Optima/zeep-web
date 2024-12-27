package utils

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type PackageMeasure struct {
	Quantity          float64 `json:"quantity"`          // Quantity of the packages
	UnitsPerPackage   float64 `json:"unitsPerPackage"`   // Number of units contained in a package
	PackageUnit       string  `json:"packageUnit"`       // Name of the unit of the package
	TotalUnitsInStock float64 `json:"totalUnitsInStock"` // Total number of units
}

func ConvertPackagesToUnits(ingredient data.StockRequestIngredient, quantityInPackages float64) float64 {
	return ingredient.StockMaterial.Package.Size * quantityInPackages
}

func ReturnPackageMeasure(ingredient data.StockRequestIngredient, quantityInPackages float64) PackageMeasure {
	return PackageMeasure{
		Quantity:          quantityInPackages,
		UnitsPerPackage:   ingredient.StockMaterial.Package.Size,
		PackageUnit:       ingredient.StockMaterial.Package.Unit.Name,
		TotalUnitsInStock: ingredient.StockMaterial.Package.Size * quantityInPackages,
	}
}

func ReturnPackageMeasureForStockMaterial(stockMaterial data.StockMaterial, quantityInPackages float64) PackageMeasure {
	return PackageMeasure{
		Quantity:          quantityInPackages,
		UnitsPerPackage:   stockMaterial.Package.Size,
		PackageUnit:       stockMaterial.Package.Unit.Name,
		TotalUnitsInStock: stockMaterial.Package.Size * quantityInPackages,
	}
}
