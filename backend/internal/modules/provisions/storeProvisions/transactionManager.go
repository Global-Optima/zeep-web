package storeProvisions

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"gorm.io/gorm"
	"time"
)

type TransactionManager interface {
	CreateStoreProvisionWithStocks(storeProvision *data.StoreProvision, ingredientIDs []uint) (storeProvisionID uint, err error)
	CompleteStoreProvision(existingProvision *data.StoreProvision) ([]data.StoreStock, error)
}

type transactionManager struct {
	db                        *gorm.DB
	repo                      StoreProvisionRepository
	storeStockRepo            storeStocks.StoreStockRepository
	ingredientRepo            ingredients.IngredientRepository
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository
}

func NewTransactionManager(
	db *gorm.DB,
	storeProvisionRepo StoreProvisionRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	ingredientRepo ingredients.IngredientRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
) TransactionManager {
	return &transactionManager{
		db:                        db,
		repo:                      storeProvisionRepo,
		storeStockRepo:            storeStockRepo,
		ingredientRepo:            ingredientRepo,
		storeInventoryManagerRepo: storeInventoryManagerRepo,
	}
}

func (m *transactionManager) CreateStoreProvisionWithStocks(storeProvision *data.StoreProvision, ingredientIDs []uint) (storeProvisionID uint, err error) {
	var id uint

	err = m.db.Transaction(func(tx *gorm.DB) error {
		var err error
		sp := m.repo.CloneWithTransaction(tx)
		id, err = sp.CreateStoreProvision(storeProvision)
		if err != nil {
			return err
		}

		storeStockRepoTx := m.storeStockRepo.CloneWithTransaction(tx)

		missingIngredientIDs, err := storeStockRepoTx.FilterMissingIngredientsIDs(storeProvision.StoreID, ingredientIDs)
		if err != nil {
			return err
		}

		missingIngredients, err := m.ingredientRepo.GetIngredientsWithDetailsByIDs(missingIngredientIDs)
		if err != nil {
			return err
		}

		newStoreStocks := make([]data.StoreStock, len(missingIngredients))
		for i, ingredient := range missingIngredients {
			newStock, err := storeStocksTypes.DefaultStockFromIngredient(storeProvision.StoreID, &ingredient)
			if err != nil {
				return err
			}
			newStoreStocks[i] = *newStock
		}

		if len(newStoreStocks) > 0 {
			_, err = m.addStocks(storeStockRepoTx, newStoreStocks)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *transactionManager) addStocks(storeStockRepo storeStocks.StoreStockRepository, stocks []data.StoreStock) ([]uint, error) {
	ids, err := storeStockRepo.AddMultipleStocks(stocks)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (m *transactionManager) CompleteStoreProvision(provision *data.StoreProvision) ([]data.StoreStock, error) {
	if provision == nil || provision.ID == 0 || provision.StoreID == 0 {
		return nil, fmt.Errorf("invalid input arguments")
	}

	if provision.Status != data.STORE_PROVISION_STATUS_PREPARING {
		wrapped := fmt.Errorf("failed to update store provision: %w", types.ErrProvisionCompleted)
		return nil, wrapped
	}

	provision.Status = data.STORE_PROVISION_STATUS_COMPLETED

	currentTime := time.Now().UTC()
	provision.CompletedAt = &currentTime

	expirationTime := currentTime.Add(time.Duration(provision.ExpirationInMinutes) * time.Minute)
	provision.ExpiresAt = &expirationTime

	var deductedStocks []data.StoreStock

	err := m.db.Transaction(func(tx *gorm.DB) error {
		var err error

		storeInventoryManagerRepoTx := m.storeInventoryManagerRepo.CloneWithTransaction(tx)
		deductedStocks, err = storeInventoryManagerRepoTx.DeductStoreStocksByStoreProvision(provision)
		if err != nil {
			wrapped := fmt.Errorf("failed to deduct store stocks: %w", err)
			return wrapped
		}

		repoTx := m.repo.CloneWithTransaction(tx)
		err = repoTx.SaveStoreProvision(provision)
		if err != nil {
			wrapped := fmt.Errorf("failed to complete store provision: %w", err)
			return wrapped
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return deductedStocks, nil
}
