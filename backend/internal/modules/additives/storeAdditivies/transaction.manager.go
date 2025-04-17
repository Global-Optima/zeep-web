package storeAdditives

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
	storeInventoryManagersTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TransactionManager interface {
	CreateStoreAdditivesWithStocks(storeID uint, storeAdditive []data.StoreAdditive, ingredientIDs []uint) ([]uint, error)
	CheckSufficientInventoryByStoreAdditiveID(
		storeID, storeAdditiveID uint,
		frozenInventory *storeInventoryManagersTypes.FrozenInventory,
	) error
}

type transactionManager struct {
	db                        *gorm.DB
	storeAdditiveRepo         StoreAdditiveRepository
	storeStockRepo            storeStocks.StoreStockRepository
	ingredientRepo            ingredients.IngredientRepository
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository
}

func NewTransactionManager(
	db *gorm.DB,
	storeAdditiveRepo StoreAdditiveRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	ingredientRepo ingredients.IngredientRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
) TransactionManager {
	return &transactionManager{
		db:                        db,
		storeAdditiveRepo:         storeAdditiveRepo,
		storeStockRepo:            storeStockRepo,
		ingredientRepo:            ingredientRepo,
		storeInventoryManagerRepo: storeInventoryManagerRepo,
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

		storeStockRepo := m.storeStockRepo.CloneWithTransaction(tx)

		missingIngredientIDs, err := m.storeStockRepo.FilterMissingIngredientsIDs(storeID, ingredientIDs)
		if err != nil {
			return err
		}

		missingIngredients, err := m.ingredientRepo.GetIngredientsWithDetailsByIDs(missingIngredientIDs)
		if err != nil {
			return err
		}

		newStoreStocks := make([]data.StoreStock, len(missingIngredients))
		for i, ingredient := range missingIngredients {
			newStock, err := storeStocksTypes.DefaultStockFromIngredient(storeID, &ingredient)
			if err != nil {
				return err
			}
			newStoreStocks[i] = *newStock
		}

		if len(newStoreStocks) > 0 {
			_, err = m.addStocks(storeStockRepo, newStoreStocks)
			if err != nil {
				return err
			}
		}

		storeAdditiveIDs := make([]uint, len(storeAdditives))
		for i, storeAdditive := range storeAdditives {
			storeAdditiveIDs[i] = storeAdditive.ID
		}

		storeInventoryManagerRepoTx := m.storeInventoryManagerRepo.CloneWithTransaction(tx)
		frozenStockMap, err := storeInventoryManagerRepoTx.CalculateFrozenInventory(
			storeID,
			&storeInventoryManagersTypes.FrozenInventoryFilter{
				IngredientIDs: ingredientIDs,
			},
		)
		if err != nil {
			return err
		}

		if err := storeInventoryManagerRepoTx.RecalculateStoreAdditives(storeAdditiveIDs, storeID, frozenStockMap); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (m *transactionManager) CheckSufficientInventoryByStoreAdditiveID(
	storeID, storeAdditiveID uint,
	frozenInventory *storeInventoryManagersTypes.FrozenInventory,
) error {
	var storeAdditive data.StoreAdditive

	err := m.db.
		Where("store_id = ?", storeID).
		Where("id = ?", storeAdditiveID).
		Preload("Additive.Ingredients").
		Preload("Additive.AdditiveProvisions").
		First(&storeAdditive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.ErrStoreAdditiveNotFound
		}
		return fmt.Errorf("%w: failed to get storeAdditive (id=%d): %w",
			types.ErrFailedToFetchStoreAdditive, storeAdditiveID, err)
	}

	requiredIngredientMap := make(map[uint]float64)
	for _, ingr := range storeAdditive.Additive.Ingredients {
		requiredIngredientMap[ingr.IngredientID] += ingr.Quantity
	}

	if err := m.storeInventoryManagerRepo.CheckStoreStocks(storeID, requiredIngredientMap, frozenInventory); err != nil {
		return err
	}

	requiredProvisionMap := make(map[uint]float64)
	for _, prov := range storeAdditive.Additive.AdditiveProvisions {
		requiredProvisionMap[prov.ProvisionID] += prov.Volume
	}

	if err := m.storeInventoryManagerRepo.CheckStoreProvisions(storeID, requiredProvisionMap, frozenInventory); err != nil {
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
