package storeInventoryManagers

import (
	"fmt"
	"slices"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeProvisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func deductStoreStocks(
	tx *gorm.DB,
	storeID uint,
	requiredIngredientQtyMap map[uint]float64,
) (map[uint]*data.StoreStock, error) {
	if len(requiredIngredientQtyMap) == 0 {
		return nil, nil
	}

	// 1) Load all relevant stocks in one query.
	ids := make([]uint, 0, len(requiredIngredientQtyMap))
	for id := range requiredIngredientQtyMap {
		ids = append(ids, id)
	}

	var stocks []data.StoreStock
	if err := tx.
		Preload("Ingredient").
		Preload("Ingredient.Unit").
		Where("store_id = ? AND ingredient_id IN ?", storeID, ids).
		Find(&stocks).Error; err != nil {
		return nil, fmt.Errorf("failed to load store stocks: %w", err)
	}

	// Map for lookup.
	stockMap := make(map[uint]*data.StoreStock, len(stocks))
	for i := range stocks {
		s := &stocks[i]
		stockMap[s.IngredientID] = s
	}

	// 2) Deduct each required amount.
	result := make(map[uint]*data.StoreStock, len(requiredIngredientQtyMap))
	for ingrID, reqQty := range requiredIngredientQtyMap {
		s, ok := stockMap[ingrID]
		if !ok {
			return nil, fmt.Errorf("stock not found for ingredient ID %d", ingrID)
		}
		if s.Quantity < reqQty {
			return nil, fmt.Errorf("%w: insufficient stock for ingredient ID %d",
				storeStocksTypes.ErrInsufficientStock, ingrID)
		}
		s.Quantity -= reqQty
		if err := tx.Save(s).Error; err != nil {
			return nil, fmt.Errorf("failed to save stock for ingredient %d: %w", ingrID, err)
		}
		logrus.Infof("deducted %.2f from ingredient %d", reqQty, ingrID)
		result[ingrID] = s
	}
	return result, nil
}

func deductStoreProvisions(
	tx *gorm.DB,
	storeID uint,
	requiredProvisionVolMap map[uint]float64,
) (map[uint][]data.StoreProvision, error) {
	if len(requiredProvisionVolMap) == 0 {
		return nil, nil
	}

	// 1) Load all relevant provisions in one query.
	provIDs := make([]uint, 0, len(requiredProvisionVolMap))
	for id := range requiredProvisionVolMap {
		provIDs = append(provIDs, id)
	}

	var allProv []data.StoreProvision
	if err := tx.
		Where("store_id = ? AND provision_id IN ? AND status = ? AND (expires_at IS NULL OR expires_at > ?)",
			storeID, provIDs, data.STORE_PROVISION_STATUS_COMPLETED, time.Now().UTC()).
		Order("created_at ASC").
		Find(&allProv).Error; err != nil {
		return nil, fmt.Errorf("failed to load store provisions: %w", err)
	}

	// Group by provisionID
	group := make(map[uint][]data.StoreProvision, len(requiredProvisionVolMap))
	for _, p := range allProv {
		group[p.ProvisionID] = append(group[p.ProvisionID], p)
	}

	// 2) Deduct for each provision ID
	result := make(map[uint][]data.StoreProvision, len(requiredProvisionVolMap))
	for provID, reqVol := range requiredProvisionVolMap {
		remain := reqVol
		used := []data.StoreProvision{}

		list := group[provID]
		for i := range list {
			if remain <= 0 {
				break
			}
			p := &list[i]
			avail := p.Volume
			delta := avail
			if delta > remain {
				delta = remain
			}
			if delta == 0 {
				continue
			}
			p.Volume -= delta
			remain -= delta
			if p.Volume == 0 {
				p.Status = data.STORE_PROVISION_STATUS_EMPTY
			}
			if err := tx.Save(p).Error; err != nil {
				return nil, fmt.Errorf("failed to update provision %d: %w", p.ID, err)
			}
			logrus.Infof("deducted %.2f from provision %d", delta, provID)
			used = append(used, *p)
		}

		if remain > 0 {
			return nil, fmt.Errorf("%w: not enough volume for provision %d",
				storeProvisionsTypes.ErrInsufficientStoreProvision, provID)
		}
		result[provID] = used
	}
	return result, nil
}

func recalculateStoreProducts(
	tx *gorm.DB,
	storeProductIDs []uint,
	frozenInventory *types.FrozenInventory,
	storeID uint,
) error {
	if len(storeProductIDs) == 0 {
		return nil
	}

	outOfStockIDs, err := getOutOfStockStoreProductIDs(tx, storeProductIDs, storeID, frozenInventory)
	if err != nil {
		return err
	}
	inStockIDs := utils.DiffSlice(storeProductIDs, outOfStockIDs)

	if err := updateStoreProductStockFlags(tx, outOfStockIDs, true); err != nil {
		return err
	}
	if err := updateStoreProductStockFlags(tx, inStockIDs, false); err != nil {
		return err
	}

	return nil
}

func recalculateStoreAdditives(
	tx *gorm.DB,
	storeAdditiveIDs []uint,
	storeID uint,
	frozenInventory *types.FrozenInventory,
) error {
	if len(storeAdditiveIDs) == 0 {
		return nil
	}

	outOfStockIDs, err := getOutOfStockStoreAdditiveIDs(tx, storeAdditiveIDs, storeID, frozenInventory)
	logrus.Infof("outOfstock additives: %v", outOfStockIDs)
	if err != nil {
		return err
	}
	inStockIDs := utils.DiffSlice(storeAdditiveIDs, outOfStockIDs)

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
	return tx.Model(&data.StoreProduct{}).
		Where("id IN ?", ids).
		Update("is_out_of_stock", isOutOfStock).Error
}

func updateStoreAdditiveStockFlags(tx *gorm.DB, ids []uint, isOutOfStock bool) error {
	if len(ids) == 0 {
		return nil
	}
	return tx.Model(&data.StoreAdditive{}).
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

func getStoreProductIDsByProvisions(tx *gorm.DB, storeID uint, provisionIDs []uint) ([]uint, error) {
	if storeID == 0 || len(provisionIDs) == 0 {
		return nil, nil
	}

	byProvisions, err := getStoreProductIDsByProvisionUsage(tx, storeID, provisionIDs)
	if err != nil {
		return nil, err
	}

	byAdditives, err := getStoreProductIDsByDefaultAdditiveProvisionUsage(tx, storeID, provisionIDs)
	if err != nil {
		return nil, err
	}

	return utils.UnionSlices(byProvisions, byAdditives), nil
}

func getStoreProductIDsByProvisionUsage(tx *gorm.DB, storeID uint, provisionIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&data.StoreProductSize{}).
		Select("DISTINCT store_product_sizes.store_product_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_provisions ON product_size_provisions.product_size_id = store_product_sizes.product_size_id").
		Where("store_products.store_id = ?", storeID).
		Where("product_size_provisions.provision_id IN ?", provisionIDs).
		Pluck("store_product_sizes.store_product_id", &ids).Error
	return ids, err
}

func getStoreProductIDsByDefaultAdditiveProvisionUsage(tx *gorm.DB, storeID uint, provisionIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&data.StoreProductSize{}).
		Select("DISTINCT store_product_sizes.store_product_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_additives ON product_size_additives.product_size_id = store_product_sizes.product_size_id AND product_size_additives.is_default = TRUE").
		Joins("JOIN additive_provisions ON additive_provisions.additive_id = product_size_additives.additive_id").
		Where("store_products.store_id = ?", storeID).
		Where("additive_provisions.provision_id IN ?", provisionIDs).
		Pluck("store_product_sizes.store_product_id", &ids).Error
	return ids, err
}

