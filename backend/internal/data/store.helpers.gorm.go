package data

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RecalculateOutOfStock(tx *gorm.DB, storeID uint, ingredientIDs []uint, productSizeIDs []uint, frozenStockMap map[uint]float64) error {
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
		mergedStoreProductIDs := utils.MergeDistinct(storeProductIDsFromPS, storeProductIDsFromIngredients)

		if err := RecalculateStoreProducts(tx, mergedStoreProductIDs); err != nil {
			return err
		}
	}

	if len(ingredientIDs) > 0 {
		storeAdditiveIDs, err = getStoreAdditiveIDsByIngredients(tx, storeID, ingredientIDs)
		if err != nil {
			return err
		}
		if err := RecalculateStoreAdditives(tx, storeAdditiveIDs); err != nil {
			return err
		}
	}

	logrus.Info("=============================================")

	return nil
}

func RecalculateStoreProducts(tx *gorm.DB, storeProductIDs []uint) error {
	if len(storeProductIDs) == 0 {
		return nil
	}

	outOfStockIDs, err := getOutOfStockStoreProductIDs(tx, storeProductIDs)
	if err != nil {
		return err
	}
	inStockIDs := utils.DiffSlice(storeProductIDs, outOfStockIDs)
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

func RecalculateStoreAdditives(tx *gorm.DB, storeAdditiveIDs []uint) error {
	if len(storeAdditiveIDs) == 0 {
		return nil
	}

	outOfStockIDs, err := getOutOfStockStoreAdditiveIDs(tx, storeAdditiveIDs)
	if err != nil {
		return err
	}
	inStockIDs := utils.DiffSlice(storeAdditiveIDs, outOfStockIDs)
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
	if storeID == 0 || len(ingredientIDs) == 0 {
		return nil, nil
	}

	byIngredients, err := getStoreProductIDsByIngredientUsage(tx, storeID, ingredientIDs)
	if err != nil {
		return nil, err
	}

	byAdditives, err := getStoreProductIDsByDefaultAdditiveUsage(tx, storeID, ingredientIDs)
	if err != nil {
		return nil, err
	}

	return utils.MergeDistinct(byIngredients, byAdditives), nil
}

func getStoreProductIDsByIngredientUsage(tx *gorm.DB, storeID uint, ingredientIDs []uint) ([]uint, error) {
	var productIDs []uint

	err := tx.Model(&StoreProductSize{}).
		Select("DISTINCT store_product_sizes.store_product_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_ingredients ON product_size_ingredients.product_size_id = store_product_sizes.product_size_id").
		Where("store_products.store_id = ?", storeID).
		Where("product_size_ingredients.ingredient_id IN ?", ingredientIDs).
		Pluck("store_product_sizes.store_product_id", &productIDs).Error

	return productIDs, err
}

func getStoreProductIDsByDefaultAdditiveUsage(tx *gorm.DB, storeID uint, ingredientIDs []uint) ([]uint, error) {
	var productIDs []uint

	err := tx.Model(&StoreProductSize{}).
		Select("DISTINCT store_product_sizes.store_product_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_additives ON product_size_additives.product_size_id = store_product_sizes.product_size_id AND product_size_additives.is_default = true").
		Joins("JOIN additive_ingredients ON additive_ingredients.additive_id = product_size_additives.additive_id").
		Where("store_products.store_id = ?", storeID).
		Where("additive_ingredients.ingredient_id IN ?", ingredientIDs).
		Pluck("store_product_sizes.store_product_id", &productIDs).Error

	return productIDs, err
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
	outByIngredients, err := getOutByIngredients(tx, storeProductIDs)
	if err != nil {
		return nil, err
	}

	outByAdditives, err := getOutByDefaultAdditives(tx, storeProductIDs)
	if err != nil {
		return nil, err
	}

	return utils.MergeDistinct(outByIngredients, outByAdditives), nil
}

func getOutByIngredients(tx *gorm.DB, storeProductIDs []uint) ([]uint, error) {
	var ids []uint

	err := tx.Model(&StoreProductSize{}).
		Select("DISTINCT store_product_sizes.store_product_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_ingredients ON product_size_ingredients.product_size_id = store_product_sizes.product_size_id").
		Joins("JOIN store_stocks ON store_stocks.ingredient_id = product_size_ingredients.ingredient_id AND store_stocks.store_id = store_products.store_id").
		Where("store_products.id IN ?", storeProductIDs).
		Where("store_stocks.quantity < product_size_ingredients.quantity").
		Pluck("store_product_sizes.store_product_id", &ids).Error

	return ids, err
}

func getOutByDefaultAdditives(tx *gorm.DB, storeProductIDs []uint) ([]uint, error) {
	var ids []uint

	err := tx.Model(&StoreProductSize{}).
		Select("DISTINCT store_product_sizes.store_product_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_additives ON product_size_additives.product_size_id = store_product_sizes.product_size_id AND product_size_additives.is_default = TRUE").
		Joins("JOIN additive_ingredients ON additive_ingredients.additive_id = product_size_additives.additive_id").
		Joins("JOIN store_additives ON store_additives.additive_id = product_size_additives.additive_id AND store_additives.store_id = store_products.store_id").
		Joins("JOIN store_stocks ON store_stocks.ingredient_id = additive_ingredients.ingredient_id AND store_stocks.store_id = store_additives.store_id").
		Where("store_products.id IN ?", storeProductIDs).
		Where("store_stocks.quantity < additive_ingredients.quantity").
		Pluck("store_product_sizes.store_product_id", &ids).Error

	return ids, err
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

func CalculateFrozenStock(tx *gorm.DB, storeID uint, ingredientIDs []uint) (map[uint]float64, error) {
	frozenStock := make(map[uint]float64)

	orders, err := loadActiveOrders(tx, storeID, ingredientIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load active orders for store %d: %w", storeID, err)
	}

	for _, order := range orders {
		for _, sub := range order.Suborders {
			if !isSuborderActive(sub) {
				continue
			}
			accumulateProductUsage(&frozenStock, sub, ingredientIDs)
			accumulateAdditiveUsage(&frozenStock, sub, ingredientIDs)
		}
	}

	return frozenStock, nil
}

func loadActiveOrders(tx *gorm.DB, storeID uint, ingredientIDs []uint) ([]Order, error) {
	var orders []Order

	query := tx.Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.Ingredients.Ingredient").
		Preload("Suborders.StoreProductSize.ProductSize.ProductSizeIngredients.Ingredient").
		Where("store_id = ?", storeID).
		Where("status IN ?", []OrderStatus{
			OrderStatusWaitingForPayment,
			OrderStatusPending,
			OrderStatusPreparing,
		})

	if len(ingredientIDs) > 0 {
		query = query.
			Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.Ingredients", "ingredient_id IN ?", ingredientIDs).
			Preload("Suborders.StoreProductSize.ProductSize.ProductSizeIngredients", "ingredient_id IN ?", ingredientIDs)
	}

	err := query.Find(&orders).Error
	return orders, err
}

func isSuborderActive(sub Suborder) bool {
	return sub.Status == SubOrderStatusPending || sub.Status == SubOrderStatusPreparing
}

func accumulateProductUsage(frozenStock *map[uint]float64, sub Suborder, ingredientIDs []uint) {
	for _, usage := range sub.StoreProductSize.ProductSize.ProductSizeIngredients {
		if ingredientIDs == nil || contains(ingredientIDs, usage.IngredientID) {
			(*frozenStock)[usage.IngredientID] += usage.Quantity
		}
	}
}

func accumulateAdditiveUsage(frozenStock *map[uint]float64, sub Suborder, ingredientIDs []uint) {
	for _, subAdd := range sub.SuborderAdditives {
		for _, ingrUsage := range subAdd.StoreAdditive.Additive.Ingredients {
			if ingredientIDs == nil || contains(ingredientIDs, ingrUsage.IngredientID) {
				(*frozenStock)[ingrUsage.IngredientID] += ingrUsage.Quantity
			}
		}
	}
}

func contains(slice []uint, item uint) bool {
	for _, id := range slice {
		if id == item {
			return true
		}
	}
	return false
}
