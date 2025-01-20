package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
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
			ID:              delivery.ID,
			StockMaterialID: delivery.StockMaterialID,
			SupplierID:      delivery.SupplierID,
			WarehouseID:     delivery.WarehouseID,
			Barcode:         delivery.Barcode,
			Quantity:        delivery.Quantity,
			DeliveryDate:    delivery.DeliveryDate,
			ExpirationDate:  delivery.ExpirationDate,
		}
	}
	return response
}

func ToWarehouseStockResponse(stock data.AggregatedWarehouseStock, pkgMeasure utils.PackageMeasure) WarehouseStockResponse {
	return WarehouseStockResponse{
		StockMaterial: StockMaterialResponse{
			*stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(&stock.StockMaterial),
			pkgMeasure,
		},
		EarliestExpirationDate: stock.EarliestExpirationDate,
	}
}

func ToStockMaterialDetails(stock data.AggregatedWarehouseStock, pkgMeasure utils.PackageMeasure, deliveries []data.SupplierWarehouseDelivery) WarehouseStockMaterialDetailsDTO {
	deliveriesDTO := make([]StockMaterialDeliveryDTO, len(deliveries))
	for i, delivery := range deliveries {
		deliveriesDTO[i] = StockMaterialDeliveryDTO{
			Supplier:       delivery.Supplier.Name,
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