func getStoreProductIDsByIngredientUsage(tx *gorm.DB, storeID uint, ingredientIDs []uint) ([]uint, error) {
	var productIDs []uint

	err := tx.Model(&data.StoreProductSize{}).
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

	err := tx.Model(&data.StoreProductSize{}).
		Select("DISTINCT store_product_sizes.store_product_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Joins("JOIN product_size_additives ON product_size_additives.product_size_id = store_product_sizes.product_size_id AND product_size_additives.is_default = true").
		Joins("JOIN additive_ingredients ON additive_ingredients.additive_id = product_size_additives.additive_id").
		Where("store_products.store_id = ?", storeID).
		Where("additive_ingredients.ingredient_id IN ?", ingredientIDs).
		Pluck("store_product_sizes.store_product_id", &productIDs).Error

	return productIDs, err
}

func getStoreAdditiveIDsByAdditives(tx *gorm.DB, storeID uint, additiveIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&data.StoreAdditive{}).
		Select("DISTINCT store_additives.id").
		Joins("JOIN additives ON additives.id = store_additives.additive_id").
		Where("store_additives.store_id = ?", storeID).
		Where("additives.id IN ?", additiveIDs).
		Pluck("store_additives.id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func getStoreAdditiveIDsByIngredients(tx *gorm.DB, storeID uint, ingredientIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&data.StoreAdditive{}).
		Select("DISTINCT store_additives.id").
		Joins("JOIN additive_ingredients ON additive_ingredients.additive_id = store_additives.additive_id").
		Where("store_additives.store_id = ?", storeID).
		Where("additive_ingredients.ingredient_id IN ?", ingredientIDs).
		Pluck("store_additives.id", &ids).Error
	return ids, err
}

func getStoreAdditiveIDsByProvisions(tx *gorm.DB, storeID uint, provisionIDs []uint) ([]uint, error) {
	if storeID == 0 || len(provisionIDs) == 0 {
		return nil, nil
	}

	var additiveIDs []uint
	err := tx.Model(&data.StoreAdditive{}).
		Distinct("store_additives.id").
		Joins("JOIN additive_provisions ON additive_provisions.additive_id = store_additives.additive_id").
		Where("store_additives.store_id = ?", storeID).
		Where("additive_provisions.provision_id IN ?", provisionIDs).
		Pluck("store_additives.id", &additiveIDs).Error

	return additiveIDs, err
}

func getOutOfStockStoreProductIDs(
	tx *gorm.DB,
	storeProductIDs []uint,
	storeID uint,
	frozenInventory *types.FrozenInventory,
) ([]uint, error) {
	if len(storeProductIDs) == 0 {
		return nil, nil
	}
	if frozenInventory == nil {
		return nil, fmt.Errorf("frozen inventory cannot be nil")
	}

	// 1) get all store_product_size IDs for these storeProductIDs
	storeProdSizeIDs, sizeToProd, err := storeProdSizeToProduct(tx, storeProductIDs)
	if err != nil {
		return nil, err
	}
	if len(storeProdSizeIDs) == 0 {
		return nil, nil // no sizes => no out-of-stock
	}

	// 2) aggregate usage by store_product_size
	usage, err := getAggregatedUsageForStoreProdSizes(tx, storeProdSizeIDs)
	if err != nil {
		return nil, err
	}
	logrus.Infof("usages: %v", usage)

	// 3) gather all needed ingredientIDs / provisionIDs
	ingredientSet := make(map[uint]struct{})
	for key := range usage.Ingredient {
		ingredientSet[key.ResourceID] = struct{}{}
	}
	provisionSet := make(map[uint]struct{})
	for key := range usage.Provision {
		provisionSet[key.ResourceID] = struct{}{}
	}

	var ingredientIDs, provisionIDs []uint
	for ingID := range ingredientSet {
		ingredientIDs = append(ingredientIDs, ingID)
	}
	for provID := range provisionSet {
		provisionIDs = append(provisionIDs, provID)
	}

	// 4) load (stock - frozen) for these ingredients/provisions
	stockMap, provisionMap, err := buildStockMaps(tx, storeID, ingredientIDs, provisionIDs, frozenInventory)
	if err != nil {
		return nil, err
	}
	// TODO remove
	if vol, exists := provisionMap[5]; exists {
		logrus.Infof("                                      requried: %v  inStock: %v; ", usage.Provision[usageKey{74, 5}], vol)
	}

	// We track which storeProducts are out-of-stock
	outSet := make(map[uint]struct{})

	// 5) check each store_product_size usage vs available
	// If any size fails => entire storeProduct is out-of-stock
	for key, needed := range usage.Ingredient {
		available := stockMap[key.ResourceID]
		if needed > available {
			outSet[sizeToProd[key.StoreProductSizeID]] = struct{}{}
		}
	}
	for key, needed := range usage.Provision {
		available := provisionMap[key.ResourceID]
		if needed > available {
			outSet[sizeToProd[key.StoreProductSizeID]] = struct{}{}
		}
	}

	// 6) convert outSet to slice
	var outOfStockIDs []uint
	for prodID := range outSet {
		outOfStockIDs = append(outOfStockIDs, prodID)
	}

	return outOfStockIDs, nil
}

func buildStockMaps(
	tx *gorm.DB,
	storeID uint,
	ingredientIDs, provisionIDs []uint,
	frozen *types.FrozenInventory,
) (map[uint]float64, map[uint]float64, error) {
	stockMap := make(map[uint]float64)
	provisionMap := make(map[uint]float64)

	if len(ingredientIDs) > 0 {
		stocks, err := getRelevantStoreStocks(tx, storeID, ingredientIDs)
		if err != nil {
			return nil, nil, err
		}
		for _, s := range stocks {
			// subtract frozen
			total := s.Quantity - frozen.Ingredients[s.IngredientID]
			if total < 0 {
				total = 0
			}
			stockMap[s.IngredientID] = total
		}
	}

	if len(provisionIDs) > 0 {
		provs, err := getRelevantStoreProvisions(tx, storeID, provisionIDs)
		if err != nil {
			return nil, nil, err
		}
		for _, p := range provs {
			provisionMap[p.ProvisionID] += p.Volume
		}
		for pid, frozenQty := range frozen.Provisions {
			if _, ok := provisionMap[pid]; ok {
				newVal := provisionMap[pid] - frozenQty
				if newVal < 0 {
					newVal = 0
				}
				provisionMap[pid] = newVal
			}
		}
	}

	return stockMap, provisionMap, nil
}

func storeProdSizeToProduct(
	tx *gorm.DB,
	storeProductIDs []uint,
) ([]uint, map[uint]uint, error) {
	if len(storeProductIDs) == 0 {
		return nil, nil, nil
	}

	type spSizeRow struct {
		ID             uint
		StoreProductID uint
	}
	var rows []spSizeRow

	err := tx.Model(&data.StoreProductSize{}).
		Select("id, store_product_id").
		Where("store_product_id IN ?", storeProductIDs).
		Find(&rows).Error
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load storeProductSizes: %w", err)
	}

	sizeIDs := make([]uint, 0, len(rows))
	sizeMap := make(map[uint]uint, len(rows))
	for _, r := range rows {
		sizeIDs = append(sizeIDs, r.ID)
		sizeMap[r.ID] = r.StoreProductID
	}

	return sizeIDs, sizeMap, nil
}

