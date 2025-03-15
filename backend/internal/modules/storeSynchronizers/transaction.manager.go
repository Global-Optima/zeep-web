package storeSynchronizers

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeSynchronizers/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
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

	unsyncData, err := m.fetchUnsynchronizedData(storeID, store.LastInventorySyncAt.UTC())
	if err != nil {
		return nil, err
	}

	if unsyncData == nil {
		return nil, fmt.Errorf("could not fetch unsyncData: nil pointer dereference")
	}

	if len(unsyncData.AdditiveIDs) == 0 && len(unsyncData.IngredientIDs) == 0 {
		return &types.SynchronizationStatus{
			LastSyncDate: store.LastInventorySyncAt.UTC().String(),
			IsSync:       true,
		}, nil
	}

	err = m.filterMissingData(storeID, unsyncData)
	if err != nil {
		return nil, err
	}

	if len(unsyncData.AdditiveIDs) > 0 || len(unsyncData.IngredientIDs) > 0 {
		return &types.SynchronizationStatus{
			LastSyncDate: store.LastInventorySyncAt.UTC().String(),
			IsSync:       false,
		}, nil
	}

	synchronizedAt, err := m.storeRepo.UpdateStoreSyncTime(storeID)
	if err != nil {
		return nil, err
	}

	return &types.SynchronizationStatus{
		LastSyncDate: synchronizedAt.UTC().String(),
		IsSync:       true,
	}, nil
}

func (m *transactionManager) SynchronizeStoreInventory(storeID uint) error {
	store, err := m.storeRepo.GetRawStoreByID(storeID)
	if err != nil {
		return err
	}

	unsyncData, err := m.fetchUnsynchronizedData(storeID, store.LastInventorySyncAt)
	if err != nil {
		return err
	}

	err = m.filterMissingData(storeID, unsyncData)
	if err != nil {
		return err
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		g := new(errgroup.Group)

		if len(unsyncData.AdditiveIDs) > 0 {
			g.Go(func() error {
				storeAdditiveList := make([]data.StoreAdditive, len(unsyncData.AdditiveIDs))
				for i, id := range unsyncData.AdditiveIDs {
					storeAdditiveList[i] = data.StoreAdditive{
						StoreID:    storeID,
						AdditiveID: id,
					}
				}
				storeAdditiveRepoTx := m.storeAdditiveRepo.CloneWithTransaction(tx)
				_, err := storeAdditiveRepoTx.CreateStoreAdditives(storeAdditiveList)
				return err
			})
		}

		if len(unsyncData.IngredientIDs) > 0 {
			g.Go(func() error {
				newStoreStocks := make([]data.StoreStock, len(unsyncData.IngredientIDs))
				for i, ingredientID := range unsyncData.IngredientIDs {
					newStoreStocks[i] = *storeStocksTypes.DefaultStockFromIngredient(storeID, ingredientID)
				}
				storeStockRepoTx := m.storeStockRepo.CloneWithTransaction(tx)
				_, err := storeStockRepoTx.AddMultipleStocks(newStoreStocks)
				return err
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		storeRepoTx := m.storeRepo.CloneWithTransaction(tx)
		_, err = storeRepoTx.UpdateStoreSyncTime(storeID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (m *transactionManager) fetchUnsynchronizedData(storeID uint, lastSyncAt time.Time) (*types.UnsyncData, error) {
	var (
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
		ingredientIDsFromProducts, err = m.repo.GetNotSynchronizedProductSizeWithAdditivesIngredients(storeID, lastSyncAt)
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
		AdditiveIDs:   additiveIDs,
		IngredientIDs: utils.MergeDistinct(ingredientIDsFromAdditives, ingredientIDsFromProducts),
	}, nil
}

func (m *transactionManager) filterMissingData(storeID uint, data *types.UnsyncData) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		var err error
		data.AdditiveIDs, err = m.storeAdditiveRepo.FilterMissingStoreAdditiveIDs(storeID, data.AdditiveIDs)
		return err
	})

	g.Go(func() error {
		var err error
		data.IngredientIDs, err = m.storeStockRepo.FilterMissingIngredientsIDs(storeID, data.IngredientIDs)
		return err
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
