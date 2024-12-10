package warehouse

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	AssignStoreToWarehouse(storeID, warehouseID uint) error
	ReassignStoreToWarehouse(storeID, newWarehouseID uint) error
	ListStoresForWarehouse(warehouseID uint) ([]data.Store, error)
}

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) AssignStoreToWarehouse(storeID, warehouseID uint) error {
	storeWarehouse := data.StoreWarehouse{
		StoreID:     storeID,
		WarehouseID: warehouseID,
	}
	return r.db.Create(&storeWarehouse).Error
}

func (r *warehouseRepository) ReassignStoreToWarehouse(storeID, newWarehouseID uint) error {
	return r.db.Model(&data.StoreWarehouse{}).
		Where("store_id = ?", storeID).
		Update("warehouse_id", newWarehouseID).Error
}

func (r *warehouseRepository) ListStoresForWarehouse(warehouseID uint) ([]data.Store, error) {
	var stores []data.Store
	err := r.db.Joins("JOIN store_warehouses ON stores.id = store_warehouses.store_id").
		Where("store_warehouses.warehouse_id = ?", warehouseID).
		Find(&stores).Error
	return stores, err
}
