package utils

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func ConvertPackagesToUnits(ingredient data.StockRequestIngredient, quantityInPackages float64) float64 {
	if ingredient.StockMaterial.Package.Unit.ID != ingredient.Ingredient.Unit.ID {
		return ingredient.StockMaterial.Package.Size * quantityInPackages * ingredient.StockMaterial.Package.Unit.ConversionFactor
	}

	return ingredient.StockMaterial.Package.Size * quantityInPackages
}

type PackageMeasure struct {
	Quantity        float64 `json:"quantity"`        // Quantity of the packages
	UnitsPerPackage float64 `json:"unitsPerPackage"` // Number of units contained in a package
	Unit            string  `json:"unit"`            // Name of the unit
	TotalUnits      float64 `json:"totalUnits"`      // Total number of units

}

func ReturnPackageMeasure(ingredient data.StockRequestIngredient, quantityInPackages float64) PackageMeasure {
	return PackageMeasure{
		Quantity:        quantityInPackages,
		UnitsPerPackage: ingredient.StockMaterial.Package.Size,
		Unit:            ingredient.StockMaterial.Package.Unit.Name,
		TotalUnits:      ConvertPackagesToUnits(ingredient, quantityInPackages),
	}
}