// Summation of direct + default-additive usage for each store_product_size
func getAggregatedUsageForStoreProdSizes(
	tx *gorm.DB,
	storeProdSizeIDs []uint,
) (*aggregatedUsage, error) {
	usage := &aggregatedUsage{
		Ingredient: make(map[usageKey]float64),
		Provision:  make(map[usageKey]float64),
	}

	// 1) Direct ingredients
	dirIng, err := getDirectIngredientUsage(tx, storeProdSizeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load direct ingredient usage: %w", err)
	}
	for _, row := range dirIng {
		key := usageKey{
			StoreProductSizeID: row.StoreProductSizeID,
			ResourceID:         row.IngredientID,
		}
		usage.Ingredient[key] += row.RequiredQuantity
	}

	// 2) Direct provisions
	dirProv, err := getDirectProvisionUsage(tx, storeProdSizeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load direct provision usage: %w", err)
	}
	for _, row := range dirProv {
		key := usageKey{
			StoreProductSizeID: row.StoreProductSizeID,
			ResourceID:         row.ProvisionID,
		}
		usage.Provision[key] += row.RequiredVolume
	}

	// 3) Default-additive ingredients
	defIng, err := getDefaultAdditiveIngredientUsage(tx, storeProdSizeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load default-additive ingredient usage: %w", err)
	}
	for _, row := range defIng {
		key := usageKey{
			StoreProductSizeID: row.StoreProductSizeID,
			ResourceID:         row.IngredientID,
		}
		usage.Ingredient[key] += row.RequiredQuantity
	}

	// 4) Default-additive provisions
	defProv, err := getDefaultAdditiveProvisionUsage(tx, storeProdSizeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load default-additive provision usage: %w", err)
	}
	for _, row := range defProv {
		key := usageKey{
			StoreProductSizeID: row.StoreProductSizeID,
			ResourceID:         row.ProvisionID,
		}
		usage.Provision[key] += row.RequiredVolume
	}

	return usage, nil
}

