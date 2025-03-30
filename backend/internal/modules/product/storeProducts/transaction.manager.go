package storeProducts

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"gorm.io/gorm"
)

type TransactionManager interface {
	CreateStoreProductWithStocks(storeID uint, storeProduct *data.StoreProduct, storeAdditives []data.StoreAdditive, ingredientIDs []uint) (storeProductID uint, storeAdditiveIDs []uint, err error)
	CreateMultipleStoreProductsWithStocks(storeID uint, storeProduct []data.StoreProduct, storeAdditives []data.StoreAdditive, ingredientIDs []uint) (storeProductIDs []uint, storeAdditiveIDs []uint, err error)
	UpdateStoreProductWithStocks(storeID, storeProductID uint, updateModels *types.StoreProductModels, ingredientIDs []uint) error
}

type transactionManager struct {
	db                *gorm.DB
	storeProductRepo  StoreProductRepository
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository
	storeStockRepo    storeStocks.StoreStockRepository
}

func NewTransactionManager(db *gorm.DB, storeProductRepo StoreProductRepository, storeAdditiveRepo storeAdditives.StoreAdditiveRepository, storeStockRepo storeStocks.StoreStockRepository) TransactionManager {
	return &transactionManager{
		db:                db,
		storeProductRepo:  storeProductRepo,
		storeAdditiveRepo: storeAdditiveRepo,
		storeStockRepo:    storeStockRepo,
	}
}

func (m *transactionManager) CreateStoreProductWithStocks(storeID uint, storeProduct *data.StoreProduct, storeAdditives []data.StoreAdditive, ingredientIDs []uint) (uint, []uint, error) {
	var id uint
	var storeAdditiveIDs []uint
	err := m.db.Transaction(func(tx *gorm.DB) error {
		var err error
		sp := m.storeProductRepo.CloneWithTransaction(tx)
		storeProduct.StoreID = storeID
		id, err = sp.CreateStoreProduct(storeProduct)
		if err != nil {
			return err
		}

		storeStockRepo := m.storeStockRepo.CloneWithTransaction(tx)

		storeAdditiveRepoTx := m.storeAdditiveRepo.CloneWithTransaction(tx)
		storeAdditiveIDs, err = storeAdditiveRepoTx.CreateStoreAdditives(storeAdditives)
		if err != nil {
			return err
		}

		missingIngredientIDs, err := storeStockRepo.FilterMissingIngredientsIDs(storeID, ingredientIDs)
		if err != nil {
			return err
		}

		newStoreStocks := make([]data.StoreStock, len(missingIngredientIDs))
		for i, ingredientID := range missingIngredientIDs {
			newStoreStocks[i] = *storeStocksTypes.DefaultStockFromIngredient(storeID, ingredientID)
		}

		if len(newStoreStocks) > 0 {
			_, err = m.addStocks(&storeStockRepo, newStoreStocks)
			if err != nil {
				return err
			}
		}

		if err := data.RecalculateOutOfStock(tx, storeID, ingredientIDs, nil); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, nil, err
	}
	return id, storeAdditiveIDs, nil
}

func (m *transactionManager) CreateMultipleStoreProductsWithStocks(storeID uint, storeProducts []data.StoreProduct, storeAdditives []data.StoreAdditive, ingredientIDs []uint) ([]uint, []uint, error) {
	var storeProductIDs, storeAdditiveIDs []uint
	err := m.db.Transaction(func(tx *gorm.DB) error {
		sp := m.storeProductRepo.CloneWithTransaction(tx)

		for _, storeProduct := range storeProducts {
			storeProduct.StoreID = storeID
		}

		for _, storeAdditive := range storeAdditives {
			storeAdditive.StoreID = storeID
		}

		var err error
		storeProductIDs, err = sp.CreateMultipleStoreProducts(storeProducts)
		if err != nil {
			return err
		}

		storeAdditiveRepoTx := m.storeAdditiveRepo.CloneWithTransaction(tx)
		storeAdditiveIDs, err = storeAdditiveRepoTx.CreateStoreAdditives(storeAdditives)
		if err != nil {
			return err
		}

		storeStockRepoTx := m.storeStockRepo.CloneWithTransaction(tx)
		missingIngredientIDs, err := storeStockRepoTx.FilterMissingIngredientsIDs(storeID, ingredientIDs)
		if err != nil {
			return err
		}

		newStoreStocks := make([]data.StoreStock, len(missingIngredientIDs))
		for i, ingredientID := range missingIngredientIDs {
			newStoreStocks[i] = *storeStocksTypes.DefaultStockFromIngredient(storeID, ingredientID)
		}

		if len(newStoreStocks) > 0 {
			_, err = m.addStocks(&storeStockRepoTx, newStoreStocks)
			if err != nil {
				return err
			}
		}

		if err := data.RecalculateOutOfStock(tx, storeID, ingredientIDs, nil); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return storeProductIDs, storeAdditiveIDs, nil
}

func (m *transactionManager) UpdateStoreProductWithStocks(storeID, storeProductID uint, updateModels *types.StoreProductModels, ingredientIDs []uint) error {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		sp := m.storeProductRepo.CloneWithTransaction(tx)
		if err := sp.UpdateStoreProductByID(storeID, storeProductID, updateModels); err != nil {
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

		if len(newStoreStocks) > 0 {
			_, err = m.addStocks(&storeStockRepo, newStoreStocks)
			if err != nil {
				return err
			}
		}

		if err := data.RecalculateOutOfStock(tx, storeID, ingredientIDs, nil); err != nil {
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
