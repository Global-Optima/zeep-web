package data

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func UpdateOutOfStockInventory(tx *gorm.DB, storeID, ingredientID uint) error {
	if err := UpdateStoreProductsOutOfStock(tx, storeID, ingredientID); err != nil {
		return err
	}

	if err := UpdateStoreAdditivesOutOfStock(tx, storeID, ingredientID); err != nil {
		return err
	}

	return nil
}

func UpdateStoreAdditivesOutOfStock(tx *gorm.DB, storeID, ingredientID uint) error {
	logrus.Infof("storeAdditives, storeID=%d, ingredientID=%d", storeID, ingredientID)
	var storeAdditives []StoreAdditive
	err := tx.Model(&StoreAdditive{}).
		Joins("JOIN additive_ingredients ai ON ai.additive_id = store_additives.additive_id").
		Where("ai.ingredient_id = ?", ingredientID).
		Where("store_additives.store_id = ?", storeID).
		Find(&storeAdditives).Error
	if err != nil {
		return err
	}

	if len(storeAdditives) == 0 {
		return nil
	}

	for _, sa := range storeAdditives {
		outOfStock, err := isStoreAdditiveOutOfStock(tx, sa)
		if err != nil {
			return err
		}

		if err := tx.Model(&StoreAdditive{}).
			Where("id = ?", sa.ID).
			Update("is_out_of_stock", outOfStock).Error; err != nil {
			return err
		}
	}

	return nil
}

func isStoreAdditiveOutOfStock(tx *gorm.DB, sa StoreAdditive) (bool, error) {
	var additiveIngredients []AdditiveIngredient
	if err := tx.Where("additive_id = ?", sa.AdditiveID).
		Preload("Ingredient").
		Find(&additiveIngredients).Error; err != nil {
		return false, err
	}

	for _, ai := range additiveIngredients {
		var stock StoreStock
		err := tx.Where("store_id = ? AND ingredient_id = ?", sa.StoreID, ai.IngredientID).
			First(&stock).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) || stock.Quantity < stock.LowStockThreshold {
			return true, nil
		}
	}

	return false, nil
}

func UpdateStoreProductsOutOfStock(tx *gorm.DB, storeID, ingredientID uint) error {
	var spsList []StoreProductSize
	err := tx.Model(&StoreProductSize{}).
		Joins("JOIN store_products sp ON sp.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_ingredients psi ON psi.product_size_id = store_product_sizes.product_size_id").
		Where("psi.ingredient_id = ?", ingredientID).
		Where("sp.store_id = ?", storeID).
		Preload("StoreProduct").
		Find(&spsList).Error
	if err != nil {
		return err
	}

	if len(spsList) == 0 {
		return nil
	}

	for _, sps := range spsList {
		sp := sps.StoreProduct

		spsOutOfStock, err := isSPSOutOfStock(tx, sps)
		if err != nil {
			return err
		}

		if spsOutOfStock {
			if err := tx.Model(&StoreProduct{}).
				Where("id = ?", sp.ID).
				Update("is_out_of_stock", true).
				Error; err != nil {
				return err
			}
		} else {
			stillOut, err := anySPSOutOfStockForProduct(tx, sp.ID)
			if err != nil {
				return err
			}
			if !stillOut {
				if err := tx.Model(&StoreProduct{}).
					Where("id = ?", sp.ID).
					Update("is_out_of_stock", false).
					Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func isSPSOutOfStock(tx *gorm.DB, sps StoreProductSize) (bool, error) {
	storeID := sps.StoreProduct.StoreID

	var psIngredients []ProductSizeIngredient
	if err := tx.Where("product_size_id = ?", sps.ProductSizeID).
		Preload("Ingredient").
		Find(&psIngredients).Error; err != nil {
		return false, err
	}

	for _, psi := range psIngredients {
		var stock StoreStock
		err := tx.Where("store_id = ? AND ingredient_id = ?", storeID, psi.IngredientID).
			First(&stock).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) || stock.Quantity < stock.LowStockThreshold {
			return true, nil
		}
	}

	var defaultAdditives []ProductSizeAdditive
	if err := tx.
		Where("product_size_id = ? AND is_default = ?", sps.ProductSizeID, true).
		Find(&defaultAdditives).Error; err != nil {
		return false, err
	}

	for _, psa := range defaultAdditives {
		var storeAdditive StoreAdditive
		if err := tx.Where("store_id = ? AND additive_id = ?", storeID, psa.AdditiveID).
			First(&storeAdditive).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return true, nil
			}
			return false, err
		}

		if storeAdditive.IsOutOfStock {
			return true, nil
		}
	}

	return false, nil
}

func anySPSOutOfStockForProduct(tx *gorm.DB, storeProductID uint) (bool, error) {
	var spsList []StoreProductSize
	if err := tx.Where("store_product_id = ?", storeProductID).
		Preload("StoreProduct").
		Find(&spsList).Error; err != nil {
		return false, err
	}

	for _, sps := range spsList {
		spsOut, err := isSPSOutOfStock(tx, sps)
		if err != nil {
			return false, err
		}
		if spsOut {
			return true, nil
		}
	}
	return false, nil
}
