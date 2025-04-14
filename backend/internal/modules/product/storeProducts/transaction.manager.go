package storeProducts

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
	storeInventoryManagersTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TransactionManager interface {
	CreateStoreProductWithStocks(storeID uint, storeProduct *data.StoreProduct, storeAdditives []data.StoreAdditive, ingredientIDs []uint) (storeProductID uint, storeAdditiveIDs []uint, err error)
	CreateMultipleStoreProductsWithStocks(storeID uint, storeProduct []data.StoreProduct, storeAdditives []data.StoreAdditive, ingredientIDs []uint) (storeProductIDs []uint, storeAdditiveIDs []uint, err error)
	UpdateStoreProductWithStocks(storeID, storeProductID uint, updateModels *types.StoreProductModels, ingredientIDs []uint) error
	CheckSufficientStoreProductSizeById(
		storeID, storeProductSizeID uint,
		frozenInventory *storeInventoryManagersTypes.FrozenInventory,
	) error
}

type transactionManager struct {
	db                        *gorm.DB
	storeProductRepo          StoreProductRepository
	storeAdditiveRepo         storeAdditives.StoreAdditiveRepository
	storeStockRepo            storeStocks.StoreStockRepository
	ingredientRepo            ingredients.IngredientRepository
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository
}

func NewTransactionManager(
	db *gorm.DB,
	storeProductRepo StoreProductRepository,
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	ingredientRepo ingredients.IngredientRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
) TransactionManager {
	return &transactionManager{
		db:                        db,
		storeProductRepo:          storeProductRepo,
		storeAdditiveRepo:         storeAdditiveRepo,
		storeStockRepo:            storeStockRepo,
		ingredientRepo:            ingredientRepo,
		storeInventoryManagerRepo: storeInventoryManagerRepo,
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

		storeStockRepoTx := m.storeStockRepo.CloneWithTransaction(tx)

		storeAdditiveRepoTx := m.storeAdditiveRepo.CloneWithTransaction(tx)
		storeAdditiveIDs, err = storeAdditiveRepoTx.CreateStoreAdditives(storeAdditives)
		if err != nil {
			return err
		}

		missingIngredientIDs, err := storeStockRepoTx.FilterMissingIngredientsIDs(storeID, ingredientIDs)
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
			_, err = m.addStocks(storeStockRepoTx, newStoreStocks)
			if err != nil {
				return err
			}
		}

		storeInventoryManagerRepoTx := m.storeInventoryManagerRepo.CloneWithTransaction(tx)
		if err := storeInventoryManagerRepoTx.RecalculateStoreInventory(storeID, &storeInventoryManagersTypes.RecalculateInput{IngredientIDs: ingredientIDs}); err != nil {
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
			_, err = m.addStocks(storeStockRepoTx, newStoreStocks)
			if err != nil {
				return err
			}
		}

		storeInventoryManagerRepoTx := m.storeInventoryManagerRepo.CloneWithTransaction(tx)
		if err := storeInventoryManagerRepoTx.RecalculateStoreInventory(storeID, &storeInventoryManagersTypes.RecalculateInput{IngredientIDs: ingredientIDs}); err != nil {
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
		storeStockRepoTx := m.storeStockRepo.CloneWithTransaction(tx)

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
			_, err = m.addStocks(storeStockRepoTx, newStoreStocks)
			if err != nil {
				return err
			}
		}

		storeInventoryManagerRepoTx := m.storeInventoryManagerRepo.CloneWithTransaction(tx)
		if err := storeInventoryManagerRepoTx.RecalculateStoreInventory(storeID, &storeInventoryManagersTypes.RecalculateInput{IngredientIDs: ingredientIDs}); err != nil {
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

func (m *transactionManager) CheckSufficientStoreProductSizeById(
	storeID, storeProductSizeID uint,
	frozenInventory *storeInventoryManagersTypes.FrozenInventory,
) error {
	var storeProductSize data.StoreProductSize

	err := m.db.Model(&data.StoreProductSize{}).
		Joins("JOIN store_products sp ON store_product_sizes.store_product_id = sp.id").
		Where("sp.store_id = ? AND store_product_sizes.id = ?", storeID, storeProductSizeID).
		Preload("ProductSize.Product").
		Preload("ProductSize.ProductSizeIngredients").
		Preload("ProductSize.ProductSizeProvisions").
		Preload("ProductSize.Additives.Additive.Ingredients").
		Preload("ProductSize.Additives.Additive.AdditiveProvisions").
		First(&storeProductSize).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.ErrStoreProductSizeNotFound
		}
		return fmt.Errorf("%w: failed to load StoreProductSize (ID=%d): %w",
			types.ErrFailedToFetchStoreProductSize, storeProductSizeID, err)
	}

	if err := m.checkProductSizeIngredients(storeID, &storeProductSize, frozenInventory); err != nil {
		return err
	}

	if err := m.checkProductSizeProvisions(storeID, &storeProductSize, frozenInventory); err != nil {
		return err
	}
	logrus.Infof("Final FROZEN STOCK CHECK: %v", frozenInventory)

	return nil
}

func (m *transactionManager) checkProductSizeIngredients(
	storeID uint,
	storeProductSize *data.StoreProductSize,
	frozenInventory *storeInventoryManagersTypes.FrozenInventory,
) error {
	requiredIngredientQuantityMap := make(map[uint]float64)
	// Check base ingredients
	for _, usage := range storeProductSize.ProductSize.ProductSizeIngredients {
		requiredIngredientQuantityMap[usage.IngredientID] += usage.Quantity
	}

	// Check ingredients from default additives
	for _, psa := range storeProductSize.ProductSize.Additives {
		if !psa.IsDefault {
			continue
		}
		for _, additiveIng := range psa.Additive.Ingredients {
			requiredIngredientQuantityMap[additiveIng.IngredientID] += additiveIng.Quantity
		}
	}

	return m.storeInventoryManagerRepo.CheckStoreStocks(storeID, requiredIngredientQuantityMap, frozenInventory)
}

func (m *transactionManager) checkProductSizeProvisions(
	storeID uint,
	storeProductSize *data.StoreProductSize,
	frozenInventory *storeInventoryManagersTypes.FrozenInventory,
) error {
	requiredProvisionQuantityMap := make(map[uint]float64)

	// Gather base product size provisions
	for _, usage := range storeProductSize.ProductSize.ProductSizeProvisions {
		requiredProvisionQuantityMap[usage.ProvisionID] += usage.Volume
	}

	// Gather provisions from default additives
	for _, psa := range storeProductSize.ProductSize.Additives {
		if !psa.IsDefault {
			continue
		}
		for _, additiveProv := range psa.Additive.AdditiveProvisions {
			requiredProvisionQuantityMap[additiveProv.ProvisionID] += additiveProv.Volume
		}
	}

	// Call the provision check
	return m.storeInventoryManagerRepo.CheckStoreProvisions(storeID, requiredProvisionQuantityMap, frozenInventory)
}
