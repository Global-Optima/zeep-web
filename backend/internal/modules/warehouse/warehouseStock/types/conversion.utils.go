package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	supplierTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/supplier/types"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
)

func DeliveriesToDeliveryResponses(deliveries []data.SupplierWarehouseDelivery) []WarehouseDeliveryDTO {
	response := make([]WarehouseDeliveryDTO, len(deliveries))
	for i, delivery := range deliveries {
		materials := make([]WarehouseDeliveryStockMaterialDTO, len(delivery.Materials))

		for i, material := range delivery.Materials {
			materials[i] = WarehouseDeliveryStockMaterialDTO{
				StockMaterial:  *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&material.StockMaterial),
				Quantity:       material.Quantity,
				Barcode:        material.Barcode,
				ExpirationDate: material.ExpirationDate,
			}
		}

		response[i] = WarehouseDeliveryDTO{
			ID:           delivery.ID,
			Supplier:     supplierTypes.ToSupplierResponse(delivery.Supplier),
			Warehouse:    *warehouseTypes.ToWarehouseDTO(delivery.Warehouse),
			Materials:    materials,
			DeliveryDate: delivery.DeliveryDate,
		}
	}
	return response
}

func ToDeliveryResponse(delivery data.SupplierWarehouseDelivery) WarehouseDeliveryDTO {
	materials := make([]WarehouseDeliveryStockMaterialDTO, len(delivery.Materials))

	for i, material := range delivery.Materials {
		materials[i] = WarehouseDeliveryStockMaterialDTO{
			StockMaterial:  *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&material.StockMaterial),
			Quantity:       material.Quantity,
			Barcode:        material.Barcode,
			ExpirationDate: material.ExpirationDate,
		}
	}

	return WarehouseDeliveryDTO{
		ID:           delivery.ID,
		Supplier:     supplierTypes.ToSupplierResponse(delivery.Supplier),
		Warehouse:    *warehouseTypes.ToWarehouseDTO(delivery.Warehouse),
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