func getDirectIngredientUsage(tx *gorm.DB, storeProdSizeIDs []uint) ([]productSizeIngredientRow, error) {
	if len(storeProdSizeIDs) == 0 {
		return nil, nil
	}
	var rows []productSizeIngredientRow

	// store_product_sizes -> product_sizes -> product_size_ingredients
	err := tx.Model(&data.StoreProductSize{}).
		Select(`
            store_product_sizes.id AS store_product_size_id,
            product_size_ingredients.ingredient_id,
            product_size_ingredients.quantity AS required_quantity
        `).
		Joins(`JOIN product_sizes ON product_sizes.id = store_product_sizes.product_size_id`).
		Joins(`JOIN product_size_ingredients ON product_size_ingredients.product_size_id = product_sizes.id`).
		Where("store_product_sizes.id IN ?", storeProdSizeIDs).
		Scan(&rows).Error

	return rows, err
}

func getDirectProvisionUsage(tx *gorm.DB, storeProdSizeIDs []uint) ([]productSizeProvisionRow, error) {
	if len(storeProdSizeIDs) == 0 {
		return nil, nil
	}
	var rows []productSizeProvisionRow

	err := tx.Model(&data.StoreProductSize{}).
		Select(`
            store_product_sizes.id AS store_product_size_id,
            product_size_provisions.provision_id,
            product_size_provisions.volume AS required_volume
        `).
		Joins(`JOIN product_sizes ON product_sizes.id = store_product_sizes.product_size_id`).
		Joins(`JOIN product_size_provisions ON product_size_provisions.product_size_id = product_sizes.id`).
		Where("store_product_sizes.id IN ?", storeProdSizeIDs).
		Scan(&rows).Error

	return rows, err
}

