package orders

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

// TODO remove?
func getOrderIngredients(tx *gorm.DB, orderID uint) ([]uint, error) {
	ingredientIDsFromProductSizes, err := getOrderProductSizeIngredients(tx, orderID)
	if err != nil {
		return nil, err
	}

	ingredientIDsFromAdditives, err := getOrderAdditiveIngredients(tx, orderID)
	if err != nil {
		return nil, err
	}

	return utils.UnionSlices(ingredientIDsFromProductSizes, ingredientIDsFromAdditives), nil
}

func getOrderProductSizeIngredients(tx *gorm.DB, orderID uint) ([]uint, error) {
	var ingredientIDsFromProductSizes []uint

	err := tx.Model(&data.ProductSizeIngredient{}).
		Distinct("product_size_ingredients.ingredient_id").
		Joins("JOIN store_product_sizes ON store_product_sizes.product_size_id = product_size_ingredients.product_size_id").
		Joins("JOIN suborders ON suborders.store_product_size_id = store_product_sizes.id").
		Where("suborders.order_id = ?", orderID).
		Pluck("product_size_ingredients.ingredient_id", &ingredientIDsFromProductSizes).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredients from product sizes: %w", err)
	}

	return ingredientIDsFromProductSizes, nil
}

func getOrderAdditiveIngredients(tx *gorm.DB, orderID uint) ([]uint, error) {
	var ingredientIDsFromAdditives []uint

	err := tx.Model(&data.AdditiveIngredient{}).
		Distinct("additive_ingredients.ingredient_id").
		Joins("JOIN store_additives ON store_additives.additive_id = additive_ingredients.additive_id").
		Joins("JOIN suborder_additives ON suborder_additives.store_additive_id = store_additives.id").
		Joins("JOIN suborders ON suborders.id = suborder_additives.suborder_id").
		Where("suborders.order_id = ?", orderID).
		Pluck("additive_ingredients.ingredient_id", &ingredientIDsFromAdditives).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredients from additives: %w", err)
	}

	return ingredientIDsFromAdditives, nil
}
