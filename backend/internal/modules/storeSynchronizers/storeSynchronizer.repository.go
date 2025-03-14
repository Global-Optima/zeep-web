package storeSynchronizers

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeSynchronizers/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type StoreSynchronizeRepository interface {
	IsSynchronizedStore(storeID uint) (bool, error)
	GetNotSynchronizedAdditivesIDs(storeID uint, lastSync time.Time) ([]uint, error)
	GetNotSynchronizedProductSizesIngredients(storeID uint, lastSync time.Time) ([]uint, error)
	GetNotSynchronizedProductSizesAdditivesIDs(storeID uint, lastSync time.Time) ([]uint, error)
}

type storeSynchronizeRepository struct {
	db                *gorm.DB
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository
	storeStockRepo    storeStocks.StoreStockRepository
}

func NewStoreSynchronizeRepository(db *gorm.DB, storeAdditiveRepo storeAdditives.StoreAdditiveRepository, storeStockRepo storeStocks.StoreStockRepository) StoreSynchronizeRepository {
	return &storeSynchronizeRepository{db: db}
}

func (r *storeSynchronizeRepository) CloneWithTransaction(tx *gorm.DB) StoreSynchronizeRepository {
	return &storeSynchronizeRepository{
		db: tx,
	}
}

func (r *storeSynchronizeRepository) IsSynchronizedStore(storeID uint) (bool, error) {
	// 1. Retrieve store's LastInventorySyncAt timestamp
	var lastSyncTime time.Time
	err := r.db.Model(&data.Store{}).
		Select("last_inventory_sync_at").
		Where("id = ?", storeID).
		Scan(&lastSyncTime).Error
	if err != nil {
		return false, fmt.Errorf("failed to fetch store sync time: %w", err)
	}

	// 2. Fetch unsynchronized product size additives IDs
	var additiveIDs []uint
	err = r.db.Model(&data.ProductSizeAdditive{}).
		Joins("JOIN product_sizes ON product_sizes.id = product_size_additives.product_size_id").
		Joins("JOIN store_product_sizes ON product_sizes.id = store_product_sizes.product_size_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Where("store_products.store_id = ?", storeID).
		Where("product_size_additives.created_at > ?", lastSyncTime).
		Pluck("product_size_additives.id", &additiveIDs).Error
	if err != nil {
		return false, fmt.Errorf("failed to fetch unsynchronized product size additives: %w", err)
	}

	// 3. Fetch unsynchronized product size ingredients IDs
	var ingredientIDs []uint
	err = r.db.Model(&data.ProductSizeIngredient{}).
		Joins("JOIN product_sizes ON product_sizes.id = product_size_ingredients.product_size_id").
		Joins("JOIN store_product_sizes ON product_sizes.id = store_product_sizes.product_size_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Where("store_products.store_id = ?", storeID).
		Where("product_size_ingredients.created_at > ?", lastSyncTime).
		Pluck("product_size_ingredients.id", &ingredientIDs).Error
	if err != nil {
		return false, fmt.Errorf("failed to fetch unsynchronized product size ingredients: %w", err)
	}

	// 4. Fetch unsynchronized additive ingredients IDs
	var additiveIngredientIDs []uint
	err = r.db.Model(&data.AdditiveIngredient{}).
		Joins("JOIN additives ON additives.id = additive_ingredients.additive_id").
		Joins("JOIN store_additives ON additives.id = store_additives.additive_id").
		Where("store_additives.store_id = ?", storeID).
		Where("additive_ingredients.created_at > ?", lastSyncTime).
		Pluck("additive_ingredients.id", &additiveIngredientIDs).Error
	if err != nil {
		return false, fmt.Errorf("failed to fetch unsynchronized additive ingredients: %w", err)
	}

	// 5. Use the Filter function to determine missing IDs
	missingAdditiveIDs, err := r.storeAdditiveRepo.FilterMissingStoreAdditiveIDs(storeID, additiveIDs)
	if err != nil {
		return false, fmt.Errorf("failed to fetch missing additive ingredients: %w", err)
	}
	if len(missingAdditiveIDs) > 0 {
		return false, nil
	}

	missingIngredientIDs, err := r.storeStockRepo.FilterMissingIngredientsIDs(storeID, utils.MergeDistinct(ingredientIDs, additiveIngredientIDs))
	if err != nil {
		return false, fmt.Errorf("failed to fetch missing ingredients: %w", err)
	}
	if len(missingIngredientIDs) > 0 {
		return false, nil
	}

	go func() {
		storeRepo
	}()
	// Store is synchronized if no missing IDs
	return true, nil
}

func (r *storeSynchronizeRepository) GetNotSynchronizedAdditivesIDs(storeID uint, lastSync time.Time) ([]uint, error) {
	var notSynchronizedProductSizesAdditivesIDs []uint
	err := r.db.Model(&data.ProductSizeAdditive{}).
		Distinct("product_size_additives.additive_id").
		Joins("JOIN store_additives ON store_additives.id = product_size_additives.product_size_id").
		Joins("JOIN additives ON additives.id = product_size_additives.additive_id").
		Joins("JOIN product_sizes ON product_sizes.id = product_size_additives.product_size_id").
		Joins("JOIN store_product_sizes ON product_sizes.id = store_product_sizes.product_size_id").
		Where("store_additives.store_id = ?", storeID).
		Where("additives.created_at > ?", lastSync).
		Pluck("product_size_additives.additive_id", &notSynchronizedProductSizesAdditivesIDs).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrNotSynchronizedAdditivesNotFound
		}
		return nil, fmt.Errorf("failed to fetch product sizes: %w", err)
	}

	return notSynchronizedProductSizesAdditivesIDs, nil
}

func (r *storeSynchronizeRepository) GetNotSynchronizedProductSizesIngredients(storeID uint, lastSync time.Time) ([]uint, error) {
	var notSynchronizedIngredients []uint

	err := r.db.Model(&data.ProductSizeIngredient{}).
		Distinct("product_size_ingredients.ingredient_id").
		Joins("JOIN product_sizes ON product_sizes.id = product_size_ingredients.product_size_id").
		Joins("JOIN ingredients ON ingredients.id = product_size_ingredients.ingredient_id").
		Joins("JOIN store_product_sizes ON product_sizes.id = store_product_sizes.product_size_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Where("store_products.store_id = ?", storeID).
		Where("ingredients.created_at > ?", lastSync).
		Pluck("product_size_ingredients.ingredient_id", &notSynchronizedIngredients).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrNotSynchronizedProductSizeIngredientsNotFound
		}
		return nil, err
	}

	return notSynchronizedIngredients, nil
}

func (r *storeSynchronizeRepository) GetNotSynchronizedProductSizesAdditivesIDs(storeID uint, lastSync time.Time) ([]uint, error) {
	var notSynchronizedProductSizesAdditivesIDs []uint
	err := r.db.Model(&data.ProductSizeAdditive{}).
		Distinct("product_size_additives.additive_id").
		Joins("JOIN product_sizes ON product_sizes.id = product_size_additives.product_size_id").
		Joins("JOIN additives ON additives.id = product_size_additives.additive_id").
		Joins("JOIN store_product_sizes ON product_sizes.id = store_product_sizes.product_size_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Where("store_products.store_id = ?", storeID).
		Where("additives.created_at > ?", lastSync).
		Pluck("product_size_additives.additive_id", &notSynchronizedProductSizesAdditivesIDs).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrNotSynchronizedProductSizeAdditivesNotFound
		}
		return nil, fmt.Errorf("failed to fetch product sizes: %w", err)
	}

	return notSynchronizedProductSizesAdditivesIDs, nil
}