func getDefaultAdditiveIngredientUsage(tx *gorm.DB, storeProdSizeIDs []uint) ([]productSizeIngredientRow, error) {
	if len(storeProdSizeIDs) == 0 {
		return nil, nil
	}
	var rows []productSizeIngredientRow

	err := tx.Model(&data.StoreProductSize{}).
		Select(`
            store_product_sizes.id AS store_product_size_id,
            additive_ingredients.ingredient_id,
            additive_ingredients.quantity AS required_quantity
        `).
		Joins(`JOIN product_sizes ON product_sizes.id = store_product_sizes.product_size_id`).
		Joins(`JOIN product_size_additives 
               ON product_size_additives.product_size_id = product_sizes.id
               AND product_size_additives.is_default = TRUE`).
		Joins(`JOIN additive_ingredients
               ON additive_ingredients.additive_id = product_size_additives.additive_id`).
		Where("store_product_sizes.id IN ?", storeProdSizeIDs).
		Scan(&rows).Error

	return rows, err
}

func getDefaultAdditiveProvisionUsage(tx *gorm.DB, storeProdSizeIDs []uint) ([]productSizeProvisionRow, error) {
	if len(storeProdSizeIDs) == 0 {
		return nil, nil
	}
	var rows []productSizeProvisionRow

	err := tx.Model(&data.StoreProductSize{}).
		Select(`
            store_product_sizes.id AS store_product_size_id,
            additive_provisions.provision_id,
            additive_provisions.volume AS required_volume
        `).
		Joins(`JOIN product_sizes ON product_sizes.id = store_product_sizes.product_size_id`).
		Joins(`JOIN product_size_additives 
               ON product_size_additives.product_size_id = product_sizes.id
               AND product_size_additives.is_default = TRUE`).
		Joins(`JOIN additive_provisions 
               ON additive_provisions.additive_id = product_size_additives.additive_id`).
		Where("store_product_sizes.id IN ?", storeProdSizeIDs).
		Scan(&rows).Error

	return rows, err
}

