package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

func ToSupplierResponse(supplier data.Supplier) SupplierResponse {
	return SupplierResponse{
		ID:           supplier.ID,
		Name:         supplier.Name,
		ContactEmail: supplier.ContactEmail,
		ContactPhone: supplier.ContactPhone,
		City:         supplier.City,
		Address:      supplier.Address,
		CreatedAt:    supplier.CreatedAt,
		UpdatedAt:    supplier.UpdatedAt,
	}
}

func ToSupplier(dto CreateSupplierDTO) data.Supplier {
	return data.Supplier{
		Name:         dto.Name,
		ContactEmail: dto.ContactEmail,
		ContactPhone: dto.ContactPhone,
		City:         dto.City,
		Address:      dto.Address,
	}
}

func ToSupplierMaterialResponse(material data.SupplierMaterial) SupplierMaterialResponse {
	var basePrice float64
	var effectiveDate time.Time

	// Extract the most recent price
	if len(material.SupplierPrices) > 0 {
		latestPrice := material.SupplierPrices[len(material.SupplierPrices)-1]
		basePrice = latestPrice.BasePrice
		effectiveDate = latestPrice.EffectiveDate
	}

	return SupplierMaterialResponse{
		StockMaterial: StockMaterialDTO{
			ID:          material.StockMaterial.ID,
			Name:        material.StockMaterial.Name,
			Description: material.StockMaterial.Description,
			Category:    material.StockMaterial.StockMaterialCategory.Name,
			SafetyStock: material.StockMaterial.SafetyStock,
			Barcode:     material.StockMaterial.Barcode,
			PackageMeasure: utils.PackageMeasure{
				UnitsPerPackage: material.StockMaterial.Package.Size,
				PackageUnit:     material.StockMaterial.Package.Unit.Name,
			},
		},
		BasePrice:     basePrice,
		EffectiveDate: effectiveDate,
	}
}

func ToSupplierMaterialResponses(materials []data.SupplierMaterial) []SupplierMaterialResponse {
	responses := make([]SupplierMaterialResponse, len(materials))
	for i, material := range materials {
		responses[i] = ToSupplierMaterialResponse(material)
	}
	return responses
}
