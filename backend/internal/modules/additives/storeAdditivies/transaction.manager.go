package storeAdditives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"gorm.io/gorm"
)

type TransactionManager interface {
	CreateStoreAdditivesWithStocks(storeID uint, storeAdditive []data.StoreAdditive, ingredientIDs []uint) ([]uint, error)
}

type transactionManager struct {
	db                *gorm.DB
	storeAdditiveRepo StoreAdditiveRepository
	storeStockRepo    storeStocks.StoreStockRepository
}

func NewTransactionManager(db *gorm.DB, storeAdditiveRepo StoreAdditiveRepository, storeStockRepo storeStocks.StoreStockRepository) TransactionManager {
	return &transactionManager{
		db:                db,
		storeAdditiveRepo: storeAdditiveRepo,
		storeStockRepo:    storeStockRepo,
	}
}

func (m *transactionManager) CreateStoreAdditivesWithStocks(storeID uint, storeAdditives []data.StoreAdditive, ingredientIDs []uint) ([]uint, error) {
	var ids []uint
	err := m.db.Transaction(func(tx *gorm.DB) error {
		sa := m.storeAdditiveRepo.CloneWithTransaction(tx)

		for _, storeAdditive := range storeAdditives {
			storeAdditive.StoreID = storeID
		}

		var err error
		ids, err = sa.CreateStoreAdditives(storeAdditives)
		if err != nil {
			return err
		}

		storeWarehouseRepo := m.storeStockRepo.CloneWithTransaction(tx)

		missingIngredientIDs, err := m.storeStockRepo.FilterMissingIngredientsIDs(storeID, ingredientIDs)
		if err != nil {
			return err
		}

		newStoreStocks := make([]data.StoreStock, len(missingIngredientIDs))
		for i, ingredientID := range missingIngredientIDs {
			newStoreStocks[i] = *storeStocksTypes.DefaultStockFromIngredient(storeID, ingredientID)
		}

		_, err = m.addStocks(&storeWarehouseRepo, newStoreStocks)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (m *transactionManager) UpdateStoreAdditivesWithStocks(storeID, storeAdditiveID uint, updateStoreAdditive *data.StoreAdditive, ingredientIDs []uint) error {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		sp := m.storeAdditiveRepo.CloneWithTransaction(tx)
		if err := sp.UpdateStoreAdditive(storeID, storeAdditiveID, updateStoreAdditive); err != nil {
			return err
		}

		storeStockRepo := m.storeStockRepo.CloneWithTransaction(tx)

		missingIngredientIDs, err := m.storeStockRepo.FilterMissingIngredientsIDs(storeID, ingredientIDs)
		if err != nil {
			return err
		}

		newStoreStocks := make([]data.StoreStock, len(missingIngredientIDs))
		for i, ingredientID := range missingIngredientIDs {
			newStoreStocks[i] = *storeStocksTypes.DefaultStockFromIngredient(storeID, ingredientID)
		}

		_, err = m.addStocks(&storeStockRepo, newStoreStocks)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *transactionManager) addStocks(storeStockRepo storeStocks.StoreStockRepository, stocks []data.StoreStock) ([]uint, error) {
	ids, err := storeStockRepo.AddMultipleStocks(stocks)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
