package storeAdditives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
	storeWarehousesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TransactionManager interface {
	CreateStoreAdditivesWithStocks(storeID uint, storeAdditive []data.StoreAdditive, dtos []storeWarehousesTypes.AddStockDTO) ([]uint, error)
}

type transactionManager struct {
	db                 *gorm.DB
	storeAdditiveRepo  StoreAdditiveRepository
	storeWarehouseRepo storeWarehouses.StoreWarehouseRepository
}

func NewTransactionManager(db *gorm.DB, storeAdditiveRepo StoreAdditiveRepository, storeWarehouseRepo storeWarehouses.StoreWarehouseRepository) TransactionManager {
	return &transactionManager{
		db:                 db,
		storeAdditiveRepo:  storeAdditiveRepo,
		storeWarehouseRepo: storeWarehouseRepo,
	}
}

func (m *transactionManager) CreateStoreAdditivesWithStocks(storeID uint, storeAdditives []data.StoreAdditive, dtos []storeWarehousesTypes.AddStockDTO) ([]uint, error) {
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

func (m *transactionManager) UpdateStoreAdditivesWithStocks(storeID, storeAdditiveID uint, updateStoreAdditive *data.StoreAdditive, dtos []storeWarehousesTypes.AddStockDTO) error {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		sp := m.storeAdditiveRepo.CloneWithTransaction(tx)
		if err := sp.UpdateStoreAdditive(storeID, storeAdditiveID, updateStoreAdditive); err != nil {
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
