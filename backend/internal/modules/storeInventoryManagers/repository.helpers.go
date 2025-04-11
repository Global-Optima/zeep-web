package storeInventoryManagers

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeProvisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"slices"
	"time"
)

func deductStoreStock(tx *gorm.DB, storeID, ingredientID uint, requiredQuantity float64) (*data.StoreStock, error) {
	var existingStock data.StoreStock
	err := tx.Preload("Ingredient").
		Preload("Ingredient.Unit").
		Where("store_id = ? AND ingredient_id = ?", storeID, ingredientID).
		First(&existingStock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("stock not found for ingredient ID %d", ingredientID)
		}
		return nil, fmt.Errorf("failed to fetch store warehouse stock: %w", err)
	}

	if existingStock.Quantity < requiredQuantity {
		return nil, fmt.Errorf("%w: insufficient stock for ingredient ID %d", storeStocksTypes.ErrInsufficientStock, ingredientID)
	}

	existingStock.Quantity -= requiredQuantity

	if err := tx.Save(&existingStock).Error; err != nil {
		return nil, fmt.Errorf("failed to save store warehouse stock for ingredient ID %d: %w", ingredientID, err)
	}

	return &existingStock, nil
}

func deductStoreProvisions(tx *gorm.DB, storeID, provisionID uint, requiredVolume float64) ([]data.StoreProvision, error) {
	if requiredVolume <= 0 {
		return nil, nil
	}

	var provisions []data.StoreProvision
	err := tx.Where("store_id = ? AND provision_id = ? AND status = ? AND (expires_at IS NULL OR expires_at > ?)",
		storeID, provisionID, data.STORE_PROVISION_STATUS_COMPLETED, time.Now().UTC()).
		Order("created_at ASC").
		Find(&provisions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load store provisions: %w", err)
	}

	var usedProvisions []data.StoreProvision
	remaining := requiredVolume

	for _, provision := range provisions {
		if remaining <= 0 {
			break
		}

		available := provision.Volume
		deduct := min(available, remaining)
		if deduct == 0 {
			continue
		}

		provision.Volume -= deduct
		remaining -= deduct

		if provision.Volume == 0 {
			provision.Status = data.STORE_PROVISION_STATUS_EMPTY
		}

		if err := tx.Save(&provision).Error; err != nil {
			return nil, fmt.Errorf("failed to update provision volume for ID %d: %w", provision.ID, err)
		}

		logrus.Infof("deducted: %v", deduct)

		usedProvisions = append(usedProvisions, provision)
	}

	if remaining > 0 {
		return nil, fmt.Errorf("%w: not enough provision volume for provision ID %d", storeProvisionsTypes.ErrInsufficientStoreProvision, provisionID)
	}

	return usedProvisions, nil
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
	outByIngredients, err := getOutByInventory(tx, storeProductIDs, storeID, frozenInventory)
	if err != nil {
		return nil, err
	}

	outByAdditives, err := getOutByDefaultAdditives(tx, storeProductIDs, storeID, frozenInventory)
	if err != nil {
		return nil, err
	}
	logrus.Info("outByInventory")

	return utils.UnionSlices(outByIngredients, outByAdditives), nil
}

func getOutByInventory(
	tx *gorm.DB,
	storeProductIDs []uint,
	storeID uint,
	frozen *types.FrozenInventory,
) ([]uint, error) {
	if len(storeProductIDs) == 0 {
		return nil, nil
	}
	if frozen == nil {
		return nil, fmt.Errorf("nil frozen inventory pointer fetched")
	}

	ingredientRows, err := getIngredientUsageRowsForStoreProducts(tx, storeProductIDs)
	if err != nil {
		return nil, err
	}

	provisionRows, err := getProvisionUsageRowsForStoreProducts(tx, storeProductIDs)
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

	// === LOAD PROVISION STOCK ===
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
		if stockMap[row.IngredientID] < row.RequiredQuantity {
			outSet[row.StoreProductOrAdditiveID] = struct{}{}
			logrus.Infof("outByIngredient %d", row.IngredientID)
		}
	}
	for _, row := range provisionRows {
		if provisionStockMap[row.ProvisionID] < row.RequiredVolume {
			logrus.Infof("provisionStock: %v / requiredVolume: %v", provisionStockMap[row.ProvisionID], row.RequiredVolume)
			outSet[row.StoreProductOrAdditiveID] = struct{}{}
			logrus.Infof("outByProvision %d", row.ProvisionID)
		}
	}

	var outIDs []uint
	for id := range outSet {
		outIDs = append(outIDs, id)
	}
	return outIDs, nil
}

