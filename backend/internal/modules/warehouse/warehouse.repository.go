package warehouse

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	AssignStoreToWarehouse(storeID, warehouseID uint) error
	ReassignStoreToWarehouse(storeID, newWarehouseID uint) error
	GetAllStoresByWarehouse(warehouseID uint) ([]data.Store, error)

	CreateWarehouse(warehouse *data.Warehouse, facilityAddress *data.FacilityAddress) error
	GetWarehouseByID(id uint) (*data.Warehouse, error)
	GetAllWarehouses() ([]data.Warehouse, error)
	UpdateWarehouse(warehouse *data.Warehouse) error
	DeleteWarehouse(id uint) error
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

func (r *warehouseRepository) GetAllStoresByWarehouse(warehouseID uint) ([]data.Store, error) {
	var stores []data.Store
	err := r.db.Joins("JOIN store_warehouses ON stores.id = store_warehouses.store_id").
		Where("store_warehouses.warehouse_id = ?", warehouseID).
		Find(&stores).Error
	return stores, err
}

func (r *warehouseRepository) CreateWarehouse(warehouse *data.Warehouse, facilityAddress *data.FacilityAddress) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(facilityAddress).Error; err != nil {
			return err
		}
		warehouse.FacilityAddressID = facilityAddress.ID
		return tx.Create(warehouse).Error
	})
}

func (r *warehouseRepository) GetWarehouseByID(id uint) (*data.Warehouse, error) {
	var warehouse data.Warehouse
	if err := r.db.Preload("FacilityAddress").First(&warehouse, id).Error; err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) GetAllWarehouses() ([]data.Warehouse, error) {
	var warehouses []data.Warehouse
	if err := r.db.Preload("FacilityAddress").Find(&warehouses).Error; err != nil {
		return nil, err
	}
	return warehouses, nil
}

func (r *warehouseRepository) UpdateWarehouse(warehouse *data.Warehouse) error {
	return r.db.Save(warehouse).Error
}

func (r *warehouseRepository) DeleteWarehouse(id uint) error {
	return r.db.Delete(&data.Warehouse{}, id).Error
}
