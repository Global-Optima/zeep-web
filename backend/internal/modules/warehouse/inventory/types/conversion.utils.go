package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

func ConvertExistingItemsToStockRequest(items []ExistingInventoryItem, db *gorm.DB) ([]data.StockRequestIngredient, error) {
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

func StocksToInventoryItems(stocks []data.WarehouseStock) *InventoryLevelsResponse {
	levels := make([]InventoryLevel, len(stocks))
	for i, stock := range stocks {
		levels[i] = InventoryLevel{
			StockMaterialID: stock.StockMaterialID,
			Name:            stock.StockMaterial.Name,
			Quantity:        stock.Quantity,
		}
	}

	return &InventoryLevelsResponse{
		Levels: levels,
	}
}

func ExpiringItemsToResponses(deliveries []data.SupplierWarehouseDelivery) []UpcomingExpirationResponse {
	response := make([]UpcomingExpirationResponse, len(deliveries))
	for i, delivery := range deliveries {
		response[i] = UpcomingExpirationResponse{
			DeliveryID:      delivery.ID,
			StockMaterialID: delivery.StockMaterialID,
			Name:            delivery.StockMaterial.Name,
			ExpirationDate:  delivery.ExpirationDate,
			Quantity:        delivery.Quantity,
		}
	}
	return response
}
