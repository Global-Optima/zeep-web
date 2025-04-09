package storeInventoryManagers

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreInventoryManagerRepository interface {
	DeductStoreInventoryByProductSize(storeID, productSizeID uint) (*types.DeductedStoreInventory, error)
	DeductStoreInventoryByAdditive(storeID, additiveID uint) (*types.DeductedStoreInventory, error)

	RecalculateOutOfStock(storeID uint, input *types.RecalculateInput) error
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

	deductedStoreInventory := &types.DeductedStoreInventory{}
	var ingredientsToRecalculate []uint

	err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, ingredient := range productSizeIngredients {
			updatedStock, err := deductStoreStock(tx, storeID, ingredient.IngredientID, ingredient.Quantity)
			if err != nil {
				return err
			}
			deductedStoreInventory.StoreStocks = append(deductedStoreInventory.StoreStocks, *updatedStock)
			if updatedStock.Quantity <= updatedStock.LowStockThreshold {
				ingredientsToRecalculate = append(ingredientsToRecalculate, updatedStock.IngredientID)
			}
		}

		for _, provision := range productSizeProvisions {
			updatedProvisions, err := deductStoreProvisions(tx, storeID, provision.ProvisionID, provision.Volume)
			if err != nil {
				return err
			}
			//TODO should handle storeProvisions with 0 volume
			deductedStoreInventory.StoreProvisions = append(deductedStoreInventory.StoreProvisions, updatedProvisions...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return deductedStoreInventory, nil
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
	var ingredientsToRecalculate []uint

	err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, ingredient := range additiveIngredients {
			updatedStock, err := deductStoreStock(tx, storeID, ingredient.IngredientID, ingredient.Quantity)
			if err != nil {
				return err
			}
			deductedStoreInventory.StoreStocks = append(deductedStoreInventory.StoreStocks, *updatedStock)
			if updatedStock.Quantity <= updatedStock.LowStockThreshold {
				ingredientsToRecalculate = append(ingredientsToRecalculate, updatedStock.IngredientID)
			}
		}

		for _, provision := range additiveProvisions {
			//TODO delete or change status for updated storeProvisions
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

func (r *storeInventoryManagerRepository) RecalculateOutOfStock(storeID uint, input *types.RecalculateInput) error {
	if storeID == 0 {
		return errors.New("failed to recalculate with invalid input parameters")
	}
	if input == nil {
		return nil
	}

	var (
		hasIngredients  = len(input.IngredientIDs) != 0
		hasProvisions   = len(input.ProvisionIDs) != 0
		hasProductSizes = len(input.ProductSizeIDs) != 0
		hasAdditives    = len(input.AdditiveIDs) != 0
	)

	if !hasIngredients && !hasProvisions && !hasProductSizes && !hasAdditives {
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

	frozenInventory := &types.FrozenInventory{
		Ingredients: make(map[uint]float64),
		Provisions:  make(map[uint]float64),
	}

	if hasProductSizes {
		storeProductIDsFromPS, err = getStoreProductIDsByProductSizes(r.db, storeID, input.ProductSizeIDs)
		if err != nil {
			return err
		}

		productSizesIngredientIDs, err = getAllIngredientIDsByProductSizes(r.db, productSizesIngredientIDs)
		if err != nil {
			return err
		}

		productSizesProvisionIDs, err = getAllProvisionIDsByProductSizes(r.db, productSizesProvisionIDs)
		if err != nil {
			return err
		}

		if len(productSizesIngredientIDs) > 0 {
			input.IngredientIDs = utils.UnionSlices(input.IngredientIDs, productSizesIngredientIDs)
		}
	}

	if hasAdditives {
		storeAdditiveIDsFromAdditives, err = getStoreAdditiveIDsByAdditives(r.db, storeID, input.AdditiveIDs)
		if err != nil {
			return err
		}
	}

	if hasIngredients || hasProvisions {
		frozenInventoryFilter := &types.FrozenInventoryFilter{
			IngredientIDs: input.IngredientIDs,
			ProvisionIDs:  input.ProvisionIDs,
		}

		frozenInventory, err = CalculateFrozenInventory(r.db, storeID, frozenInventoryFilter)
		if err != nil {
			return err
		}
	}

	if hasIngredients {
		storeProductIDsFromIngredients, err = getStoreProductIDsByIngredients(r.db, storeID, input.IngredientIDs)
		if err != nil {
			return err
		}

		storeAdditiveIDsFromIngredients, err = getStoreAdditiveIDsByIngredients(r.db, storeID, input.IngredientIDs)
		if err != nil {
			return err
		}
	}

	if hasProvisions {
		storeProductIDsFromProvisions, err = getStoreProductIDsByProvisions(r.db, storeID, input.ProvisionIDs)
		if err != nil {
			return err
		}

		storeAdditiveIDsFromProvisions, err = getStoreAdditiveIDsByProvisions(r.db, storeID, input.ProvisionIDs)
		if err != nil {
			return err
		}
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
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
	})
}

func (r *storeInventoryManagerRepository) CalculateFrozenInventory(storeID uint, filter *types.FrozenInventoryFilter) (*types.FrozenInventory, error) {
	return CalculateFrozenInventory(r.db, storeID, filter)
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

func (r *storeInventoryManagerRepository) getAdditiveProvisions(provisionID uint) ([]data.AdditiveProvision, error) {
	var additiveProvisions []data.AdditiveProvision
	err := r.db.Preload("Provision.Unit").
		Where("provision_id = ?", provisionID).
		Find(&additiveProvisions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch additive provisions: %w", err)
	}
	return additiveProvisions, nil
}
