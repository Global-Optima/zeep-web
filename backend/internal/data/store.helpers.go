package data

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RecalculateOutOfStock(tx *gorm.DB, storeID uint, ingredientIDs []uint, productSizeIDs []uint) error {
	if storeID == 0 || (len(ingredientIDs) == 0 && len(productSizeIDs) == 0) {
		return nil
	}

	var storeProductIDsFromPS,
		storeProductIDsFromIngredients,
		storeAdditiveIDs []uint
	var err error

	if len(productSizeIDs) > 0 {
		storeProductIDsFromPS, err = getStoreProductIDsByProductSizes(tx, storeID, productSizeIDs)
		if err != nil {
			return err
		}
	}

	if len(ingredientIDs) > 0 {
		storeProductIDsFromIngredients, err = getStoreProductIDsByIngredients(tx, storeID, ingredientIDs)
		if err != nil {
			return err
		}
	}

	if len(storeProductIDsFromPS) > 0 || len(storeProductIDsFromIngredients) > 0 {
		mergedStoreProductIDs := unionUintSlices(storeProductIDsFromPS, storeProductIDsFromIngredients)

		if err := recalculateStoreProducts(tx, mergedStoreProductIDs); err != nil {
			return err
		}
	}

	if len(ingredientIDs) > 0 {
		storeAdditiveIDs, err = getStoreAdditiveIDsByIngredients(tx, storeID, ingredientIDs)
		if err != nil {
			return err
		}
		if err := recalculateStoreAdditives(tx, storeAdditiveIDs); err != nil {
			return err
		}
	}

	logrus.Info("=============================================")

	return nil
}

func recalculateStoreProducts(tx *gorm.DB, storeProductIDs []uint) error {
	if len(storeProductIDs) == 0 {
		return nil
	}

	outOfStockIDs, err := getOutOfStockStoreProductIDs(tx, storeProductIDs)
	if err != nil {
		return err
	}
	inStockIDs := diffUintSlice(storeProductIDs, outOfStockIDs)
	logrus.Infof("outOfStockProducts: %v", outOfStockIDs)
	logrus.Infof("inStockProducts: %v", inStockIDs)

	if err := updateStoreProductStockFlags(tx, outOfStockIDs, true); err != nil {
		return err
	}
	if err := updateStoreProductStockFlags(tx, inStockIDs, false); err != nil {
		return err
	}

	return nil
}

func recalculateStoreAdditives(tx *gorm.DB, additiveIDs []uint) error {
	if len(additiveIDs) == 0 {
		return nil
	}

	outOfStockIDs, err := getOutOfStockStoreAdditiveIDs(tx, additiveIDs)
	if err != nil {
		return err
	}
	inStockIDs := diffUintSlice(additiveIDs, outOfStockIDs)
	logrus.Infof("outOfStockAdditives: %v", outOfStockIDs)
	logrus.Infof("inStockAdditives: %v", inStockIDs)

	if err := updateStoreAdditiveStockFlags(tx, outOfStockIDs, true); err != nil {
		return err
	}
	if err := updateStoreAdditiveStockFlags(tx, inStockIDs, false); err != nil {
		return err
	}

	return nil
}

func updateStoreProductStockFlags(tx *gorm.DB, ids []uint, isOutOfStock bool) error {
	if len(ids) == 0 {
		return nil
	}
	return tx.Model(&StoreProduct{}).
		Where("id IN ?", ids).
		Update("is_out_of_stock", isOutOfStock).Error
}

func updateStoreAdditiveStockFlags(tx *gorm.DB, ids []uint, isOutOfStock bool) error {
	if len(ids) == 0 {
		return nil
	}
	return tx.Model(&StoreAdditive{}).
		Where("id IN ?", ids).
		Update("is_out_of_stock", isOutOfStock).Error
}

func getStoreProductIDsByIngredients(tx *gorm.DB, storeID uint, ingredientIDs []uint) ([]uint, error) {
	var byIngredients []uint
	err := tx.Model(&StoreProductSize{}).
		Select("DISTINCT store_product_sizes.store_product_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_ingredients ON product_size_ingredients.product_size_id = store_product_sizes.product_size_id").
		Where("store_products.store_id = ?", storeID).
		Where("product_size_ingredients.ingredient_id IN ?", ingredientIDs).
		Pluck("store_product_sizes.store_product_id", &byIngredients).Error
	if err != nil {
		return nil, err
	}

	var byAdditives []uint
	err = tx.Model(&StoreProductSize{}).
		Select("DISTINCT store_product_sizes.store_product_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_additives ON product_size_additives.product_size_id = store_product_sizes.product_size_id AND product_size_additives.is_default = true").
		Joins("JOIN additive_ingredients ON additive_ingredients.additive_id = product_size_additives.additive_id").
		Where("store_products.store_id = ?", storeID).
		Where("additive_ingredients.ingredient_id IN ?", ingredientIDs).
		Pluck("store_product_sizes.store_product_id", &byAdditives).Error
	if err != nil {
		return nil, err
	}

	return unionUintSlices(byIngredients, byAdditives), nil
}

