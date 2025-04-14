package storeInventoryManagers

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeProvisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StoreInventoryManagerRepository interface {
	CheckStoreStocks(
		storeID uint,
		requiredIngredientQuantityMap map[uint]float64,
		frozenInventory *types.FrozenInventory,
	) error
	CheckStoreProvisions(
		storeID uint,
		requiredProvisionVolumeMap map[uint]float64,
		frozenInventory *types.FrozenInventory,
	) error

	DeductStoreInventoryByProductSize(storeID, productSizeID uint) (*types.DeductedStoreInventory, error)
	DeductStoreInventoryByAdditive(storeID, additiveID uint) (*types.DeductedStoreInventory, error)
	DeductStoreStocksByStoreProvision(storeProvision *data.StoreProvision) ([]data.StoreStock, error)

	RecalculateStoreAdditives(storeAdditiveIDs []uint, storeID uint, frozenInventory *types.FrozenInventory) error
	RecalculateStoreInventory(storeID uint, input *types.RecalculateInput) error
	CalculateFrozenInventory(storeID uint, filter *types.FrozenInventoryFilter) (*types.FrozenInventory, error)

	CloneWithTransaction(tx *gorm.DB) StoreInventoryManagerRepository
}

type storeInventoryManagerRepository struct {
	db *gorm.DB
}

func NewStoreInventoryManagerRepository(
	db *gorm.DB,
) StoreInventoryManagerRepository {
	return &storeInventoryManagerRepository{
		db: db,
	}
}

func (r *storeInventoryManagerRepository) CloneWithTransaction(tx *gorm.DB) StoreInventoryManagerRepository {
	return &storeInventoryManagerRepository{
		db: tx,
	}
}