func getOutOfStockStoreAdditiveIDs(
	tx *gorm.DB,
	storeAdditiveIDs []uint,
	storeID uint,
	frozen *types.FrozenInventory,
) ([]uint, error) {
	if len(storeAdditiveIDs) == 0 {
		return nil, nil
	}
	if frozen == nil {
		return nil, fmt.Errorf("nil frozen inventory pointer fetched")
	}

	var ingredientRows []additiveIngredientUsageRow
	err := tx.Model(&data.StoreAdditive{}).
		Select(`DISTINCT store_additives.id AS store_product_or_additive_id,
                additive_ingredients.ingredient_id,
                additive_ingredients.quantity AS required_quantity`).
		Joins(`JOIN additive_ingredients
               ON additive_ingredients.additive_id = store_additives.additive_id`).
		Where(`store_additives.id IN ?`, storeAdditiveIDs).
		Scan(&ingredientRows).Error
	if err != nil {
		return nil, err
	}

	var provisionRows []additiveProvisionUsageRow
	err = tx.Model(&data.StoreAdditive{}).
		Select(`DISTINCT store_additives.id AS store_product_or_additive_id,
                additive_provisions.provision_id AS provision_id,
                additive_provisions.volume AS required_volume`).
		Joins(`JOIN additive_provisions
               ON additive_provisions.additive_id = store_additives.additive_id`).
		Where(`store_additives.id IN ?`, storeAdditiveIDs).
		Scan(&provisionRows).Error
	if err != nil {
		return nil, err
	}

	// === LOAD INGREDIENT STOCK ===
	ingredientSet := make(map[uint]struct{})
	for _, row := range ingredientRows {
		ingredientSet[row.IngredientID] = struct{}{}
	}
	var neededIngredientIDs []uint
	for id := range ingredientSet {
		neededIngredientIDs = append(neededIngredientIDs, id)
	}

	stocks, err := getRelevantStoreStocks(tx, storeID, neededIngredientIDs)
	if err != nil {
		return nil, err
	}

	stockMap := make(map[uint]float64)
	for _, s := range stocks {
		stockMap[s.IngredientID] = s.Quantity - frozen.Ingredients[s.IngredientID]
	}

	provisionSet := make(map[uint]struct{})
	for _, row := range provisionRows {
		provisionSet[row.ProvisionID] = struct{}{}
	}
	var neededProvisionIDs []uint
	for id := range provisionSet {
		neededProvisionIDs = append(neededProvisionIDs, id)
	}

	provisions, err := getRelevantStoreProvisions(tx, storeID, neededProvisionIDs)
	if err != nil {
		return nil, err
	}

	provisionStockMap := make(map[uint]float64)
	for _, p := range provisions {
		provisionStockMap[p.ProvisionID] += p.Volume
	}

	for pid, frozenQty := range frozen.Provisions {
		provisionStockMap[pid] -= frozenQty
		if provisionStockMap[pid] < 0 {
			provisionStockMap[pid] = 0
		}
	}

	// === DETECT OUT OF STOCK ===
	outSet := make(map[uint]struct{})

	for _, row := range ingredientRows {
		available := stockMap[row.IngredientID]
		if available < row.RequiredQuantity {
			outSet[row.StoreAdditiveID] = struct{}{}
		}
	}

	for _, row := range provisionRows {
		available := provisionStockMap[row.ProvisionID]
		if available < row.RequiredVolume {
			outSet[row.StoreAdditiveID] = struct{}{}
		}
	}

	var outIDs []uint
	for id := range outSet {
		outIDs = append(outIDs, id)
	}

	return outIDs, nil
}

func getRelevantStoreStocks(tx *gorm.DB, storeID uint, ingredientIDs []uint) ([]data.StoreStock, error) {
	var stocks []data.StoreStock
	err := tx.Model(&data.StoreStock{}).
		Where("store_id = ?", storeID).
		Where("ingredient_id IN ?", ingredientIDs).
		Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func getRelevantStoreProvisions(tx *gorm.DB, storeID uint, provisionIDs []uint) ([]data.StoreProvision, error) {
	var provisions []data.StoreProvision
	err := tx.Model(&data.StoreProvision{}).
		Where("store_id = ?", storeID).
		Where("provision_id IN ?", provisionIDs).
		Where("status = ?", data.STORE_PROVISION_STATUS_COMPLETED).
		Where("expires_at IS NULL OR expires_at > ?", time.Now().UTC()).
		Find(&provisions).Error
	if err != nil {
		return nil, err
	}
	return provisions, nil
}

func getStoreProductIDsByProductSizes(tx *gorm.DB, storeID uint, productSizeIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&data.StoreProduct{}).
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

	directIngredientIDs, err := getProductSizeDirectIngredientsByProductSizeIDs(tx, productSizeIDs)
	if err != nil {
		return nil, err
	}

	additiveIngredientIDs, err := getDefaultAdditiveIngredientsByProductSizeIDs(tx, productSizeIDs)
	if err != nil {
		return nil, err
	}

	return utils.UnionSlices(directIngredientIDs, additiveIngredientIDs), nil
}

