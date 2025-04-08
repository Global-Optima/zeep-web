package data

import (
	"fmt"
	"github.com/pkg/errors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type RecalculateInput struct {
	IngredientIDs  []uint
	ProvisionIDs   []uint
	ProductSizeIDs []uint
	AdditiveIDs    []uint
}

type FrozenInventory struct {
	Ingredients map[uint]float64
	Provisions  map[uint]float64
}

type ingredientUsageRow struct {
	EntityID         uint
	IngredientID     uint
	RequiredQuantity float64
}

type provisionUsageRow struct {
	EntityID       uint
	ProvisionID    uint
	RequiredVolume float64
}

// RecalculateOutOfStock allows nil values for ingredientIDs and productSizeIDs if no need to check by any of them
// TODO refactor and use variables to store states of checks like hasProductSizeProvisions := len(productSizesProvisionIDs) > 0
func RecalculateOutOfStock(tx *gorm.DB, storeID uint, input *RecalculateInput) error {
	if tx == nil || storeID == 0 {
		return errors.New("failed to recalculate with invalid input parameters")
	}

	if input == nil || len(input.IngredientIDs) == 0 && len(input.ProductSizeIDs) == 0 && len(input.AdditiveIDs) == 0 && len(input.ProvisionIDs) == 0 {
		return nil
	}

	var productSizesIngredientIDs,
		productSizesProvisionIDs,
		storeProductIDsFromPS,
		storeProductIDsFromIngredients,
		storeProductIDsFromProvisions,
		storeAdditiveIDsFromAdditives,
		storeAdditiveIDsFromIngredients,
		storeAdditiveIDsFromProvisions []uint
	var err error

	frozenInventory := &FrozenInventory{
		Ingredients: make(map[uint]float64),
		Provisions:  make(map[uint]float64),
	}

	if len(input.ProductSizeIDs) > 0 {
		storeProductIDsFromPS, err = getStoreProductIDsByProductSizes(tx, storeID, input.ProductSizeIDs)
		if err != nil {
			return err
		}

		productSizesIngredientIDs, err = getAllIngredientIDsByProductSizes(tx, productSizesIngredientIDs)
		if err != nil {
			return err
		}

		productSizesProvisionIDs, err = getAllProvisionIDsByProductSizes(tx, productSizesProvisionIDs)
		if err != nil {
			return err
		}

		if len(productSizesIngredientIDs) > 0 {
			input.IngredientIDs = utils.UnionSlices(input.IngredientIDs, productSizesIngredientIDs)
		}
	}

	if len(input.AdditiveIDs) > 0 {
		storeAdditiveIDsFromAdditives, err = getStoreAdditiveIDsByAdditives(tx, storeID, input.AdditiveIDs)
		if err != nil {
			return err
		}
	}

	if len(input.IngredientIDs) > 0 || len(input.ProvisionIDs) > 0 {
		frozenInventoryFilter := &FrozenInventoryFilter{
			IngredientIDs: input.IngredientIDs,
			ProvisionIDs:  input.ProvisionIDs,
		}

		frozenInventory, err = CalculateFrozenInventory(tx, storeID, frozenInventoryFilter)
		if err != nil {
			return err
		}
	}

	if len(input.IngredientIDs) > 0 {
		storeProductIDsFromIngredients, err = getStoreProductIDsByIngredients(tx, storeID, input.IngredientIDs)
		if err != nil {
			return err
		}

		storeAdditiveIDsFromIngredients, err = getStoreAdditiveIDsByIngredients(tx, storeID, input.IngredientIDs)
		if err != nil {
			return err
		}
	}

	if len(input.ProvisionIDs) > 0 {
		storeProductIDsFromProvisions, err = getStoreProductIDsByProvisions(tx, storeID, input.ProvisionIDs)
		if err != nil {
			return err
		}

		storeAdditiveIDsFromProvisions, err = getStoreAdditiveIDsByProvisions(tx, storeID, input.ProvisionIDs)
		if err != nil {
			return err
		}
	}

	if len(storeProductIDsFromPS) > 0 || len(storeProductIDsFromIngredients) > 0 {
		if err := RecalculateStoreProducts(
			tx,
			utils.UnionSlices(storeProductIDsFromPS, storeProductIDsFromIngredients, storeProductIDsFromProvisions),
			frozenInventory,
			storeID,
		); err != nil {
			return err
		}
	}

	if len(storeAdditiveIDsFromAdditives) > 0 || len(storeAdditiveIDsFromIngredients) > 0 {
		if err := RecalculateStoreAdditives(
			tx,
			utils.UnionSlices(storeAdditiveIDsFromAdditives, storeAdditiveIDsFromIngredients, storeAdditiveIDsFromProvisions),
			storeID,
			frozenInventory,
		); err != nil {
			return err
		}
	}

	return nil
}

