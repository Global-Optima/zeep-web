package storeProducts

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"gorm.io/gorm"
)

type TransactionManager interface {
	CreateStoreProductWithStocks(storeID uint, storeProduct *data.StoreProduct, storeAdditives []data.StoreAdditive, stockDTOs []storeStocksTypes.AddStoreStockDTO) (storeProductID uint, storeAdditiveIDs []uint, err error)
	CreateMultipleStoreProductsWithStocks(storeID uint, storeProduct []data.StoreProduct, storeAdditives []data.StoreAdditive, stockDTOs []storeStocksTypes.AddStoreStockDTO) (storeProductIDs []uint, storeAdditiveIDs []uint, err error)
	UpdateStoreProductWithStocks(storeID, storeProductID uint, updateModels *types.StoreProductModels, dtos []storeStocksTypes.AddStoreStockDTO) error
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

func (m *transactionManager) CreateStoreProductWithStocks(storeID uint, storeProduct *data.StoreProduct, storeAdditives []data.StoreAdditive, stockDTOs []storeStocksTypes.AddStoreStockDTO) (uint, []uint, error) {
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

		if err := m.addStocks(&storeStockRepo, storeID, stockDTOs); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, nil, err
	}
	return id, storeAdditiveIDs, nil
}

func (m *transactionManager) CreateMultipleStoreProductsWithStocks(storeID uint, storeProducts []data.StoreProduct, storeAdditives []data.StoreAdditive, storeStockDTOs []storeStocksTypes.AddStoreStockDTO) ([]uint, []uint, error) {
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

		storeWarehouseRepoTx := m.storeStockRepo.CloneWithTransaction(tx)
		if err := m.addStocks(&storeWarehouseRepoTx, storeID, storeStockDTOs); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return storeProductIDs, storeAdditiveIDs, nil
}

func (m *transactionManager) UpdateStoreProductWithStocks(storeID, storeProductID uint, updateModels *types.StoreProductModels, dtos []storeStocksTypes.AddStoreStockDTO) error {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		sp := m.storeProductRepo.CloneWithTransaction(tx)
		if err := sp.UpdateStoreProductByID(storeID, storeProductID, updateModels); err != nil {
			return err
		}
		storeWarehouseRepo := m.storeStockRepo.CloneWithTransaction(tx)

		if err := m.addStocks(&storeWarehouseRepo, storeID, dtos); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *transactionManager) addStocks(storeStockRepo storeStocks.StoreStockRepository, storeID uint, dtos []storeStocksTypes.AddStoreStockDTO) error {
	for _, dto := range dtos {
		_, err := storeStockRepo.AddStock(storeID, &dto)
		if err != nil {
			switch {
			case errors.Is(err, moduleErrors.ErrAlreadyExists):
				continue
			}
			return err
		}
	}
	return nil
}