func getAllProvisionIDsByProductSizes(tx *gorm.DB, productSizeIDs []uint) ([]uint, error) {
	if len(productSizeIDs) == 0 {
		return nil, nil
	}

	directProvisionIDs, err := getProductSizeDirectProvisionsByProductSizeIDs(tx, productSizeIDs)
	if err != nil {
		return nil, err
	}

	additiveProvisionIDs, err := getDefaultAdditiveProvisionsByProductSizeIDs(tx, productSizeIDs)
	if err != nil {
		return nil, err
	}

	return utils.UnionSlices(directProvisionIDs, additiveProvisionIDs), nil
}

func getAllIngredientIDsByAdditives(tx *gorm.DB, additiveIDs []uint) ([]uint, error) {
	if len(additiveIDs) == 0 {
		return nil, nil
	}

	var ingredientIDs []uint
	err := tx.Model(&data.AdditiveIngredient{}).
		Distinct("additive_ingredients.ingredient_id").
		Where("additive_ingredients.additive_id IN ?", additiveIDs).
		Pluck("additive_ingredients.ingredient_id", &ingredientIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch direct ingredientIDs: %w", err)
	}

	return ingredientIDs, nil
}

func getAllProvisionIDsByAdditives(tx *gorm.DB, additiveIDs []uint) ([]uint, error) {
	if len(additiveIDs) == 0 {
		return nil, nil
	}

	var provisionIDs []uint
	err := tx.Model(&data.AdditiveProvision{}).
		Distinct("additive_provisions.provision_id").
		Where("additive_provisions.additive_id IN ?", additiveIDs).
		Pluck("additive_provisions.provision_id", &provisionIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch direct provisionIDs: %w", err)
	}

	return provisionIDs, nil
}

func getProductSizeDirectProvisionsByProductSizeIDs(tx *gorm.DB, productSizeIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&data.ProductSizeProvision{}).
		Distinct("product_size_provisions.provision_id").
		Where("product_size_id IN ?", productSizeIDs).
		Pluck("provision_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provisions from product sizes: %w", err)
	}
	return ids, nil
}

func getDefaultAdditiveProvisionsByProductSizeIDs(tx *gorm.DB, productSizeIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&data.AdditiveProvision{}).
		Distinct("additive_provisions.provision_id").
		Joins("JOIN product_size_additives ON product_size_additives.additive_id = additive_provisions.additive_id").
		Where("product_size_additives.product_size_id IN ?", productSizeIDs).
		Where("product_size_additives.is_default = TRUE").
		Pluck("additive_provisions.provision_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provisions from default additives: %w", err)
	}
	return ids, nil
}

func getProductSizeDirectIngredientsByProductSizeIDs(tx *gorm.DB, productSizeIDs []uint) ([]uint, error) {
	var directIDs []uint
	err := tx.Model(&data.ProductSizeIngredient{}).
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
	err := tx.Model(&data.AdditiveIngredient{}).
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

// calculateFrozenInventory allows nil value for ingredientIDs if no need to filter ingredients
func calculateFrozenInventory(tx *gorm.DB, storeID uint, filter *types.FrozenInventoryFilter) (*types.FrozenInventory, error) {
	frozen := &types.FrozenInventory{
		Ingredients: make(map[uint]float64),
		Provisions:  make(map[uint]float64),
	}

	orders, err := loadActiveOrders(tx, storeID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to load active orders for store %d: %w", storeID, err)
	}

	for _, order := range orders {
		for _, sub := range order.Suborders {
			if !isSuborderActive(sub) {
				continue
			}
			accumulateProductUsage(frozen, sub, filter)
			accumulateAdditiveUsage(frozen, sub, filter)
		}
	}

	return frozen, nil
}

func loadActiveOrders(tx *gorm.DB, storeID uint, filter *types.FrozenInventoryFilter) ([]data.Order, error) {
	var orders []data.Order

	// Base query: always preload everything (no WHERE conditions on these)
	query := tx.
		Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.Ingredients.Ingredient").
		Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.AdditiveProvisions.Provision").
		Preload("Suborders.StoreProductSize.ProductSize.ProductSizeIngredients.Ingredient").
		Preload("Suborders.StoreProductSize.ProductSize.ProductSizeProvisions.Provision").
		Preload("Suborders.StoreProductSize.ProductSize.Additives.Additive.Ingredients.Ingredient").
		Preload("Suborders.StoreProductSize.ProductSize.Additives.Additive.AdditiveProvisions.Provision").
		Where("store_id = ?", storeID).
		Where("status IN ?", []data.OrderStatus{
			data.OrderStatusWaitingForPayment,
			data.OrderStatusPending,
			data.OrderStatusPreparing,
		})

	// If filter is provided, narrow certain relationships by ingredient/provision IDs
	if filter != nil {
		if len(filter.IngredientIDs) > 0 {
			query = query.
				Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.Ingredients",
					"ingredient_id IN ?", filter.IngredientIDs).
				Preload("Suborders.StoreProductSize.ProductSize.ProductSizeIngredients",
					"ingredient_id IN ?", filter.IngredientIDs).
				Preload("Suborders.StoreProductSize.ProductSize.Additives.Additive.Ingredients",
					"ingredient_id IN ?", filter.IngredientIDs)
		}

		if len(filter.ProvisionIDs) > 0 {
			query = query.
				Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.AdditiveProvisions",
					"provision_id IN ?", filter.ProvisionIDs).
				Preload("Suborders.StoreProductSize.ProductSize.ProductSizeProvisions",
					"provision_id IN ?", filter.ProvisionIDs).
				Preload("Suborders.StoreProductSize.ProductSize.Additives.Additive.AdditiveProvisions",
					"provision_id IN ?", filter.ProvisionIDs)
		}
	}

	err := query.Find(&orders).Error
	return orders, err
}