func RecalculateStoreProducts(
	tx *gorm.DB,
	storeProductIDs []uint,
	frozenInventory *FrozenInventory,
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

func RecalculateStoreAdditives(
	tx *gorm.DB,
	storeAdditiveIDs []uint,
	storeID uint,
	frozenInventory *FrozenInventory,
) error {
	if len(storeAdditiveIDs) == 0 {
		return nil
	}

	outOfStockIDs, err := getOutOfStockStoreAdditiveIDs(tx, storeAdditiveIDs, storeID, frozenInventory)
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
	err := tx.Model(&StoreProductSize{}).
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
	err := tx.Model(&StoreProductSize{}).
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

func getStoreAdditiveIDsByAdditives(tx *gorm.DB, storeID uint, additiveIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&StoreAdditive{}).
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
	err := tx.Model(&StoreAdditive{}).
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
	err := tx.Model(&StoreAdditive{}).
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
	frozenInventory *FrozenInventory,
) ([]uint, error) {
	outByIngredients, err := getOutByInventory(tx, storeProductIDs, storeID, frozenInventory)
	if err != nil {
		return nil, err
	}

	outByAdditives, err := getOutByDefaultAdditives(tx, storeProductIDs, storeID, frozenInventory)
	if err != nil {
		return nil, err
	}

	return utils.UnionSlices(outByIngredients, outByAdditives), nil
}

func getOutByInventory(
	tx *gorm.DB,
	storeProductIDs []uint,
	storeID uint,
	frozen *FrozenInventory,
) ([]uint, error) {
	if len(storeProductIDs) == 0 {
		return nil, nil
	}
	if frozen == nil {
		return nil, fmt.Errorf("nil frozen inventory pointer fetched")
	}

	// === INGREDIENTS ===
	var ingredientRows []ingredientUsageRow
	err := tx.Model(&StoreProductSize{}).
		Select(`DISTINCT store_product_sizes.store_product_id AS entity_id,
                product_size_ingredients.ingredient_id,
                product_size_ingredients.quantity AS required_quantity`).
		Joins(`JOIN store_products ON store_products.id = store_product_sizes.store_product_id`).
		Joins(`JOIN product_size_ingredients ON product_size_ingredients.product_size_id = store_product_sizes.product_size_id`).
		Where(`store_products.id IN ?`, storeProductIDs).
		Scan(&ingredientRows).Error
	if err != nil {
		return nil, err
	}

	// === PROVISIONS ===
	var provisionRows []provisionUsageRow
	err = tx.Model(&StoreProductSize{}).
		Select(`DISTINCT store_product_sizes.store_product_id AS entity_id,
                product_size_provisions.provision_id,
                product_size_provisions.volume AS required_volume`).
		Joins(`JOIN store_products ON store_products.id = store_product_sizes.store_product_id`).
		Joins(`JOIN product_size_provisions ON product_size_provisions.product_size_id = store_product_sizes.product_size_id`).
		Where(`store_products.id IN ?`, storeProductIDs).
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

	var stocks []StoreStock
	err = tx.Model(&StoreStock{}).
		Where("store_id = ?", storeID).
		Where("ingredient_id IN ?", neededIngredientIDs).
		Find(&stocks).Error
	if err != nil {
		return nil, err
	}

	stockMap := make(map[uint]float64)
	for _, s := range stocks {
		stockMap[s.IngredientID] = s.Quantity - frozen.Ingredients[s.IngredientID]
	}

	// === LOAD PROVISION STOCK ===
	provisionSet := make(map[uint]struct{})
	for _, row := range provisionRows {
		provisionSet[row.ProvisionID] = struct{}{}
	}
	var neededProvisionIDs []uint
	for id := range provisionSet {
		neededProvisionIDs = append(neededProvisionIDs, id)
	}

	var provisions []StoreProvision
	err = tx.Model(&StoreProvision{}).
		Where("store_id = ?", storeID).
		Where("provision_id IN ?", neededProvisionIDs).
		Where("status = ?", STORE_PROVISION_STATUS_COMPLETED).
		Where("expires_at IS NULL OR expires_at > ?", time.Now().UTC()).
		Find(&provisions).Error
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
		if stockMap[row.IngredientID] < row.RequiredQuantity {
			outSet[row.EntityID] = struct{}{}
		}
	}
	for _, row := range provisionRows {
		if provisionStockMap[row.ProvisionID] < row.RequiredVolume {
			outSet[row.EntityID] = struct{}{}
		}
	}

	var outIDs []uint
	for id := range outSet {
		outIDs = append(outIDs, id)
	}
	return outIDs, nil
}

