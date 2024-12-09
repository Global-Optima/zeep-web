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
		err := db.Where("sku_id = ?", item.SKU_ID).First(&mapping).Error
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
