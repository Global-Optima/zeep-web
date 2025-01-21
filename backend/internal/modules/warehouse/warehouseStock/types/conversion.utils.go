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
		response[i] = DeliveryResponse{
			ID: delivery.ID,
			Material: WarehouseDeliveryStockMaterialDTO{
				StockMaterial: *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&delivery.StockMaterial),
				Package:       stockMaterialPackageTypes.ToStockMaterialPackageResponse(&delivery.Package),
				Quantity:      delivery.Quantity,
			},
			Supplier:       supplierTypes.ToSupplierResponse(delivery.Supplier),
			Warehouse:      *warehouseTypes.ToWarehouseResponse(delivery.Warehouse),
			Barcode:        delivery.Barcode,
			DeliveryDate:   delivery.DeliveryDate,
			ExpirationDate: delivery.ExpirationDate,
		}
	}
	return response
}

func ToDeliveryResponse(delivery data.SupplierWarehouseDelivery) DeliveryResponse {
	return DeliveryResponse{
		ID: delivery.ID,
		Material: WarehouseDeliveryStockMaterialDTO{
			StockMaterial: *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&delivery.StockMaterial),
			Package:       stockMaterialPackageTypes.ToStockMaterialPackageResponse(&delivery.Package),
			Quantity:      delivery.Quantity,
		},
		Supplier:       supplierTypes.ToSupplierResponse(delivery.Supplier),
		Warehouse:      *warehouseTypes.ToWarehouseResponse(delivery.Warehouse),
		Barcode:        delivery.Barcode,
		DeliveryDate:   delivery.DeliveryDate,
		ExpirationDate: delivery.ExpirationDate,
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

func ToStockMaterialDetails(stock data.AggregatedWarehouseStock, deliveries []data.SupplierWarehouseDelivery) WarehouseStockMaterialDetailsDTO {
	deliveriesDTO := make([]StockMaterialDeliveryDTO, len(deliveries))
	for i, delivery := range deliveries {
		deliveriesDTO[i] = StockMaterialDeliveryDTO{
			Supplier:       supplierTypes.ToSupplierResponse(delivery.Supplier),
			Quantity:       delivery.Quantity,
			DeliveryDate:   delivery.DeliveryDate,
			ExpirationDate: delivery.ExpirationDate,
		}
	}

	return WarehouseStockMaterialDetailsDTO{
		StockMaterial:          *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&stock.StockMaterial),
		Quantity:               stock.TotalQuantity,
		EarliestExpirationDate: stock.EarliestExpirationDate,
		Deliveries:             deliveriesDTO,
	}
}
