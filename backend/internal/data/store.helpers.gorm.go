package data

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type usageRow struct {
	EntityID     uint
	IngredientID uint
	RequiredQty  float64
}

func RecalculateOutOfStock(tx *gorm.DB, storeID uint, ingredientIDs []uint, productSizeIDs []uint) error {
	start := time.Now()
	logrus.Info("=================Start Recalculation===================")
	logrus.Infof("initial values: storeID=%v ingredientIDs=%v productSizeIDs=%v", storeID, ingredientIDs, productSizeIDs)

	if storeID == 0 || (len(ingredientIDs) == 0 && len(productSizeIDs) == 0) {
		return nil
	}

	var productSizesIngredientIDs,
		storeProductIDsFromPS,
		storeProductIDsFromIngredients,
		storeAdditiveIDs []uint
	var err error

	if len(productSizeIDs) > 0 {
		storeProductIDsFromPS, err = getStoreProductIDsByProductSizes(tx, storeID, productSizeIDs)
		if err != nil {
			return err
		}

		productSizesIngredientIDs, err = getAllIngredientIDsByProductSizes(tx, productSizesIngredientIDs)
		if err != nil {
			return err
		}

		if len(productSizesIngredientIDs) > 0 {
			ingredientIDs = utils.UnionSlices(ingredientIDs, productSizesIngredientIDs)
		}
		logrus.Infof("ingredients from product size: %v", productSizesIngredientIDs)
		logrus.Infof("ingredients initial: %v", ingredientIDs)
		logrus.Infof("ingredients diff: %v", utils.DiffSlice(ingredientIDs, productSizesIngredientIDs))
	}

	if len(ingredientIDs) > 0 {
		frozenStockMap, err := CalculateFrozenStock(tx, storeID, ingredientIDs)
		if err != nil {
			return err
		}
		logrus.Infof("FrozenStockMap: %v", frozenStockMap)

		storeProductIDsFromIngredients, err = getStoreProductIDsByIngredients(tx, storeID, ingredientIDs)
		if err != nil {
			return err
		}

		if len(storeProductIDsFromPS) > 0 || len(storeProductIDsFromIngredients) > 0 {
			mergedStoreProductIDs := utils.UnionSlices(storeProductIDsFromPS, storeProductIDsFromIngredients)

			if err := RecalculateStoreProducts(tx, mergedStoreProductIDs, frozenStockMap, storeID); err != nil {
				return err
			}
		}

		storeAdditiveIDs, err = getStoreAdditiveIDsByIngredients(tx, storeID, ingredientIDs)
		if err != nil {
			return err
		}

		if len(storeAdditiveIDs) > 0 {
			if err := RecalculateStoreAdditives(tx, storeAdditiveIDs, storeID, frozenStockMap); err != nil {
				return err
			}
		}
	}

	logrus.Infof("=============Estimated time: %v==================", time.Since(start))

	return nil
}

