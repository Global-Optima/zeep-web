package stores

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
	"gorm.io/gorm"
)

const (
	ACTIVE_STORE_STATUS   = "ACTIVE"
	DISABLED_STORE_STATUS = "DISABLED"
)

type StoreRepository interface {
	GetAllStores(filter *types.StoreFilter) ([]data.Store, error)
	GetAllStoresForNotifications() ([]data.Store, error)
	CreateStore(store *data.Store) (uint, error)
	GetStoreByID(storeID uint) (*data.Store, error)
	GetRawStoreByID(storeID uint) (*data.Store, error)
	GetStores(filter *types.StoreFilter) ([]data.Store, error)
	UpdateStore(storeID uint, updateModels *types.StoreUpdateModels) error
	UpdateStoreSyncTime(storeID uint) error
	DeleteStore(storeID uint, hardDelete bool) error
	CreateFacilityAddress(facilityAddress *data.FacilityAddress) (*data.FacilityAddress, error)
	GetFacilityAddressByAddress(address string) (*data.FacilityAddress, error)
	CloneWithTransaction(tx *gorm.DB) storeRepository
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) CloneWithTransaction(tx *gorm.DB) storeRepository {
	return storeRepository{
		db: tx,
	}
}

func (r *storeRepository) GetAllStores(filter *types.StoreFilter) ([]data.Store, error) {
	var stores []data.Store

	query := r.db.Preload("FacilityAddress").
		Preload("Franchisee").
		Preload("Warehouse").
		Preload("Warehouse.Region").
		Preload("Warehouse.FacilityAddress").
		Preload("FacilityAddress")

	if filter == nil {
		return nil, errors.New("filter is nil")
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ? OR contact_phone ILIKE = ? OR contact_email ILIKE = ?",
			searchTerm, searchTerm, searchTerm)
	}

	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", filter.WarehouseID)
	}

	if filter.IsFranchisee != nil {
		if *filter.IsFranchisee {
			query = query.Where("franchisee_id IS NOT NULL")
		} else {
			query = query.Where("franchisee_id IS NULL")
		}
	}

	if err := query.Scopes(filter.Sort.SortGorm()).Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *storeRepository) GetAllStoresForNotifications() ([]data.Store, error) {
	var stores []data.Store

	query := r.db.Preload("FacilityAddress").
		Preload("FacilityAddress").
		Preload("Franchisee").
		Preload("Warehouse").
		Preload("Warehouse.Region").
		Preload("Warehouse.FacilityAddress")

	if err := query.Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *storeRepository) CreateStore(store *data.Store) (uint, error) {
	if err := r.db.Create(store).Error; err != nil {
		return 0, err
	}
	return store.ID, nil
}

func (r *storeRepository) GetRawStoreByID(storeID uint) (*data.Store, error) {
	var store data.Store
	if err := r.db.Where(&data.Store{BaseEntity: data.BaseEntity{ID: storeID}}).
		First(&store).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) GetStoreByID(storeID uint) (*data.Store, error) {
	var store data.Store
	if err := r.db.Preload("FacilityAddress").
		Preload("Franchisee").
		Preload("Warehouse").
		Preload("Warehouse.Region").
		Preload("Warehouse.FacilityAddress").
		Where("id = ?", storeID).First(&store).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) GetStores(filter *types.StoreFilter) ([]data.Store, error) {
	var stores []data.Store
	query := r.db.Model(&data.Store{}).
		Preload("FacilityAddress").
		Preload("Franchisee").
		Preload("Warehouse").
		Preload("Warehouse.Region").
		Preload("Warehouse.FacilityAddress")

	if filter == nil {
		return nil, errors.New("filter is nil")
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ? OR contact_phone ILIKE ? OR contact_email ILIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	if filter.IsFranchisee != nil {
		if *filter.IsFranchisee {
			query = query.Where("franchisee_id IS NOT NULL")
		} else {
			query = query.Where("franchisee_id IS NULL")
		}
	}

	if filter.FranchiseeID != nil {
		query = query.Where("franchisee_id = ?", *filter.FranchiseeID)
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Store{})
	if err != nil {
		return nil, err
	}
	query.Find(&stores)

	return stores, nil
}

func (r *storeRepository) UpdateStore(storeID uint, updateModels *types.StoreUpdateModels) error {
	if updateModels == nil {
		return fmt.Errorf("nothing to update")
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&updateModels.Store).Error; err != nil {
			return err
		}

		if updateModels.FacilityAddress != nil {
			if err := tx.Save(&updateModels.FacilityAddress).Error; err != nil {
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

func (r *storeRepository) DeleteStore(storeID uint, hardDelete bool) error {
	if hardDelete {
		if err := r.db.Delete(&data.Store{}, storeID).Error; err != nil {
			return err
		}
	} else {
		if err := r.db.Model(&data.Store{}).Where("id = ?", storeID).Update("status", DISABLED_STORE_STATUS).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *storeRepository) CreateFacilityAddress(facilityAddress *data.FacilityAddress) (*data.FacilityAddress, error) {
	if err := r.db.Create(facilityAddress).Error; err != nil {
		return nil, err
	}
	return facilityAddress, nil
}

func (r *storeRepository) GetFacilityAddressByAddress(address string) (*data.FacilityAddress, error) {
	var facilityAddress data.FacilityAddress
	if err := r.db.Where("trim(lower(address)) = ?", strings.ToLower(strings.TrimSpace(address))).First(&facilityAddress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &facilityAddress, nil
}

func (r *storeRepository) UpdateStoreSyncTime(storeID uint) error {
	updatedStore := &data.Store{LastInventorySyncAt: time.Now()}

	return r.db.Where(&data.Store{BaseEntity: data.BaseEntity{ID: storeID}}).
		Updates(updatedStore).Error
}
