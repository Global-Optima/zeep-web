package storeSynchronizers

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeSynchronizers/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type TransactionManager interface {
	SynchronizeStoreInventory(storeID uint) error
	GetSynchronizationStatus(storeID uint) (*types.SynchronizationStatus, error)
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

func (m *transactionManager) GetSynchronizationStatus(storeID uint) (*types.SynchronizationStatus, error) {
	store, err := m.storeRepo.GetRawStoreByID(storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch store sync time: %w", err)
	}

	syncStatus := &types.SynchronizationStatus{
		LastSyncDate: store.LastInventorySyncAt.UTC().String(),
		IsSync:       true,
	}

	unsyncData, err := m.fetchUnsynchronizedData(storeID, store.LastInventorySyncAt.UTC())
	if err != nil {
		return nil, err
	}

	if unsyncData == nil {
		return nil, fmt.Errorf("could not fetch unsyncData: nil pointer dereference")
	}
	logrus.Info(unsyncData)
	if len(unsyncData.AdditiveIDs) > 0 || len(unsyncData.IngredientIDs) > 0 || len(unsyncData.ProductSizeIDs) > 0 {
		syncStatus.IsSync = false
	}

	return syncStatus, nil
}

func (m *transactionManager) SynchronizeStoreInventory(storeID uint) error {
	store, err := m.storeRepo.GetRawStoreByID(storeID)
	if err != nil {
		return err
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		productSizeIDs, err := m.repo.GetNotSynchronizedProductSizesIDs(storeID, store.LastInventorySyncAt)
		if err != nil {
			return err
		}

		if err := m.synchronizeAdditives(tx, storeID, store.LastInventorySyncAt); err != nil {
			return err
		}

		synchronizedIngredientIDS, err := m.synchronizeIngredients(tx, storeID, store.LastInventorySyncAt)
		if err != nil {
			return err
		}

		err = data.RecalculateOutOfStock(tx, storeID, synchronizedIngredientIDS, productSizeIDs)
		if err != nil {
			return err
		}

		storeRepoTx := m.storeRepo.CloneWithTransaction(tx)
		_, err = storeRepoTx.UpdateStoreSyncTime(storeID)
		if err != nil {
			return fmt.Errorf("failed to update store sync time: %w", err)
		}
		return nil
	})
}

func (m *transactionManager) fetchUnsynchronizedData(storeID uint, lastSyncAt time.Time) (*types.UnsyncData, error) {
	var (
		productSizeIDs             []uint
		additiveIDs                []uint
		ingredientIDsFromProducts  []uint
		ingredientIDsFromAdditives []uint
	)

	g := new(errgroup.Group)

	g.Go(func() error {
		var err error
		additiveIDs, err = m.repo.GetNotSynchronizedProductSizesAdditivesIDs(storeID, lastSyncAt)
		return err
	})

	g.Go(func() error {
		var err error
		ingredientIDsFromProducts, err = m.repo.GetNotSynchronizedProductSizeIngredientsIDs(storeID, lastSyncAt)
		return err
	})

	g.Go(func() error {
		var err error
		productSizeIDs, err = m.repo.GetNotSynchronizedProductSizesIDs(storeID, lastSyncAt)
		return err
	})

	g.Go(func() error {
		var err error
		ingredientIDsFromAdditives, err = m.repo.GetNotSynchronizedAdditiveIngredientsIDs(storeID, lastSyncAt)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &types.UnsyncData{
		ProductSizeIDs: productSizeIDs,
		AdditiveIDs:    additiveIDs,
		IngredientIDs:  utils.UnionSlices(ingredientIDsFromAdditives, ingredientIDsFromProducts),
	}, nil
}

func (m *transactionManager) synchronizeAdditives(tx *gorm.DB, storeID uint, lastSyncAt time.Time) error {
	additiveIDs, err := m.repo.GetNotSynchronizedProductSizesAdditivesIDs(
		storeID, lastSyncAt,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch unsynced additiveIDs: %w", err)
	}

	missingAdditives, err := m.storeAdditiveRepo.FilterMissingStoreAdditiveIDs(storeID, additiveIDs)
	if err != nil {
		return fmt.Errorf("failed to filter missing store_additives: %w", err)
	}

	if len(missingAdditives) > 0 {
		storeAdditiveList := make([]data.StoreAdditive, len(missingAdditives))
		for i, addID := range missingAdditives {
			storeAdditiveList[i] = data.StoreAdditive{
				StoreID:    storeID,
				AdditiveID: addID,
			}
		}
		storeAdditiveRepoTx := m.storeAdditiveRepo.CloneWithTransaction(tx)
		_, err := storeAdditiveRepoTx.CreateStoreAdditives(storeAdditiveList)
		if err != nil {
			return fmt.Errorf("failed to add store_additives: %w", err)
		}
	}

	return nil
}

func (m *transactionManager) synchronizeIngredients(tx *gorm.DB, storeID uint, lastSyncAt time.Time) ([]uint, error) {
	var productSizeIngredientIDs, additiveIngredientIDs []uint

	var err error
	productSizeIngredientIDs, err = m.repo.GetNotSynchronizedProductSizeIngredientsIDs(
		storeID, lastSyncAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product-size ingredients: %w", err)
	}

	additiveIngredientIDs, err = m.repo.GetNotSynchronizedAdditiveIngredientsIDs(
		storeID, lastSyncAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch additive-based ingredients: %w", err)
	}

	mergedIngredientIDs := utils.UnionSlices(productSizeIngredientIDs, additiveIngredientIDs)

	missingIngredientIDs, err := m.storeStockRepo.FilterMissingIngredientsIDs(storeID, mergedIngredientIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to filter missing store_stocks: %w", err)
	}

	if len(missingIngredientIDs) > 0 {
		newStoreStocks := make([]data.StoreStock, len(missingIngredientIDs))
		for i, ingID := range missingIngredientIDs {
			newStoreStocks[i] = *storeStocksTypes.DefaultStockFromIngredient(storeID, ingID)
		}
		storeStockRepoTx := m.storeStockRepo.CloneWithTransaction(tx)
		_, err := storeStockRepoTx.AddMultipleStocks(newStoreStocks)
		if err != nil {
			return nil, fmt.Errorf("failed to add store_stocks: %w", err)
		}
	}
	return mergedIngredientIDs, nil
}
