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
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type TransactionManager interface {
	SynchronizeStoreInventory(storeID uint) error
	IsSynchronizedStore(storeID uint) (bool, error)
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
	unsyncData, err := m.fetchUnsynchronizedData(storeID)
	if err != nil {
		return false, err
	}

	if unsyncData == nil {
		return false, fmt.Errorf("could not fetch unsyncData: nil pointer dereference")
	}

	logrus.Info("1: ", unsyncData)
	if len(unsyncData.AdditiveIDs) == 0 && len(unsyncData.IngredientIDs) == 0 {
		return true, nil
	}

	err = m.filterMissingData(storeID, unsyncData)
	if err != nil {
		return false, err
	}
	logrus.Info("2: ", unsyncData)

	if len(unsyncData.AdditiveIDs) > 0 || len(unsyncData.IngredientIDs) > 0 {
		return false, nil
	}

	err = m.storeRepo.UpdateStoreSyncTime(storeID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *transactionManager) SynchronizeStoreInventory(storeID uint) error {
	store, err := m.storeRepo.GetRawStoreByID(storeID)
	if err != nil {
		return err
	}

	// Fetch unsynchronized and missing data concurrently
	unsyncData, err := m.fetchUnsynchronizedData(storeID)
	if err != nil {
		return err
	}

	err = m.filterMissingData(storeID, unsyncData)
	if err != nil {
		return err
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		g := new(errgroup.Group)

		// Add missing additives
		if len(unsyncData.AdditiveIDs) > 0 {
			g.Go(func() error {
				storeAdditiveList := make([]data.StoreAdditive, len(unsyncData.AdditiveIDs))
				for i, id := range unsyncData.AdditiveIDs {
					storeAdditiveList[i] = data.StoreAdditive{
						StoreID:    storeID,
						AdditiveID: id,
					}
				}
				m.storeAdditiveRepo.CloneWithTransaction(tx)
				_, err := m.storeAdditiveRepo.CreateStoreAdditives(storeAdditiveList)
				return err
			})
		}

		// Add missing store stocks
		if len(unsyncData.IngredientIDs) > 0 {
			g.Go(func() error {
				newStoreStocks := make([]data.StoreStock, len(unsyncData.IngredientIDs))
				for i, id := range unsyncData.IngredientIDs {
					newStoreStocks[i] = data.StoreStock{
						StoreID:           storeID,
						IngredientID:      id,
						Quantity:          0,
						LowStockThreshold: storeStocks.DEFAULT_LOW_STOCK_THRESHOLD,
					}
				}
				m.storeStockRepo.CloneWithTransaction(tx)
				_, err := m.storeStockRepo.AddMultipleStocks(newStoreStocks)
				return err
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		store.LastInventorySyncAt = time.Now()
		m.storeRepo.CloneWithTransaction(tx)
		return m.storeRepo.UpdateStore(storeID, &storesTypes.StoreUpdateModels{Store: store})
	})
}

func (m *transactionManager) fetchUnsynchronizedData(storeID uint) (*types.UnsyncData, error) {
	var (
		additiveIDs                []uint
		ingredientIDsFromProducts  []uint
		ingredientIDsFromAdditives []uint
	)

	store, err := m.storeRepo.GetRawStoreByID(storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch store sync time: %w", err)
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		var err error
		additiveIDs, err = m.repo.GetNotSynchronizedProductSizesAdditivesIDs(storeID, store.LastInventorySyncAt)
		return err
	})

	g.Go(func() error {
		var err error
		ingredientIDsFromProducts, err = m.repo.GetNotSynchronizedProductSizeWithAdditivesIngredients(storeID, store.LastInventorySyncAt)
		return err
	})

	g.Go(func() error {
		var err error
		ingredientIDsFromAdditives, err = m.repo.GetNotSynchronizedAdditiveIngredientsIDs(storeID, store.LastInventorySyncAt)
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