func getStoreAdditiveIDsByIngredients(tx *gorm.DB, storeID uint, ingredientIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&StoreAdditive{}).
		Select("DISTINCT store_additives.id").
		Joins("JOIN additive_ingredients ON additive_ingredients.additive_id = store_additives.additive_id").
		Where("store_additives.store_id = ?", storeID).
		Where("additive_ingredients.ingredient_id IN ?", ingredientIDs).
		Pluck("store_additives.id", &ids).Error
	return ids, err
}

func getOutOfStockStoreProductIDs(tx *gorm.DB, storeProductIDs []uint) ([]uint, error) {
	var outByIngredients []uint
	var outByAdditives []uint

	err := tx.Model(&StoreProductSize{}).
		Select("DISTINCT store_product_sizes.store_product_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_ingredients ON product_size_ingredients.product_size_id = store_product_sizes.product_size_id").
		Joins("JOIN store_stocks ON store_stocks.ingredient_id = product_size_ingredients.ingredient_id AND store_stocks.store_id = store_products.store_id").
		Where("store_products.id IN ?", storeProductIDs).
		Where("store_stocks.quantity < product_size_ingredients.quantity").
		Pluck("store_product_sizes.store_product_id", &outByIngredients).Error
	if err != nil {
		return nil, err
	}

	err = tx.Model(&StoreProductSize{}).
		Select("DISTINCT store_product_sizes.store_product_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_additives ON product_size_additives.product_size_id = store_product_sizes.product_size_id AND product_size_additives.is_default = TRUE").
		Joins("JOIN additive_ingredients ON additive_ingredients.additive_id = product_size_additives.additive_id").
		Joins("JOIN store_stocks ON store_stocks.ingredient_id = additive_ingredients.ingredient_id").
		Where("store_products.id IN ?", storeProductIDs).
		Where("store_stocks.quantity < additive_ingredients.quantity").
		Pluck("store_product_sizes.store_product_id", &outByAdditives).Error
	if err != nil {
		return nil, err
	}

	logrus.Infof("OutByIngredients: %v, OutByAdditives: %v", outByIngredients, outByAdditives)
	return unionUintSlices(outByIngredients, outByAdditives), nil
}

func getOutOfStockStoreAdditiveIDs(tx *gorm.DB, storeAdditiveIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&StoreAdditive{}).
		Select("DISTINCT store_additives.id").
		Joins("JOIN additive_ingredients ON additive_ingredients.additive_id = store_additives.additive_id").
		Joins("JOIN store_stocks ON store_stocks.ingredient_id = additive_ingredients.ingredient_id AND store_stocks.store_id = store_additives.store_id").
		Where("store_additives.id IN ?", storeAdditiveIDs).
		Where("store_stocks.quantity < additive_ingredients.quantity").
		Pluck("store_additives.id", &ids).Error
	return ids, err
}

func getStoreProductIDsByProductSizes(tx *gorm.DB, storeID uint, productSizeIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&StoreProduct{}).
		Distinct("store_products.id").
		Joins("JOIN store_product_sizes ON store_products.id = store_product_sizes.store_product_id").
		Where("store_products.store_id = ?", storeID).
		Where("store_product_sizes.product_size_id IN ?", productSizeIDs).
		Pluck("store_product_id", &ids).Error

	return ids, err
}

func diffUintSlice(all, subset []uint) []uint {
	m := make(map[uint]struct{}, len(subset))
	for _, id := range subset {
		m[id] = struct{}{}
	}

	var diff []uint
	for _, id := range all {
		if _, exists := m[id]; !exists {
			diff = append(diff, id)
		}
	}
	return diff
}

func unionUintSlices(a, b []uint) []uint {
	seen := make(map[uint]struct{})
	for _, v := range a {
		seen[v] = struct{}{}
	}
	for _, v := range b {
		seen[v] = struct{}{}
	}
	result := make([]uint, 0, len(seen))
	for id := range seen {
		result = append(result, id)
	}
	return result
}