func getOutByDefaultAdditives(
	tx *gorm.DB,
	storeProductIDs []uint,
	storeID uint,
	frozen *FrozenInventory,
) ([]uint, error) {
	if len(storeProductIDs) == 0 {
		return nil, nil
	}
	if frozen == nil {
		return nil, fmt.Errorf("nil frozen inventory pointer fetched")
	}

	// === INGREDIENTS ===
	var ingredientRows []ingredientUsageRow
	err := tx.Model(&StoreProductSize{}).
		Select(`DISTINCT store_product_sizes.store_product_id AS entity_id,
                additive_ingredients.ingredient_id,
                additive_ingredients.quantity AS required_quantity`).
		Joins(`JOIN store_products ON store_products.id = store_product_sizes.store_product_id`).
		Joins(`JOIN product_size_additives
               ON product_size_additives.product_size_id = store_product_sizes.product_size_id
               AND product_size_additives.is_default = TRUE`).
		Joins(`JOIN additive_ingredients
               ON additive_ingredients.additive_id = product_size_additives.additive_id`).
		Where(`store_products.id IN ?`, storeProductIDs).
		Scan(&ingredientRows).Error
	if err != nil {
		return nil, err
	}

	// === PROVISIONS ===
	var provisionRows []provisionUsageRow
	err = tx.Model(&StoreProductSize{}).
		Select(`DISTINCT store_product_sizes.store_product_id AS entity_id,
                additive_provisions.provision_id,
                additive_provisions.volume AS required_volume`).
		Joins(`JOIN store_products ON store_products.id = store_product_sizes.store_product_id`).
		Joins(`JOIN product_size_additives ON product_size_additives.product_size_id = store_product_sizes.product_size_id AND product_size_additives.is_default = TRUE`).
		Joins(`JOIN additive_provisions ON additive_provisions.additive_id = product_size_additives.additive_id`).
		Where(`store_products.id IN ?`, storeProductIDs).
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

	var stocks []StoreStock
	err = tx.Model(&StoreStock{}).
		Where("store_id = ?", storeID).
		Where("ingredient_id IN ?", neededIngredientIDs).
		Find(&stocks).Error
	if err != nil {
		return nil, err
	}

	stockMap := make(map[uint]float64)
	for _, s := range stocks {
		stockMap[s.IngredientID] = s.Quantity - frozen.Ingredients[s.IngredientID]
	}

	// === LOAD PROVISION STOCK ===
	provisionSet := make(map[uint]struct{})
	for _, row := range provisionRows {
		provisionSet[row.ProvisionID] = struct{}{}
	}
	var neededProvisionIDs []uint
	for id := range provisionSet {
		neededProvisionIDs = append(neededProvisionIDs, id)
	}

	var provisions []StoreProvision
	err = tx.Model(&StoreProvision{}).
		Where("store_id = ?", storeID).
		Where("provision_id IN ?", neededProvisionIDs).
		Where("status = ?", STORE_PROVISION_STATUS_COMPLETED).
		Where("expires_at IS NULL OR expires_at > ?", time.Now().UTC()).
		Find(&provisions).Error
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
		if stockMap[row.IngredientID] < row.RequiredQuantity {
			outSet[row.EntityID] = struct{}{}
		}
	}
	for _, row := range provisionRows {
		if provisionStockMap[row.ProvisionID] < row.RequiredVolume {
			outSet[row.EntityID] = struct{}{}
		}
	}

	var outIDs []uint
	for id := range outSet {
		outIDs = append(outIDs, id)
	}
	return outIDs, nil
}