func (r *storeInventoryManagerRepository) DeductStoreInventoryByProductSize(storeID, productSizeID uint) (*types.DeductedStoreInventory, error) {
	productSizeIngredients, err := r.getProductSizeIngredients(productSizeID)
	if err != nil {
		return nil, err
	}

	productSizeProvisions, err := r.getProductSizeProvisions(productSizeID)
	if err != nil {
		return nil, err
	}

	defaultProductSizeAdditiveIngredients, err := r.getProductSizeDefaultAdditiveIngredients(productSizeID)
	if err != nil {
		return nil, err
	}

	defaultProductSizeAdditiveProvisions, err := r.getProductSizeDefaultAdditiveProvisions(productSizeID)
	if err != nil {
		return nil, err
	}

	ingredientMap := make(map[uint]float64)
	for _, ing := range productSizeIngredients {
		ingredientMap[ing.IngredientID] += ing.Quantity
	}
	for _, ing := range defaultProductSizeAdditiveIngredients {
		ingredientMap[ing.IngredientID] += ing.Quantity
	}

	provisionMap := make(map[uint]float64)
	for _, prov := range productSizeProvisions {
		provisionMap[prov.ProvisionID] += prov.Volume
	}
	for _, prov := range defaultProductSizeAdditiveProvisions {
		provisionMap[prov.ProvisionID] += prov.Volume
	}

	deducted := &types.DeductedStoreInventory{}
	err = r.db.Transaction(func(tx *gorm.DB) error {
		for ingID, qty := range ingredientMap {
			updatedStock, err := deductStoreStock(tx, storeID, ingID, qty)
			if err != nil {
				return err
			}
			deducted.StoreStocks = append(deducted.StoreStocks, *updatedStock)
		}

		for provID, vol := range provisionMap {
			updatedProvisions, err := deductStoreProvisions(tx, storeID, provID, vol)
			if err != nil {
				return err
			}
			deducted.StoreProvisions = append(deducted.StoreProvisions, updatedProvisions...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return deducted, nil
}

func (r *storeInventoryManagerRepository) DeductStoreInventoryByAdditive(storeID, additiveID uint) (*types.DeductedStoreInventory, error) {
	additiveIngredients, err := r.getAdditiveIngredients(additiveID)
	if err != nil {
		return nil, err
	}

	additiveProvisions, err := r.getAdditiveProvisions(additiveID)
	if err != nil {
		return nil, err
	}

	deductedStoreInventory := &types.DeductedStoreInventory{}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, ingredient := range additiveIngredients {
			updatedStock, err := deductStoreStock(tx, storeID, ingredient.IngredientID, ingredient.Quantity)
			if err != nil {
				return err
			}
			deductedStoreInventory.StoreStocks = append(deductedStoreInventory.StoreStocks, *updatedStock)
		}

		for _, provision := range additiveProvisions {
			updatedProvisions, err := deductStoreProvisions(tx, storeID, provision.ProvisionID, provision.Volume)
			if err != nil {
				return err
			}
			deductedStoreInventory.StoreProvisions = append(deductedStoreInventory.StoreProvisions, updatedProvisions...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return deductedStoreInventory, nil
}

func (r *storeInventoryManagerRepository) DeductStoreStocksByStoreProvision(storeProvision *data.StoreProvision) ([]data.StoreStock, error) {
	if storeProvision == nil || storeProvision.StoreID == 0 || storeProvision.ID == 0 {
		return nil, fmt.Errorf("invalid input parameters")
	}

	storeProvisionIngredients, err := r.getStoreProvisionIngredients(storeProvision.ID)
	if err != nil {
		return nil, err
	}

	var deductedStocks []data.StoreStock

	err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, ingredient := range storeProvisionIngredients {
			updatedStock, err := deductStoreStock(tx, storeProvision.StoreID, ingredient.IngredientID, ingredient.Quantity)
			if err != nil {
				return err
			}
			deductedStocks = append(deductedStocks, *updatedStock)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return deductedStocks, nil
}

func (r *storeInventoryManagerRepository) RecalculateStoreInventory(storeID uint, input *types.RecalculateInput) error {
	logrus.Infof("=================RECALCULATION STARTS============================")
	start := time.Now()

	if storeID == 0 {
		return errors.New("failed to recalculate with invalid input parameters")
	}
	if input == nil {
		return nil
	}

	ctx := &recalculateContext{
		hasIngredients:  len(input.IngredientIDs) != 0,
		hasProvisions:   len(input.ProvisionIDs) != 0,
		hasProductSizes: len(input.ProductSizeIDs) != 0,
		hasAdditives:    len(input.AdditiveIDs) != 0,
	}

	if !ctx.hasIngredients && !ctx.hasProvisions && !ctx.hasProductSizes && !ctx.hasAdditives {
		return nil
	}

	if ctx.hasProductSizes {
		if err := r.gatherProductSizeData(ctx, storeID, input.ProductSizeIDs); err != nil {
			return err
		}
	}

	if ctx.hasAdditives {
		if err := r.gatherAdditiveData(ctx, storeID, input.AdditiveIDs); err != nil {
			return err
		}
	}

	ctx.totalIngredientIDs = utils.UnionSlices(
		input.IngredientIDs,
		ctx.productSizesIngredientIDs,
		ctx.additiveIngredientIDs,
	)
	ctx.totalProvisionIDs = utils.UnionSlices(
		input.ProvisionIDs,
		ctx.productSizesProvisionIDs,
		ctx.additiveProvisionIDs,
	)

	frozenInventory, err := r.buildFrozenInventory(ctx, storeID)
	if err != nil {
		return err
	}

	if err := r.gatherStoreProductAndAdditiveIDsFromResources(ctx, storeID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if len(ctx.storeProductIDsFromPS) > 0 || len(ctx.storeProductIDsFromIngredients) > 0 {
			err := recalculateStoreProducts(
				tx,
				utils.UnionSlices(
					ctx.storeProductIDsFromPS,
					ctx.storeProductIDsFromIngredients,
					ctx.storeProductIDsFromProvisions,
				),
				frozenInventory,
				storeID,
			)
			if err != nil {
				return err
			}
		}

		if len(ctx.storeAdditiveIDsFromAdditives) > 0 || len(ctx.storeAdditiveIDsFromIngredients) > 0 {
			err := recalculateStoreAdditives(
				tx,
				utils.UnionSlices(
					ctx.storeAdditiveIDsFromAdditives,
					ctx.storeAdditiveIDsFromIngredients,
					ctx.storeAdditiveIDsFromProvisions,
				),
				storeID,
				frozenInventory,
			)
			if err != nil {
				return err
			}
		}

		logrus.Infof("=====================Total time taken is: %v=====================", time.Since(start))
		return nil
	})
}

func (r *storeInventoryManagerRepository) gatherProductSizeData(rCtx *recalculateContext, storeID uint, productSizeIDs []uint) error {
	var err error

	rCtx.storeProductIDsFromPS, err = getStoreProductIDsByProductSizes(r.db, storeID, productSizeIDs)
	if err != nil {
		return err
	}

	rCtx.productSizesIngredientIDs, err = getAllIngredientIDsByProductSizes(r.db, productSizeIDs)
	if err != nil {
		return err
	}

	rCtx.productSizesProvisionIDs, err = getAllProvisionIDsByProductSizes(r.db, productSizeIDs)
	if err != nil {
		return err
	}

	return nil
}

func (r *storeInventoryManagerRepository) gatherAdditiveData(rCtx *recalculateContext, storeID uint, additiveIDs []uint) error {
	var err error

	rCtx.storeAdditiveIDsFromAdditives, err = getStoreAdditiveIDsByAdditives(r.db, storeID, additiveIDs)
	if err != nil {
		return err
	}

	rCtx.additiveIngredientIDs, err = getAllIngredientIDsByAdditives(r.db, additiveIDs)
	if err != nil {
		return err
	}

	rCtx.additiveProvisionIDs, err = getAllProvisionIDsByAdditives(r.db, additiveIDs)
	if err != nil {
		return err
	}
	return nil
}

func (r *storeInventoryManagerRepository) buildFrozenInventory(ctx *recalculateContext, storeID uint) (*types.FrozenInventory, error) {
	if len(ctx.totalIngredientIDs) == 0 && len(ctx.totalProvisionIDs) == 0 {
		return &types.FrozenInventory{
			Ingredients: make(map[uint]float64),
			Provisions:  make(map[uint]float64),
		}, nil
	}

	frozenInventoryFilter := &types.FrozenInventoryFilter{
		IngredientIDs: ctx.totalIngredientIDs,
		ProvisionIDs:  ctx.totalProvisionIDs,
	}
	logrus.Infof("Filter: %v", frozenInventoryFilter)

	frozen, err := calculateFrozenInventory(r.db, storeID, frozenInventoryFilter)
	if err != nil {
		return nil, err
	}
	logrus.Infof("FROZEN INVENTORY FETCHED(filtered): %v", frozen)

	return frozen, nil
}

func (r *storeInventoryManagerRepository) gatherStoreProductAndAdditiveIDsFromResources(ctx *recalculateContext, storeID uint) error {
	var err error

	if len(ctx.totalIngredientIDs) > 0 {
		ctx.storeProductIDsFromIngredients, err = getStoreProductIDsByIngredients(r.db, storeID, ctx.totalIngredientIDs)
		if err != nil {
			return err
		}

		ctx.storeAdditiveIDsFromIngredients, err = getStoreAdditiveIDsByIngredients(r.db, storeID, ctx.totalIngredientIDs)
		if err != nil {
			return err
		}
	}

	if len(ctx.totalProvisionIDs) > 0 {
		ctx.storeProductIDsFromProvisions, err = getStoreProductIDsByProvisions(r.db, storeID, ctx.totalProvisionIDs)
		if err != nil {
			return err
		}

		ctx.storeAdditiveIDsFromProvisions, err = getStoreAdditiveIDsByProvisions(r.db, storeID, ctx.totalProvisionIDs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *storeInventoryManagerRepository) RecalculateStoreAdditives(
	storeAdditiveIDs []uint,
	storeID uint,
	frozenInventory *types.FrozenInventory,
) error {
	return recalculateStoreAdditives(r.db, storeAdditiveIDs, storeID, frozenInventory)
}

func (r *storeInventoryManagerRepository) CalculateFrozenInventory(storeID uint, filter *types.FrozenInventoryFilter) (*types.FrozenInventory, error) {
	return calculateFrozenInventory(r.db, storeID, filter)
}

func (r *storeInventoryManagerRepository) CheckStoreStocks(
	storeID uint,
	requiredIngredientQuantityMap map[uint]float64,
	frozenInventory *types.FrozenInventory,
) error {
	if len(requiredIngredientQuantityMap) == 0 {
		return nil
	}

	ingredientIDs := make([]uint, 0, len(requiredIngredientQuantityMap))
	for id := range requiredIngredientQuantityMap {
		ingredientIDs = append(ingredientIDs, id)
	}

	stocks, err := getRelevantStoreStocks(r.db, storeID, ingredientIDs)
	if err != nil {
		return err
	}

	stockMap := make(map[uint]data.StoreStock)
	for _, stock := range stocks {
		stockMap[stock.IngredientID] = stock
	}

	for ingredientID, requiredQty := range requiredIngredientQuantityMap {
		stock, ok := stockMap[ingredientID]
		if !ok {
			return fmt.Errorf("%w: stock not found for ingredient ID %d", storeStocksTypes.ErrInsufficientStock, ingredientID)
		}

		frozen := frozenInventory.Ingredients[ingredientID]
		if stock.Quantity < frozen {
			return fmt.Errorf("%w: insufficient stock for ingredient ID %d: already pending %.2f, need %.2f, have %.2f",
				storeStocksTypes.ErrInsufficientStock, ingredientID, frozen, requiredQty, stock.Quantity)
		}

		effectiveAvailable := stock.Quantity - frozen
		if effectiveAvailable < requiredQty {
			return fmt.Errorf("%w: insufficient effective available for ingredient ID %d: need %.2f more",
				storeStocksTypes.ErrInsufficientStock, ingredientID, requiredQty)
		}

		frozenInventory.Ingredients[ingredientID] += requiredQty
	}

	return nil
}

func (r *storeInventoryManagerRepository) CheckStoreProvisions(
	storeID uint,
	requiredProvisionVolumeMap map[uint]float64,
	frozenInventory *types.FrozenInventory,
) error {
	if len(requiredProvisionVolumeMap) == 0 {
		return nil
	}

	var provisionIDs []uint
	for provID := range requiredProvisionVolumeMap {
		provisionIDs = append(provisionIDs, provID)
	}

	provisions, err := getRelevantStoreProvisions(r.db, storeID, provisionIDs)
	if err != nil {
		return fmt.Errorf("failed to load relevant storeProvisions: %w", err)
	}

	grouped := make(map[uint]float64)
	for _, sp := range provisions {
		grouped[sp.ProvisionID] += sp.Volume
	}

	for provID, reqVol := range requiredProvisionVolumeMap {
		totalVolume := grouped[provID]
		frozenUsed := frozenInventory.Provisions[provID]

		effectiveVolume := totalVolume - frozenUsed
		if effectiveVolume < reqVol {
			return fmt.Errorf(
				"%w: insufficient volume for provisionID=%d, needed=%.2f, have=%.2f (after frozen=%.2f)",
				storeProvisionsTypes.ErrInsufficientStoreProvision,
				provID, reqVol, totalVolume, frozenUsed,
			)
		}
		frozenInventory.Provisions[provID] += reqVol
	}

	return nil
}

func (r *storeInventoryManagerRepository) getProductSizeIngredients(productSizeID uint) ([]data.ProductSizeIngredient, error) {
	var productSizeIngredients []data.ProductSizeIngredient
	err := r.db.Preload("Ingredient").
		Preload("Ingredient.Unit").
		Where("product_size_id = ?", productSizeID).
		Find(&productSizeIngredients).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product size ingredients: %w", err)
	}
	return productSizeIngredients, nil
}

func (r *storeInventoryManagerRepository) getProductSizeDefaultAdditiveIngredients(productSizeID uint) ([]data.AdditiveIngredient, error) {
	var additiveIngredients []data.AdditiveIngredient
	err := r.db.Joins("JOIN product_size_additives ON product_size_additives.additive_id = additive_ingredients.additive_id").
		Preload("Ingredient.Unit").
		Where("product_size_additives.product_size_id = ? AND product_size_additives.is_default = TRUE", productSizeID).
		Find(&additiveIngredients).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product size additive ingredients: %w", err)
	}
	return additiveIngredients, nil
}

func (r *storeInventoryManagerRepository) getProductSizeDefaultAdditiveProvisions(productSizeID uint) ([]data.AdditiveProvision, error) {
	var additiveProvisions []data.AdditiveProvision
	err := r.db.Joins("JOIN product_size_additives ON product_size_additives.additive_id = additive_provisions.additive_id").
		Preload("Provision.Unit").
		Where("product_size_additives.product_size_id = ? AND product_size_additives.is_default = TRUE", productSizeID).
		Find(&additiveProvisions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product size additive provisions: %w", err)
	}
	return additiveProvisions, nil
}

func (r *storeInventoryManagerRepository) getProductSizeProvisions(productSizeID uint) ([]data.ProductSizeProvision, error) {
	var productSizeProvisions []data.ProductSizeProvision
	err := r.db.Preload("Provision").
		Preload("Provision.Unit").
		Where("product_size_id = ?", productSizeID).
		Find(&productSizeProvisions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product size provisions: %w", err)
	}
	return productSizeProvisions, nil
}

func (r *storeInventoryManagerRepository) getAdditiveIngredients(additiveID uint) ([]data.AdditiveIngredient, error) {
	var additiveIngredients []data.AdditiveIngredient
	err := r.db.Preload("Ingredient.Unit").
		Where("additive_id = ?", additiveID).
		Find(&additiveIngredients).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch additive ingredients: %w", err)
	}
	return additiveIngredients, nil
}

func (r *storeInventoryManagerRepository) getAdditiveProvisions(additiveID uint) ([]data.AdditiveProvision, error) {
	var additiveProvisions []data.AdditiveProvision
	err := r.db.Preload("Provision.Unit").
		Where("additive_id = ?", additiveID).
		Find(&additiveProvisions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch additive provisions: %w", err)
	}
	return additiveProvisions, nil
}

func (r *storeInventoryManagerRepository) getStoreProvisionIngredients(storeProvisionID uint) ([]data.StoreProvisionIngredient, error) {
	var storeProvisionIngredients []data.StoreProvisionIngredient
	err := r.db.Preload("Ingredient.Unit").
		Where("store_provision_id = ?", storeProvisionID).
		Find(&storeProvisionIngredients).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch store provision ingredients: %w", err)
	}
	return storeProvisionIngredients, nil
}
