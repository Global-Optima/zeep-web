package warehouse

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	AssignStoreToWarehouse(storeID, warehouseID uint) error
	ReassignStoreToWarehouse(storeID, newWarehouseID uint) error
	GetAllStoresByWarehouse(warehouseID uint, pagination *utils.Pagination) ([]data.Store, error)

	CreateWarehouse(warehouse *data.Warehouse, facilityAddress *data.FacilityAddress) error
	GetWarehouseByID(id uint) (*data.Warehouse, error)
	GetAllWarehouses(filter *types.WarehouseFilter) ([]data.Warehouse, error)
	GetAllWarehousesForNotifications() ([]data.Warehouse, error)
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
	return r.db.Model(&data.Store{}).
		Where("id = ?", storeID).
		Update("warehouse_id", warehouseID).Error
}

func (r *warehouseRepository) ReassignStoreToWarehouse(storeID, newWarehouseID uint) error {
	return r.db.Model(&data.Store{}).
		Where("id = ?", storeID).
		Update("warehouse_id", newWarehouseID).Error
}

func (r *warehouseRepository) GetAllStoresByWarehouse(warehouseID uint, pagination *utils.Pagination) ([]data.Store, error) {
	var stores []data.Store

	query := r.db.Preload("Warehouse").Where("warehouse_id = ?", warehouseID)

	if _, err := utils.ApplyPagination(query, pagination, &data.Store{}); err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	if err := query.Find(&stores).Error; err != nil {
		return nil, err
	}

	return stores, nil
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
	if err := r.db.
		Preload("FacilityAddress").
		Preload("Region").
		First(&warehouse, id).Error; err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) GetAllWarehouses(filter *types.WarehouseFilter) ([]data.Warehouse, error) {
	var warehouses []data.Warehouse

	query := r.db.Model(&data.Warehouse{}).
		Preload("FacilityAddress").
		Preload("Region")

	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	if filter.RegionID != nil {
		query = query.Where("region_id = ?", *filter.RegionID)
	}

	if filter.Search != nil {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ?", searchTerm)
	}

	if err := query.Scopes(filter.Sort.SortGorm()).Find(&warehouses).Error; err != nil {
		return nil, err
	}

	return warehouses, nil
}

func (r *warehouseRepository) GetAllWarehousesForNotifications() ([]data.Warehouse, error) {
	var warehouses []data.Warehouse

	query := r.db.Model(&data.Warehouse{}).
		Preload("FacilityAddress").
		Preload("Region")

	if err := query.Find(&warehouses).Error; err != nil {
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
