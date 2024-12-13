package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

func ConvertInventoryItemsToStockRequest(items []InventoryItem, db *gorm.DB) ([]data.StockRequestIngredient, error) {
	converted := make([]data.StockRequestIngredient, len(items))

	for i, item := range items {
		var mapping data.IngredientsMapping
		err := db.Where("stock_material_id = ?", item.SKU_ID).First(&mapping).Error
		if err != nil {
			return nil, fmt.Errorf("failed to find ingredient mapping for SKU_ID %d: %w", item.SKU_ID, err)
		}

		converted[i] = data.StockRequestIngredient{
			IngredientID: mapping.IngredientID,
			Quantity:     item.Quantity,
		}
	}

	return converted, nil
}

func DeliveriesToDeliveryResponses(deliveries []data.Delivery) []DeliveryResponse {
	response := make([]DeliveryResponse, len(deliveries))
	for i, delivery := range deliveries {
		response[i] = DeliveryResponse{
			ID:             delivery.ID,
			SKU_ID:         delivery.StockMaterialID,
			Source:         delivery.SupplierID,
			Target:         delivery.WarehouseID,
			Barcode:        delivery.Barcode,
			Quantity:       delivery.Quantity,
			DeliveryDate:   delivery.DeliveryDate,
			ExpirationDate: delivery.ExpirationDate,
		}
	}
	return response
}

func StocksToInventoryItems(stocks []data.WarehouseStock) []InventoryItem {
	response := make([]InventoryItem, len(stocks))
	for i, stock := range stocks {
		response[i] = InventoryItem{
			SKU_ID:   stock.StockMaterialID,
			Quantity: stock.Quantity,
		}
	}
	return response
}

func ExpiringItemsToResponses(deliveries []data.Delivery) []UpcomingExpirationResponse {
	response := make([]UpcomingExpirationResponse, len(deliveries))
	for i, delivery := range deliveries {
		response[i] = UpcomingExpirationResponse{
			SKU_ID:         delivery.StockMaterialID,
			Name:           delivery.StockMaterial.Name,
			ExpirationDate: delivery.ExpirationDate,
			Quantity:       delivery.Quantity,
		}
	}
	return response
}
