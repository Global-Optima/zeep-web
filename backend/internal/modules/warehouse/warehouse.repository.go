package warehouse

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	AssignStoreToWarehouse(storeID, warehouseID uint) error
	GetAllStoresByWarehouse(warehouseID uint, pagination *utils.Pagination) ([]data.Store, error)

	CreateWarehouse(warehouse *data.Warehouse, facilityAddress *data.FacilityAddress) error
	GetWarehouseByID(id uint) (*data.Warehouse, error)
	GetWarehouses(filter *types.WarehouseFilter) ([]data.Warehouse, error)
	GetAllWarehouses(filter *types.WarehouseFilter) ([]data.Warehouse, error)
	GetAllWarehousesForNotifications() ([]data.Warehouse, error)
	UpdateWarehouse(id uint, updateModels *types.WarehouseUpdateModels) error
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
	err := r.db.
		Preload("FacilityAddress").
		Preload("Region").
		First(&warehouse, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrWarehouseNotFound
		}
		return nil, types.ErrFailedToFetchWarehouse
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

	if filter.Search != nil {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ?", searchTerm)
	}

	if err := query.Scopes(filter.Sort.SortGorm()).Find(&warehouses).Error; err != nil {
		return nil, err
	}

	return warehouses, nil
}

func (r *warehouseRepository) GetWarehouses(filter *types.WarehouseFilter) ([]data.Warehouse, error) {
	var warehouses []data.Warehouse

	query := r.db.Model(&data.Warehouse{}).
		Preload("FacilityAddress").
		Preload("Region")

	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	if filter.RegionID != nil {
		query.Where("region_id = ?", *filter.RegionID)
	}

	if filter.Search != nil {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ?", searchTerm)
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Warehouse{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&warehouses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrWarehouseNotFound
		}
		return nil, types.ErrFailedToFetchWarehouse
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

func (r *warehouseRepository) UpdateWarehouse(id uint, updateModels *types.WarehouseUpdateModels) error {
	if updateModels == nil {
		return fmt.Errorf("nothing to update")
	}

	existingWarehouse, err := r.GetWarehouseByID(id)
	if err != nil {
		return err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels.Warehouse != nil {
			query := tx.Model(&data.Warehouse{}).Where(&data.Warehouse{BaseEntity: data.BaseEntity{ID: id}})

			if err := query.Updates(updateModels.Warehouse).Error; err != nil {
				return err
			}
		}

		if updateModels.FacilityAddress != nil {
			if err := tx.Where(&data.FacilityAddress{BaseEntity: data.BaseEntity{ID: existingWarehouse.FacilityAddress.ID}}).
				Updates(updateModels.FacilityAddress).Error; err != nil {
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

func (r *warehouseRepository) DeleteWarehouse(id uint) error {
	result := r.db.Delete(&data.Warehouse{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return types.ErrWarehouseNotFound
	}
	return nil
}