func getIngredientUsageRowsForStoreProducts(tx *gorm.DB, storeProductIDs []uint) ([]types.IngredientUsageRow, error) {
	var ingredientRows []types.IngredientUsageRow
	err := tx.Model(&data.StoreProductSize{}).
		Select(`DISTINCT store_product_sizes.store_product_id AS store_product_or_additive_id,
                product_size_ingredients.ingredient_id,
                product_size_ingredients.quantity AS required_quantity`).
		Joins(`JOIN store_products ON store_products.id = store_product_sizes.store_product_id`).
		Joins(`JOIN product_size_ingredients ON product_size_ingredients.product_size_id = store_product_sizes.product_size_id`).
		Where(`store_products.id IN ?`, storeProductIDs).
		Scan(&ingredientRows).Error
	if err != nil {
		return nil, err
	}
	return ingredientRows, nil
}
func getProvisionUsageRowsForStoreProducts(tx *gorm.DB, storeProductIDs []uint) ([]types.ProvisionUsageRow, error) {
	var provisionRows []types.ProvisionUsageRow
	err := tx.Model(&data.StoreProductSize{}).
		Select(`DISTINCT store_product_sizes.store_product_id AS store_product_or_additive_id,
                product_size_provisions.provision_id,
                product_size_provisions.volume AS required_volume`).
		Joins(`JOIN store_products ON store_products.id = store_product_sizes.store_product_id`).
		Joins(`JOIN product_size_provisions ON product_size_provisions.product_size_id = store_product_sizes.product_size_id`).
		Where(`store_products.id IN ?`, storeProductIDs).
		Scan(&provisionRows).Error
	if err != nil {
		return nil, err
	}
	return provisionRows, nil
}

func getOutByDefaultAdditives(
	tx *gorm.DB,
	storeProductIDs []uint,
	storeID uint,
	frozen *types.FrozenInventory,
) ([]uint, error) {
	if len(storeProductIDs) == 0 {
		return nil, nil
	}
	if frozen == nil {
		return nil, fmt.Errorf("nil frozen inventory pointer fetched")
	}

	ingredientRows, err := getIngredientUsageRowsForDefaultAdditives(tx, storeProductIDs)
	if err != nil {
		return nil, err
	}

	provisionRows, err := getProvisionUsageRowsForDefaultAdditives(tx, storeProductIDs)
	if err != nil {
		return nil, err
	}

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

	// === LOAD PROVISION STOCK ===
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
		if stockMap[row.IngredientID] < row.RequiredQuantity {
			outSet[row.StoreProductOrAdditiveID] = struct{}{}
		}
	}
	for _, row := range provisionRows {
		if provisionStockMap[row.ProvisionID] < row.RequiredVolume {
			outSet[row.StoreProductOrAdditiveID] = struct{}{}
		}
	}

	var outIDs []uint
	for id := range outSet {
		outIDs = append(outIDs, id)
	}
	return outIDs, nil
}

func getIngredientUsageRowsForDefaultAdditives(tx *gorm.DB, storeProductIDs []uint) ([]types.IngredientUsageRow, error) {
	var ingredientRows []types.IngredientUsageRow
	err := tx.Model(&data.StoreProductSize{}).
		Select(`DISTINCT store_product_sizes.store_product_id AS store_product_or_additive_id,
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
	return ingredientRows, nil
}

func getProvisionUsageRowsForDefaultAdditives(tx *gorm.DB, storeProductIDs []uint) ([]types.ProvisionUsageRow, error) {
	var provisionRows []types.ProvisionUsageRow
	err := tx.Model(&data.StoreProductSize{}).
		Select(`DISTINCT store_product_sizes.store_product_id AS store_product_or_additive_id,
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
	return provisionRows, nil
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

	var ingredientRows []types.IngredientUsageRow
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

	var provisionRows []types.ProvisionUsageRow
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
			outSet[row.StoreProductOrAdditiveID] = struct{}{}
		}
	}

	for _, row := range provisionRows {
		available := provisionStockMap[row.ProvisionID]
		if available < row.RequiredVolume {
			outSet[row.StoreProductOrAdditiveID] = struct{}{}
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

	query := tx.Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.Ingredients.Ingredient").
		Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.AdditiveProvisions.Provision").
		Preload("Suborders.StoreProductSize.ProductSize.ProductSizeIngredients.Ingredient").
		Preload("Suborders.StoreProductSize.ProductSize.ProductSizeProvisions.Provision").
		Where("store_id = ?", storeID).
		Where("status IN ?", []data.OrderStatus{
			data.OrderStatusWaitingForPayment,
			data.OrderStatusPending,
			data.OrderStatusPreparing,
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
