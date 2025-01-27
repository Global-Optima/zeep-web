package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
)

func ToStockRequestResponse(request *data.StockRequest) StockRequestResponse {
	items := make([]StockRequestMaterial, len(request.Ingredients))
	for i, ingredient := range request.Ingredients {
		items[i] = StockRequestMaterial{
			StockMaterial: *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&ingredient.StockMaterial),
			Quantity:      ingredient.Quantity,
		}
	}

	facilityAddress := &storeTypes.FacilityAddressDTO{
		ID:      request.Store.FacilityAddressID,
		Address: request.Store.FacilityAddress.Address,
	}

	if request.Store.FacilityAddress.Longitude == nil || request.Store.FacilityAddress.Latitude == nil {
		facilityAddress.Longitude = 0
		facilityAddress.Latitude = 0
	} else {
		facilityAddress.Longitude = *request.Store.FacilityAddress.Longitude
		facilityAddress.Latitude = *request.Store.FacilityAddress.Latitude
	}

	return StockRequestResponse{
		RequestID:        request.ID,
		StoreComment:     request.StoreComment,
		WarehouseComment: request.WarehouseComment,
		Store: storeTypes.StoreDTO{
			ID:              request.StoreID,
			Name:            request.Store.Name,
			IsFranchise:     request.Store.IsFranchise,
			ContactPhone:    request.Store.ContactPhone,
			ContactEmail:    request.Store.ContactEmail,
			FacilityAddress: facilityAddress,
			StoreHours:      request.Store.StoreHours,
		},
		Warehouse:      *warehouseTypes.ToWarehouseResponse(request.Warehouse),
		Status:         request.Status,
		StockMaterials: items,
		CreatedAt:      request.CreatedAt,
		UpdatedAt:      request.UpdatedAt,
	}
}
