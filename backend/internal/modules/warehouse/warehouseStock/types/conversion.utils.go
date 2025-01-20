package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	supplierTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/supplier/types"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

func ConvertExistingItemsToStockRequest(items []ExistingWarehouseStockMaterial, db *gorm.DB) ([]data.StockRequestIngredient, error) {
	converted := make([]data.StockRequestIngredient, len(items))

	for i, item := range items {
		var stockMaterial data.StockMaterial
		err := db.Preload("Ingredient").First(&stockMaterial, "id = ?", item.StockMaterialID).Error
		if err != nil {
			return nil, fmt.Errorf("failed to find StockMaterial for ID %d: %w", item.StockMaterialID, err)
		}

		converted[i] = data.StockRequestIngredient{
			IngredientID: stockMaterial.IngredientID,
			Quantity:     item.Quantity,
		}
	}

	return converted, nil
}

func DeliveriesToDeliveryResponses(deliveries []data.SupplierWarehouseDelivery) []DeliveryResponse {
	response := make([]DeliveryResponse, len(deliveries))
	for i, delivery := range deliveries {
		response[i] = DeliveryResponse{
			ID:             delivery.ID,
			StockMaterial:  *stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&delivery.StockMaterial),
			Supplier:       supplierTypes.ToSupplierResponse(delivery.Supplier),
			Warehouse:      *warehouseTypes.ToWarehouseResponse(delivery.Warehouse),
			Barcode:        delivery.Barcode,
			Quantity:       delivery.Quantity,
			DeliveryDate:   delivery.DeliveryDate,
			ExpirationDate: delivery.ExpirationDate,
		}
	}
	return response
}

func ToWarehouseStockResponse(stock data.AggregatedWarehouseStock, pkgMeasure utils.PackageMeasureWithQuantity) WarehouseStockResponse {
	return WarehouseStockResponse{
		StockMaterial: StockMaterialResponse{
			*stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&stock.StockMaterial),
			pkgMeasure,
		},
		EarliestExpirationDate: stock.EarliestExpirationDate,
	}
}

func ToStockMaterialDetails(stock data.AggregatedWarehouseStock, pkgMeasure utils.PackageMeasureWithQuantity, deliveries []data.SupplierWarehouseDelivery) WarehouseStockMaterialDetailsDTO {
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
		PackageMeasure:         pkgMeasure,
		EarliestExpirationDate: stock.EarliestExpirationDate,
		Deliveries:             deliveriesDTO,
	}
}
