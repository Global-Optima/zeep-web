package utils

func ConvertPackagesToUnits(packageSize, quantityInPackages, conversionFactor float64) float64 {
	return packageSize * quantityInPackages * conversionFactor
}
