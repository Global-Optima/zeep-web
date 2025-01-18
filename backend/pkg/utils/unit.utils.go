package utils

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type PackageMeasure struct {
	Quantity          float64 `json:"quantity"`          // Quantity of the packages
	UnitsPerPackage   float64 `json:"unitsPerPackage"`   // Number of units contained in a package
	PackageUnit       string  `json:"packageUnit"`       // Name of the unit of the package
	TotalUnitsInStock float64 `json:"totalUnitsInStock"` // Total number of units
}

func ConvertPackagesToUnits(stockMaterial data.StockMaterial, quantityInPackages float64) (float64, error) {
	if stockMaterial.Package == nil {
		return 0, fmt.Errorf("package not found for stock material ID: %d", stockMaterial.ID)
	}
	if stockMaterial.Package.Unit.ID != stockMaterial.Ingredient.Unit.ID {
		return stockMaterial.Package.Size * quantityInPackages * stockMaterial.Ingredient.Unit.ConversionFactor, nil
	}
	return stockMaterial.Package.Size * quantityInPackages, nil
}

func ReturnPackageMeasure(ingredient data.StockRequestIngredient, quantityInPackages float64) (PackageMeasure, error) {
	if ingredient.StockMaterial.Package == nil {
		return PackageMeasure{}, fmt.Errorf("package not found for stock material ID: %d", ingredient.StockMaterialID)
	}
	return PackageMeasure{
		Quantity:          quantityInPackages,
		UnitsPerPackage:   ingredient.StockMaterial.Package.Size,
		PackageUnit:       ingredient.StockMaterial.Package.Unit.Name,
		TotalUnitsInStock: ingredient.StockMaterial.Package.Size * quantityInPackages,
	}, nil
}

func ReturnPackageMeasureForStockMaterial(stockMaterial data.StockMaterial, quantityInPackages float64) (PackageMeasure, error) {
	if stockMaterial.Package == nil {
		return PackageMeasure{}, fmt.Errorf("package not found for stock material ID: %d", stockMaterial.ID)
	}
	return PackageMeasure{
		Quantity:          quantityInPackages,
		UnitsPerPackage:   stockMaterial.Package.Size,
		PackageUnit:       stockMaterial.Package.Unit.Name,
		TotalUnitsInStock: stockMaterial.Package.Size * quantityInPackages,
	}, nil
}
