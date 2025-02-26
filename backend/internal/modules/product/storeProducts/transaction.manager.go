package storeProducts

import (
	"errors"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	storeWarehousesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"gorm.io/gorm"
)

type TransactionManager interface {
	CreateStoreProductWithStocks(storeID uint, storeProduct *data.StoreProduct, dtos []storeWarehousesTypes.AddStoreStockDTO) (uint, error)
	CreateMultipleStoreProductsWithStocks(storeID uint, storeProducts []data.StoreProduct, dtos []storeWarehousesTypes.AddStoreStockDTO) ([]uint, error)
	UpdateStoreProductWithStocks(storeID, storeProductID uint, updateModels *types.StoreProductModels, dtos []storeWarehousesTypes.AddStoreStockDTO) error
}

type transactionManager struct {
	db               *gorm.DB
	storeProductRepo StoreProductRepository
	storeStockRepo   storeStocks.StoreStockRepository
}

func NewTransactionManager(db *gorm.DB, storeProductRepo StoreProductRepository, storeStockRepo storeStocks.StoreStockRepository) TransactionManager {
	return &transactionManager{
		db:               db,
		storeProductRepo: storeProductRepo,
		storeStockRepo:   storeStockRepo,
	}
}

func (m *transactionManager) CreateStoreProductWithStocks(storeID uint, storeProduct *data.StoreProduct, stockDTOs []storeWarehousesTypes.AddStoreStockDTO) (uint, error) {
	var id uint
	err := m.db.Transaction(func(tx *gorm.DB) error {
		var err error
		sp := m.storeProductRepo.CloneWithTransaction(tx)
		storeProduct.StoreID = storeID
		id, err = sp.CreateStoreProduct(storeProduct)
		if err != nil {
			return err
		}

		storeWarehouseRepo := m.storeStockRepo.CloneWithTransaction(tx)

		if err := m.addStocks(&storeWarehouseRepo, storeID, stockDTOs); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *transactionManager) CreateMultipleStoreProductsWithStocks(storeID uint, storeProducts []data.StoreProduct, dtos []storeWarehousesTypes.AddStoreStockDTO) ([]uint, error) {
	var ids []uint
	err := m.db.Transaction(func(tx *gorm.DB) error {
		sp := m.storeProductRepo.CloneWithTransaction(tx)

		for _, storeProduct := range storeProducts {
			storeProduct.StoreID = storeID
		}

		var err error
		ids, err = sp.CreateMultipleStoreProducts(storeProducts)
		if err != nil {
			return err
		}

		storeWarehouseRepo := m.storeStockRepo.CloneWithTransaction(tx)
		if err := m.addStocks(&storeWarehouseRepo, storeID, dtos); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (m *transactionManager) UpdateStoreProductWithStocks(storeID, storeProductID uint, updateModels *types.StoreProductModels, dtos []storeWarehousesTypes.AddStoreStockDTO) error {
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

func (m *transactionManager) addStocks(storeWarehouseRepo storeStocks.StoreStockRepository, storeID uint, dtos []storeWarehousesTypes.AddStoreStockDTO) error {
	for _, dto := range dtos {
		_, err := storeWarehouseRepo.AddStock(storeID, &dto)
		if err != nil {
			switch {
			case errors.Is(err, storeWarehousesTypes.ErrStockAlreadyExists):
				continue
			}
			return err
		}
	}
	return nil
}