func RecalculateStoreProducts(
	tx *gorm.DB,
	storeProductIDs []uint,
	frozenStock map[uint]float64,
	storeID uint,
) error {
	if len(storeProductIDs) == 0 {
		return nil
	}

	outOfStockIDs, err := getOutOfStockStoreProductIDs(tx, storeProductIDs, storeID, frozenStock)
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

func RecalculateStoreAdditives(
	tx *gorm.DB,
	storeAdditiveIDs []uint,
	storeID uint,
	frozenStock map[uint]float64,
) error {
	if len(storeAdditiveIDs) == 0 {
		return nil
	}

	outOfStockIDs, err := getOutOfStockStoreAdditiveIDs(tx, storeAdditiveIDs, storeID, frozenStock)
	if err != nil {
		return err
	}
	inStockIDs := utils.DiffSlice(storeAdditiveIDs, outOfStockIDs)

	// Далее стандартно
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

	return utils.UnionSlices(byIngredients, byAdditives), nil
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

func getOutOfStockStoreProductIDs(
	tx *gorm.DB,
	storeProductIDs []uint,
	storeID uint,
	frozenStock map[uint]float64,
) ([]uint, error) {
	outByIngredients, err := getOutByIngredients(tx, storeProductIDs, storeID, frozenStock)
	if err != nil {
		return nil, err
	}

	outByAdditives, err := getOutByDefaultAdditives(tx, storeProductIDs, storeID, frozenStock)
	if err != nil {
		return nil, err
	}

	return utils.UnionSlices(outByIngredients, outByAdditives), nil
}

func getOutByIngredients(
	tx *gorm.DB,
	storeProductIDs []uint,
	storeID uint,
	frozenStock map[uint]float64,
) ([]uint, error) {
	var usageRows []usageRow
	err := tx.Model(&StoreProductSize{}).
		Select(`DISTINCT store_product_sizes.store_product_id AS entity_id,
                product_size_ingredients.ingredient_id,
                product_size_ingredients.quantity AS required_qty`).
		Joins(`JOIN store_products
               ON store_products.id = store_product_sizes.store_product_id`).
		Joins(`JOIN product_size_ingredients
               ON product_size_ingredients.product_size_id = store_product_sizes.product_size_id`).
		Where(`store_products.id IN ?`, storeProductIDs).
		Scan(&usageRows).Error
	if err != nil {
		return nil, err
	}
	if len(usageRows) == 0 {
		return nil, nil
	}

	ingrSet := make(map[uint]struct{})
	for _, row := range usageRows {
		ingrSet[row.IngredientID] = struct{}{}
	}
	var neededIDs []uint
	for ingrID := range ingrSet {
		neededIDs = append(neededIDs, ingrID)
	}

	var stocks []StoreStock
	err = tx.Model(&StoreStock{}).
		Where("store_id = ?", storeID).
		Where("ingredient_id IN ?", neededIDs).
		Find(&stocks).Error
	if err != nil {
		return nil, err
	}

	stockMap := make(map[uint]float64, len(stocks))
	for _, s := range stocks {
		adjustedQty := s.Quantity - frozenStock[s.IngredientID]
		stockMap[s.IngredientID] = adjustedQty
	}

	outSet := make(map[uint]struct{})
	for _, row := range usageRows {
		available := stockMap[row.IngredientID]
		logrus.Infof("%v < %v", available, row.RequiredQty)
		if available < row.RequiredQty {
			outSet[row.EntityID] = struct{}{}
		}
	}

	var outIDs []uint
	for spID := range outSet {
		logrus.Infof("out ID %v", spID)
		outIDs = append(outIDs, spID)
	}
	logrus.Infof("TOTAL OUT IDS: %v", outIDs)

	return outIDs, nil
}

func getOutByDefaultAdditives(
	tx *gorm.DB,
	storeProductIDs []uint,
	storeID uint,
	frozenStock map[uint]float64,
) ([]uint, error) {
	var usageRows []usageRow
	err := tx.Model(&StoreProductSize{}).
		Select(`DISTINCT store_product_sizes.store_product_id AS entity_id,
                additive_ingredients.ingredient_id,
                additive_ingredients.quantity AS required_qty`).
		Joins(`JOIN store_products
               ON store_products.id = store_product_sizes.store_product_id`).
		Joins(`JOIN product_size_additives
               ON product_size_additives.product_size_id = store_product_sizes.product_size_id
               AND product_size_additives.is_default = TRUE`).
		Joins(`JOIN additive_ingredients
               ON additive_ingredients.additive_id = product_size_additives.additive_id`).
		Where(`store_products.id IN ?`, storeProductIDs).
		Scan(&usageRows).Error
	if err != nil {
		return nil, err
	}
	if len(usageRows) == 0 {
		return nil, nil
	}

	ingrSet := make(map[uint]struct{})
	for _, row := range usageRows {
		ingrSet[row.IngredientID] = struct{}{}
	}
	var neededIDs []uint
	for ingrID := range ingrSet {
		neededIDs = append(neededIDs, ingrID)
	}

	var stocks []StoreStock
	err = tx.Model(&StoreStock{}).
		Where("store_id = ?", storeID).
		Where("ingredient_id IN ?", neededIDs).
		Find(&stocks).Error
	if err != nil {
		return nil, err
	}

	stockMap := make(map[uint]float64, len(stocks))
	for _, s := range stocks {
		adjusted := s.Quantity - frozenStock[s.IngredientID]
		stockMap[s.IngredientID] = adjusted
	}

	outSet := make(map[uint]struct{})
	for _, row := range usageRows {
		available := stockMap[row.IngredientID]
		if available < row.RequiredQty {
			outSet[row.EntityID] = struct{}{}
		}
	}

	var outIDs []uint
	for spID := range outSet {
		outIDs = append(outIDs, spID)
	}
	return outIDs, nil
}

func getOutOfStockStoreAdditiveIDs(
	tx *gorm.DB,
	storeAdditiveIDs []uint,
	storeID uint,
	frozenStock map[uint]float64,
) ([]uint, error) {
	if len(storeAdditiveIDs) == 0 {
		return nil, nil
	}

	var usageRows []usageRow
	err := tx.Model(&StoreAdditive{}).
		Select(`DISTINCT store_additives.id AS entity_id,
                additive_ingredients.ingredient_id,
                additive_ingredients.quantity AS required_qty`).
		Joins(`JOIN additive_ingredients
               ON additive_ingredients.additive_id = store_additives.additive_id`).
		Where(`store_additives.id IN ?`, storeAdditiveIDs).
		Scan(&usageRows).Error
	if err != nil {
		return nil, err
	}
	if len(usageRows) == 0 {
		return nil, nil
	}

	// 2) Собираем ingredientIDs, которые реально встречаются
	ingrSet := make(map[uint]struct{})
	for _, row := range usageRows {
		ingrSet[row.IngredientID] = struct{}{}
	}
	var neededIngredientIDs []uint
	for ingrID := range ingrSet {
		neededIngredientIDs = append(neededIngredientIDs, ingrID)
	}

	var stocks []StoreStock
	err = tx.Model(&StoreStock{}).
		Where("store_id = ?", storeID).
		Where("ingredient_id IN ?", neededIngredientIDs).
		Find(&stocks).Error
	if err != nil {
		return nil, err
	}

	stockMap := make(map[uint]float64, len(stocks))
	for _, s := range stocks {
		adjusted := s.Quantity - frozenStock[s.IngredientID] // вычитаем замороженное
		stockMap[s.IngredientID] = adjusted
	}

	outOfStockSet := make(map[uint]struct{})
	for _, row := range usageRows {
		available := stockMap[row.IngredientID]
		if available < row.RequiredQty {
			outOfStockSet[row.EntityID] = struct{}{}
		}
	}

	var outIDs []uint
	for saID := range outOfStockSet {
		outIDs = append(outIDs, saID)
	}

	return outIDs, nil
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

func getAllIngredientIDsByProductSizes(tx *gorm.DB, productSizeIDs []uint) ([]uint, error) {
	if len(productSizeIDs) == 0 {
		return nil, nil
	}

	directIDs, err := getProductSizeDirectIngredientsByProductSizeIDs(tx, productSizeIDs)
	if err != nil {
		return nil, err
	}

	additiveIngredientIDs, err := getDefaultAdditiveIngredientsByProductSizeIDs(tx, productSizeIDs)
	if err != nil {
		return nil, err
	}

	return utils.UnionSlices(directIDs, additiveIngredientIDs), nil
}

func getProductSizeDirectIngredientsByProductSizeIDs(tx *gorm.DB, productSizeIDs []uint) ([]uint, error) {
	var directIDs []uint
	err := tx.Model(&ProductSizeIngredient{}).
		Distinct("product_size_ingredients.ingredient_id").
		Where("product_size_ingredients.product_size_id IN ?", productSizeIDs).
		Pluck("product_size_ingredients.ingredient_id", &directIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch direct ingredientIDs: %w", err)
	}

	return directIDs, nil
}

func getDefaultAdditiveIngredientsByProductSizeIDs(tx *gorm.DB, productSizeIDs []uint) ([]uint, error) {
	var additiveIngredientIDs []uint
	err := tx.Model(&AdditiveIngredient{}).
		Distinct("additive_ingredients.ingredient_id").
		Joins("JOIN product_size_additives ON product_size_additives.additive_id = additive_ingredients.additive_id").
		Where("product_size_additives.product_size_id IN ?", productSizeIDs).
		Where("product_size_additives.is_default = TRUE").
		Pluck("additive_ingredients.ingredient_id", &additiveIngredientIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch additive ingredientIDs: %w", err)
	}
	return additiveIngredientIDs, nil
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