func getOutOfStockStoreAdditiveIDs(
	tx *gorm.DB,
	storeAdditiveIDs []uint,
	storeID uint,
	frozen *FrozenInventory,
) ([]uint, error) {
	if len(storeAdditiveIDs) == 0 {
		return nil, nil
	}
	if frozen == nil {
		return nil, fmt.Errorf("nil frozen inventory pointer fetched")
	}

	// === INGREDIENTS ===
	var ingredientRows []ingredientUsageRow
	err := tx.Model(&StoreAdditive{}).
		Select(`DISTINCT store_additives.id AS entity_id,
                additive_ingredients.ingredient_id,
                additive_ingredients.quantity AS required_quantity`).
		Joins(`JOIN additive_ingredients
               ON additive_ingredients.additive_id = store_additives.additive_id`).
		Where(`store_additives.id IN ?`, storeAdditiveIDs).
		Scan(&ingredientRows).Error
	if err != nil {
		return nil, err
	}

	// === PROVISIONS ===
	var provisionRows []provisionUsageRow
	err = tx.Model(&StoreAdditive{}).
		Select(`DISTINCT store_additives.id AS entity_id,
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

	var stocks []StoreStock
	err = tx.Model(&StoreStock{}).
		Where("store_id = ?", storeID).
		Where("ingredient_id IN ?", neededIngredientIDs).
		Find(&stocks).Error
	if err != nil {
		return nil, err
	}

	stockMap := make(map[uint]float64)
	for _, s := range stocks {
		stockMap[s.IngredientID] = s.Quantity - frozen.Ingredients[s.IngredientID]
	}

	// === LOAD PROVISION STOCK ===
	provisionSet := make(map[uint]struct{})
	for _, row := range provisionRows {
		provisionSet[row.ProvisionID] = struct{}{}
	}
	var neededProvisionIDs []uint
	for id := range provisionSet {
		neededProvisionIDs = append(neededProvisionIDs, id)
	}

	var provisions []StoreProvision
	err = tx.Model(&StoreProvision{}).
		Where("store_id = ?", storeID).
		Where("provision_id IN ?", neededProvisionIDs).
		Where("status = ?", STORE_PROVISION_STATUS_COMPLETED).
		Where("expires_at IS NULL OR expires_at > ?", time.Now().UTC()).
		Find(&provisions).Error
	if err != nil {
		return nil, err
	}

	provisionStockMap := make(map[uint]float64)
	for _, p := range provisions {
		provisionStockMap[p.ProvisionID] += p.Volume
	}
	// Вычитаем замороженные объёмы
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
			outSet[row.EntityID] = struct{}{}
		}
	}

	for _, row := range provisionRows {
		available := provisionStockMap[row.ProvisionID]
		if available < row.RequiredVolume {
			outSet[row.EntityID] = struct{}{}
		}
	}

	var outIDs []uint
	for id := range outSet {
		outIDs = append(outIDs, id)
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

func getAllProvisionIDsByProductSizes(tx *gorm.DB, productSizeIDs []uint) ([]uint, error) {
	if len(productSizeIDs) == 0 {
		return nil, nil
	}

	directIDs, err := getProductSizeDirectProvisionsByProductSizeIDs(tx, productSizeIDs)
	if err != nil {
		return nil, err
	}

	additiveIDs, err := getDefaultAdditiveProvisionsByProductSizeIDs(tx, productSizeIDs)
	if err != nil {
		return nil, err
	}

	return utils.UnionSlices(directIDs, additiveIDs), nil
}

func getProductSizeDirectProvisionsByProductSizeIDs(tx *gorm.DB, productSizeIDs []uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&ProductSizeProvision{}).
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
	err := tx.Model(&AdditiveProvision{}).
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

type FrozenInventoryFilter struct {
	IngredientIDs []uint
	ProvisionIDs  []uint
}

// CalculateFrozenInventory allows nil value for ingredientIDs if no need to filter ingredients
func CalculateFrozenInventory(tx *gorm.DB, storeID uint, filter *FrozenInventoryFilter) (*FrozenInventory, error) {
	frozen := &FrozenInventory{
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

func loadActiveOrders(tx *gorm.DB, storeID uint, filter *FrozenInventoryFilter) ([]Order, error) {
	var orders []Order

	query := tx.Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.Ingredients.Ingredient").
		Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.AdditiveProvisions.Provision").
		Preload("Suborders.StoreProductSize.ProductSize.ProductSizeIngredients.Ingredient").
		Preload("Suborders.StoreProductSize.ProductSize.ProductSizeProvisions.Provision").
		Where("store_id = ?", storeID).
		Where("status IN ?", []OrderStatus{
			OrderStatusWaitingForPayment,
			OrderStatusPending,
			OrderStatusPreparing,
		})

	if filter != nil {
		if len(filter.IngredientIDs) > 0 {
			query = query.
				Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.Ingredients", "ingredient_id IN ?", filter.IngredientIDs).
				Preload("Suborders.StoreProductSize.ProductSize.ProductSizeIngredients", "ingredient_id IN ?", filter.IngredientIDs)
		}
		if len(filter.ProvisionIDs) > 0 {
			query = query.
				Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.AdditiveProvisions", "provision_id IN ?", filter.ProvisionIDs).
				Preload("Suborders.StoreProductSize.ProductSize.ProductSizeProvisions", "provision_id IN ?", filter.ProvisionIDs)
		}
	}

	err := query.Find(&orders).Error
	return orders, err
}

func isSuborderActive(sub Suborder) bool {
	return sub.Status == SubOrderStatusPending || sub.Status == SubOrderStatusPreparing
}

func accumulateProductUsage(frozen *FrozenInventory, sub Suborder, filter *FrozenInventoryFilter) {
	for _, ing := range sub.StoreProductSize.ProductSize.ProductSizeIngredients {
		if filter == nil || len(filter.IngredientIDs) == 0 || contains(filter.IngredientIDs, ing.IngredientID) {
			frozen.Ingredients[ing.IngredientID] += ing.Quantity
		}
	}
	for _, prov := range sub.StoreProductSize.ProductSize.ProductSizeProvisions {
		if filter == nil || len(filter.ProvisionIDs) == 0 || contains(filter.ProvisionIDs, prov.ProvisionID) {
			frozen.Provisions[prov.ProvisionID] += prov.Volume
		}
	}
}

func accumulateAdditiveUsage(frozen *FrozenInventory, sub Suborder, filter *FrozenInventoryFilter) {
	for _, subAdd := range sub.SuborderAdditives {
		for _, ing := range subAdd.StoreAdditive.Additive.Ingredients {
			if filter == nil || len(filter.IngredientIDs) == 0 || contains(filter.IngredientIDs, ing.IngredientID) {
				frozen.Ingredients[ing.IngredientID] += ing.Quantity
			}
		}
		for _, prov := range subAdd.StoreAdditive.Additive.AdditiveProvisions {
			if filter == nil || len(filter.ProvisionIDs) == 0 || contains(filter.ProvisionIDs, prov.ProvisionID) {
				frozen.Provisions[prov.ProvisionID] += prov.Volume
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
