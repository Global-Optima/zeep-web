package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	supplierTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/supplier/types"
	stockMaterialPackageTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialPackage/types"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
)

func DeliveriesToDeliveryResponses(deliveries []data.SupplierWarehouseDelivery) []DeliveryResponse {
	response := make([]DeliveryResponse, len(deliveries))
	for i, delivery := range deliveries {
		materials := make([]WarehouseDeliveryStockMaterialDTO, len(delivery.Materials))

		for i, material := range delivery.Materials {
			materials[i] = WarehouseDeliveryStockMaterialDTO{
				StockMaterial:  *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&material.StockMaterial),
				Package:        stockMaterialPackageTypes.ToStockMaterialPackageResponse(&material.Package),
				Quantity:       material.Quantity,
				Barcode:        material.Barcode,
				ExpirationDate: material.ExpirationDate,
			}
		}

		response[i] = DeliveryResponse{
			ID:           delivery.ID,
			Supplier:     supplierTypes.ToSupplierResponse(delivery.Supplier),
			Warehouse:    *warehouseTypes.ToWarehouseResponse(delivery.Warehouse),
			Materials:    materials,
			DeliveryDate: delivery.DeliveryDate,
		}
	}
	return response
}

func ToDeliveryResponse(delivery data.SupplierWarehouseDelivery) DeliveryResponse {
	materials := make([]WarehouseDeliveryStockMaterialDTO, len(delivery.Materials))

	for i, material := range delivery.Materials {
		materials[i] = WarehouseDeliveryStockMaterialDTO{
			StockMaterial:  *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&material.StockMaterial),
			Package:        stockMaterialPackageTypes.ToStockMaterialPackageResponse(&material.Package),
			Quantity:       material.Quantity,
			Barcode:        material.Barcode,
			ExpirationDate: material.ExpirationDate,
		}
	}

	return DeliveryResponse{
		ID:           delivery.ID,
		Supplier:     supplierTypes.ToSupplierResponse(delivery.Supplier),
		Warehouse:    *warehouseTypes.ToWarehouseResponse(delivery.Warehouse),
		Materials:    materials,
		DeliveryDate: delivery.DeliveryDate,
	}
}

func ToWarehouseStockResponse(stock data.AggregatedWarehouseStock) WarehouseStockResponse {
	return WarehouseStockResponse{
		StockMaterial: StockMaterialResponse{
			*stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&stock.StockMaterial),
		},
		Quantity:               stock.TotalQuantity,
		EarliestExpirationDate: stock.EarliestExpirationDate,
	}
}