func isSuborderActive(sub data.Suborder) bool {
	return sub.Status == data.SubOrderStatusPending || sub.Status == data.SubOrderStatusPreparing
}

func accumulateProductUsage(frozen *types.FrozenInventory, sub data.Suborder, filter *types.FrozenInventoryFilter) {
	for _, ing := range sub.StoreProductSize.ProductSize.ProductSizeIngredients {
		if filter == nil || len(filter.IngredientIDs) == 0 || slices.Contains(filter.IngredientIDs, ing.IngredientID) {
			frozen.Ingredients[ing.IngredientID] += ing.Quantity
		}
	}
	for _, prov := range sub.StoreProductSize.ProductSize.ProductSizeProvisions {
		if filter == nil || len(filter.ProvisionIDs) == 0 || slices.Contains(filter.ProvisionIDs, prov.ProvisionID) {
			frozen.Provisions[prov.ProvisionID] += prov.Volume
		}
	}

	accumulateDefaultAdditivesUsage(frozen, sub.StoreProductSize.ProductSize, filter)
}

func accumulateDefaultAdditivesUsage(
	frozen *types.FrozenInventory,
	pSize data.ProductSize,
	filter *types.FrozenInventoryFilter,
) {
	// pSize.Additives = []ProductSizeAdditive where each has an Additive
	for _, psa := range pSize.Additives {
		if !psa.IsDefault {
			continue
		}
		// For each default additive, accumulate ingredients/provisions
		for _, ing := range psa.Additive.Ingredients {
			if filter == nil || len(filter.IngredientIDs) == 0 || slices.Contains(filter.IngredientIDs, ing.IngredientID) {
				frozen.Ingredients[ing.IngredientID] += ing.Quantity
			}
		}
		for _, prov := range psa.Additive.AdditiveProvisions {
			if filter == nil || len(filter.ProvisionIDs) == 0 || slices.Contains(filter.ProvisionIDs, prov.ProvisionID) {
				frozen.Provisions[prov.ProvisionID] += prov.Volume
			}
		}
	}
}

func accumulateAdditiveUsage(frozen *types.FrozenInventory, sub data.Suborder, filter *types.FrozenInventoryFilter) {
	for _, subAdd := range sub.SuborderAdditives {
		for _, ing := range subAdd.StoreAdditive.Additive.Ingredients {
			if filter == nil || len(filter.IngredientIDs) == 0 || slices.Contains(filter.IngredientIDs, ing.IngredientID) {
				frozen.Ingredients[ing.IngredientID] += ing.Quantity
			}
		}
		for _, prov := range subAdd.StoreAdditive.Additive.AdditiveProvisions {
			if filter == nil || len(filter.ProvisionIDs) == 0 || slices.Contains(filter.ProvisionIDs, prov.ProvisionID) {
				frozen.Provisions[prov.ProvisionID] += prov.Volume
			}
		}
	}
}
