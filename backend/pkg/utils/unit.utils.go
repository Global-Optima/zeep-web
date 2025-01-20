package utils

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type PackageMeasure struct {
	UnitsPerPackage float64 `json:"unitsPerPackage"` // Number of units contained in a package
	PackageUnit     string  `json:"packageUnit"`     // Name of the unit of the package
}

type PackageMeasureWithQuantity struct {
	PackageMeasure
	Quantity          float64 `json:"quantity"`          // Quantity of the packages
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

func ReturnPackageMeasureWithQuantity(ingredient data.StockRequestIngredient, quantityInPackages float64) (PackageMeasureWithQuantity, error) {
	if ingredient.StockMaterial.Package == nil {
		return PackageMeasureWithQuantity{}, fmt.Errorf("package not found for stock material ID: %d", ingredient.StockMaterialID)
	}
	return PackageMeasureWithQuantity{
		PackageMeasure: PackageMeasure{
			UnitsPerPackage: ingredient.StockMaterial.Package.Size,
			PackageUnit:     ingredient.StockMaterial.Package.Unit.Name,
		},
		TotalUnitsInStock: ingredient.StockMaterial.Package.Size * quantityInPackages,
		Quantity:          quantityInPackages,
	}, nil
}

func ReturnPackageMeasureForStockMaterialWithQuantity(stockMaterial data.StockMaterial, quantityInPackages float64) (PackageMeasureWithQuantity, error) {
	if stockMaterial.Package == nil {
		return PackageMeasureWithQuantity{}, fmt.Errorf("package not found for stock material ID: %d", stockMaterial.ID)
	}
	return PackageMeasureWithQuantity{
		PackageMeasure: PackageMeasure{
			UnitsPerPackage: stockMaterial.Package.Size,
			PackageUnit:     stockMaterial.Package.Unit.Name,
		},
		Quantity:          quantityInPackages,
		TotalUnitsInStock: stockMaterial.Package.Size * quantityInPackages,
	}, nil
}

func ReturnPackageMeasures(stockMaterial data.StockMaterial) (PackageMeasure, error) {
	if stockMaterial.Package == nil {
		return PackageMeasure{}, fmt.Errorf("package not found for stock material ID: %d", stockMaterial.ID)
	}
	return PackageMeasure{
		UnitsPerPackage: stockMaterial.Package.Size,
		PackageUnit:     stockMaterial.Package.Unit.Name,
	}, nil
}
