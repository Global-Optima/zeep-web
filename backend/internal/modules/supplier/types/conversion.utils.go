package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
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

	// Extract the most recent price
	if len(material.SupplierPrices) > 0 {
		latestPrice := material.SupplierPrices[len(material.SupplierPrices)-1]
		basePrice = latestPrice.BasePrice
	}

	return SupplierMaterialResponse{
		StockMaterial: SupplierStockMaterialDTO{
			*stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&material.StockMaterial),
		},
		BasePrice: basePrice,
	}
}

func ToSupplierMaterialResponses(materials []data.SupplierMaterial) []SupplierMaterialResponse {
	responses := make([]SupplierMaterialResponse, len(materials))
	for i, material := range materials {
		responses[i] = ToSupplierMaterialResponse(material)
	}
	return responses
}
