package storeProvisions

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"gorm.io/gorm"
)

type TransactionManager interface {
	CreateStoreProvisionWithStocks(storeProvision *data.StoreProvision, ingredientIDs []uint) (storeProvisionID uint, err error)
}

type transactionManager struct {
	db                 *gorm.DB
	storeProvisionRepo StoreProvisionRepository
	storeStockRepo     storeStocks.StoreStockRepository
	ingredientRepo     ingredients.IngredientRepository
}

func NewTransactionManager(
	db *gorm.DB,
	storeProvisionRepo StoreProvisionRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	ingredientRepo ingredients.IngredientRepository,
) TransactionManager {
	return &transactionManager{
		db:                 db,
		storeProvisionRepo: storeProvisionRepo,
		storeStockRepo:     storeStockRepo,
		ingredientRepo:     ingredientRepo,
	}
}

func (m *transactionManager) CreateStoreProvisionWithStocks(storeProvision *data.StoreProvision, ingredientIDs []uint) (storeProvisionID uint, err error) {
	var id uint

	err = m.db.Transaction(func(tx *gorm.DB) error {
		var err error
		sp := m.storeProvisionRepo.CloneWithTransaction(tx)
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
