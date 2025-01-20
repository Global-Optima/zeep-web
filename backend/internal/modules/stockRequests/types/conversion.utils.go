package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

func ToStockRequestResponse(request *data.StockRequest) StockRequestResponse {
	items := make([]StockRequestMaterial, len(request.Ingredients))
	for i, ingredient := range request.Ingredients {
		var packageMeasures utils.PackageMeasureWithQuantity

		packageMeasures, err := utils.ReturnPackageMeasureWithQuantity(ingredient, ingredient.Quantity)
		if err != nil {
			panic(fmt.Sprintf("Critical error: %v", err))
		}

		items[i] = StockRequestMaterial{
			StockMaterial:              *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&ingredient.StockMaterial),
			PackageMeasureWithQuantity: packageMeasures,
		}
	}

	return StockRequestResponse{
		RequestID:        request.ID,
		StoreComment:     request.StoreComment,
		WarehouseComment: request.WarehouseComment,
		Store: storeTypes.StoreDTO{
			ID:           request.StoreID,
			Name:         request.Store.Name,
			IsFranchise:  request.Store.IsFranchise,
			ContactPhone: request.Store.ContactPhone,
			ContactEmail: request.Store.ContactEmail,
			FacilityAddress: &storeTypes.FacilityAddressDTO{
				ID:        request.Store.FacilityAddressID,
				Address:   request.Store.FacilityAddress.Address,
				Longitude: *request.Store.FacilityAddress.Longitude,
				Latitude:  *request.Store.FacilityAddress.Latitude,
			},
			StoreHours: request.Store.StoreHours,
		},
		Warehouse:      *warehouseTypes.ToWarehouseResponse(request.Warehouse),
		Status:         request.Status,
		StockMaterials: items,
		CreatedAt:      request.CreatedAt,
		UpdatedAt:      request.UpdatedAt,
	}
}
