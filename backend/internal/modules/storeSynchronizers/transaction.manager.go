package storeSynchronizers

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeSynchronizers/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	storesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type TransactionManager interface {
	SynchronizeStoreInventory(storeID uint) error
}

type transactionManager struct {
	db                *gorm.DB
	repo              StoreSynchronizeRepository
	storeRepo         stores.StoreRepository
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository
	storeStockRepo    storeStocks.StoreStockRepository
}

func NewTransactionManager(db *gorm.DB, repo StoreSynchronizeRepository, storeRepo stores.StoreRepository, storeAdditiveRepo storeAdditives.StoreAdditiveRepository, storeStockRepo storeStocks.StoreStockRepository) TransactionManager {
	return &transactionManager{
		db:                db,
		repo:              repo,
		storeRepo:         storeRepo,
		storeAdditiveRepo: storeAdditiveRepo,
		storeStockRepo:    storeStockRepo,
	}
}

func (m *transactionManager) IsSynchronizedStore(storeID uint) (bool, error) {
	// 1. Retrieve store's LastInventorySyncAt timestamp
	var lastSyncTime time.Time
	err := m.db.Model(&data.Store{}).
		Select("last_inventory_sync_at").
		Where("id = ?", storeID).
		Scan(&lastSyncTime).Error
	if err != nil {
		return false, fmt.Errorf("failed to fetch store sync time: %w", err)
	}

	// 2. Fetch unsynchronized product size additives IDs
	var additiveIDs []uint
	err = m.db.Model(&data.ProductSizeAdditive{}).
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
	err = m.db.Model(&data.ProductSizeIngredient{}).
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
	err = m.db.Model(&data.AdditiveIngredient{}).
		Joins("JOIN additives ON additives.id = additive_ingredients.additive_id").
		Joins("JOIN store_additives ON additives.id = store_additives.additive_id").
		Where("store_additives.store_id = ?", storeID).
		Where("additive_ingredients.created_at > ?", lastSyncTime).
		Pluck("additive_ingredients.id", &additiveIngredientIDs).Error
	if err != nil {
		return false, fmt.Errorf("failed to fetch unsynchronized additive ingredients: %w", err)
	}

	// 5. Use the Filter function to determine missing IDs
	missingAdditiveIDs, err := m.storeAdditiveRepo.FilterMissingStoreAdditiveIDs(storeID, additiveIDs)
	if err != nil {
		return false, fmt.Errorf("failed to fetch missing additive ingredients: %w", err)
	}
	if len(missingAdditiveIDs) > 0 {
		return false, nil
	}

	mergedIngredientIDs := utils.MergeDistinct(ingredientIDs, additiveIngredientIDs)
	missingIngredientIDs, err := m.storeStockRepo.FilterMissingIngredientsIDs(storeID, mergedIngredientIDs)
	if err != nil {
		return false, fmt.Errorf("failed to fetch missing ingredients: %w", err)
	}
	if len(missingIngredientIDs) > 0 {
		return false, nil
	}

	// Store is synchronized if no missing IDs
	return true, nil
}

func (m *transactionManager) SynchronizeStoreInventory(storeID uint) error {
	store, err := m.storeRepo.GetRawStoreByID(storeID)
	if err != nil {
		return err
	}

	productSizeAdditivesIDs, err := m.repo.GetNotSynchronizedProductSizesAdditivesIDs(storeID, store.LastInventorySyncAt)
	if err != nil && !errors.Is(err, types.ErrNotSynchronizedProductSizeAdditivesNotFound) {
		return err
	}

	ingredientsFromProductSizesIDs, err := m.repo.GetNotSynchronizedProductSizesIngredients(storeID, store.LastInventorySyncAt)
	if err != nil && !errors.Is(err, types.ErrNotSynchronizedProductSizeIngredientsNotFound) {
		return err
	}

	ingredientsFromAdditives, err := m.repo.GetNotSynchronizedAdditivesIDs(storeID, store.LastInventorySyncAt)
	if err != nil && !errors.Is(err, types.ErrNotSynchronizedAdditivesNotFound) {
		return err
	}

	missingAdditivesIDs, err := m.storeAdditiveRepo.FilterMissingStoreAdditiveIDs(storeID, productSizeAdditivesIDs)
	if err != nil {
		return err
	}

	ingredientsToAdd := utils.MergeDistinct(ingredientsFromProductSizesIDs, ingredientsFromAdditives)

	missingIngredientsIDs, err := m.storeStockRepo.FilterMissingIngredientsIDs(storeID, ingredientsToAdd)
	if err != nil {
		return err
	}

	newStoreStocks := make([]data.StoreStock, len(missingIngredientsIDs))
	for i := 0; i < len(missingIngredientsIDs); i++ {
		newStoreStocks[i] = data.StoreStock{
			StoreID:      storeID,
			IngredientID: missingIngredientsIDs[i],
		}
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		if len(missingAdditivesIDs) > 0 {
			storeAdditiveList := make([]data.StoreAdditive, len(missingAdditivesIDs))
			for i := 0; i < len(storeAdditiveList); i++ {
				storeAdditive := data.StoreAdditive{
					StoreID:    storeID,
					AdditiveID: missingAdditivesIDs[i],
				}
				storeAdditiveList[i] = storeAdditive
			}

			m.storeAdditiveRepo.CloneWithTransaction(tx)
			_, err = m.storeAdditiveRepo.CreateStoreAdditives(storeAdditiveList)
			if err != nil {
				return err
			}
		}

		if len(newStoreStocks) > 0 {
			m.storeStockRepo.CloneWithTransaction(tx)
			_, err = m.storeStockRepo.AddMultipleStocks(newStoreStocks)
			if err != nil {
				return err
			}
		}

		store.LastInventorySyncAt = time.Now()
		m.storeRepo.CloneWithTransaction(tx)
		err = m.storeRepo.UpdateStore(storeID, &storesTypes.StoreUpdateModels{Store: store})
		if err != nil {
			return err
		}
		return nil
	})
}
