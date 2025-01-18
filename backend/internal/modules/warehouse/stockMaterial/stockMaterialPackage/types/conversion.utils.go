package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
)

func ToStockMaterialPackageResponse(pkg *data.StockMaterialPackage) StockMaterialPackageResponse {
	return StockMaterialPackageResponse{
		ID:   pkg.ID,
		Size: pkg.Size,
		Unit: types.UnitResponse{
			ID:               pkg.Unit.ID,
			Name:             pkg.Unit.Name,
			ConversionFactor: pkg.Unit.ConversionFactor,
		},
		StockMaterial: StockMaterialDTO{
			ID:   pkg.StockMaterial.ID,
			Name: pkg.StockMaterial.Name,
		},
		CreatedAt: pkg.CreatedAt,
		UpdatedAt: pkg.UpdatedAt,
	}
}

func ToStockMaterialPackageResponses(pkgs []data.StockMaterialPackage) []StockMaterialPackageResponse {
	responses := make([]StockMaterialPackageResponse, len(pkgs))
	for i, pkg := range pkgs {
		responses[i] = ToStockMaterialPackageResponse(&pkg)
	}
	return responses
}
