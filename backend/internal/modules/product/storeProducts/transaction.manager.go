package storeProducts

import (
	"errors"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
	storeWarehousesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionManager interface {
	CreateStoreProductWithStocks(storeID uint, storeProduct *data.StoreProduct, dtos []storeWarehousesTypes.AddStockDTO) (uint, error)
	CreateMultipleStoreProductsWithStocks(storeID uint, storeProducts []data.StoreProduct, dtos []storeWarehousesTypes.AddStockDTO) ([]uint, error)
	UpdateStoreProductWithStocks(storeID, storeProductID uint, updateModels *types.StoreProductModels, dtos []storeWarehousesTypes.AddStockDTO) error
}

type transactionManager struct {
	db                 *gorm.DB
	storeProductRepo   StoreProductRepository
	storeWarehouseRepo storeWarehouses.StoreWarehouseRepository
}

func NewTransactionManager(db *gorm.DB, storeProductRepo StoreProductRepository, storeWarehouseRepo storeWarehouses.StoreWarehouseRepository) TransactionManager {
	return &transactionManager{
		db:                 db,
		storeProductRepo:   storeProductRepo,
		storeWarehouseRepo: storeWarehouseRepo,
	}
}

func (m *transactionManager) CreateStoreProductWithStocks(storeID uint, storeProduct *data.StoreProduct, dtos []storeWarehousesTypes.AddStockDTO) (uint, error) {
	var id uint
	err := m.db.Transaction(func(tx *gorm.DB) error {
		var err error
		sp := m.storeProductRepo.CloneWithTransaction(tx)
		storeProduct.StoreID = storeID
		id, err = sp.CreateStoreProduct(storeProduct)
		if err != nil {
			return err
		}

		sw := m.storeWarehouseRepo.CloneWithTransaction(tx)
		logrus.Info(dtos)
		for _, dto := range dtos {
			logrus.Info(dto)
			_, err := sw.AddStock(storeID, &dto)

			if err != nil {
				switch {
				case errors.Is(err, storeWarehousesTypes.ErrStockAlreadyExists):
					continue
				}
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

func (m *transactionManager) CreateMultipleStoreProductsWithStocks(storeID uint, storeProducts []data.StoreProduct, dtos []storeWarehousesTypes.AddStockDTO) ([]uint, error) {
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

		sw := m.storeWarehouseRepo.CloneWithTransaction(tx)
		for _, dto := range dtos {
			_, err := sw.AddStock(storeID, &dto)

			if err != nil {
				switch {
				case errors.Is(err, storeWarehousesTypes.ErrStockAlreadyExists):
					continue
				}
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (m *transactionManager) UpdateStoreProductWithStocks(storeID, storeProductID uint, updateModels *types.StoreProductModels, dtos []storeWarehousesTypes.AddStockDTO) error {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		sp := m.storeProductRepo.CloneWithTransaction(tx)
		if err := sp.UpdateStoreProductByID(storeID, storeProductID, updateModels); err != nil {
			return err
		}
		sw := m.storeWarehouseRepo.CloneWithTransaction(tx)
		for _, dto := range dtos {
			_, err := sw.AddStock(storeID, &dto)
			if err != nil {
				switch {
				case errors.Is(err, storeWarehousesTypes.ErrStockAlreadyExists):
					continue
				}
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
