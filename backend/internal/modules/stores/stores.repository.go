package stores

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type StoreRepository interface {
	GetAllStores(searchTerm string) ([]data.Store, error)
	CreateStore(store *data.Store) (*data.Store, error)
	GetStoreByID(storeID uint) (*data.Store, error)
	UpdateStore(store *data.Store) (*data.Store, error)
	DeleteStore(storeID uint, hardDelete bool) error
	CreateFacilityAddress(facilityAddress *data.FacilityAddress) (*data.FacilityAddress, error)
	GetFacilityAddressByAddress(address string) (*data.FacilityAddress, error)
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) GetAllStores(searchTerm string) ([]data.Store, error) {
	var stores []data.Store

	query := r.db.Preload("FacilityAddress").Where("status = ?", "active").Preload("FacilityAddress")

	if searchTerm != "" {
		query = query.Where("name ILIKE ? OR CAST(id AS TEXT) = ?", "%"+searchTerm+"%", searchTerm)
	}

	if err := query.Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *storeRepository) CreateStore(store *data.Store) (*data.Store, error) {
	if err := r.db.Create(store).Error; err != nil {
		return nil, err
	}
	return store, nil
}

func (r *storeRepository) GetStoreByID(storeID uint) (*data.Store, error) {
	var store data.Store
	if err := r.db.Preload("FacilityAddress").Where("id = ?", storeID).First(&store).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) UpdateStore(store *data.Store) (*data.Store, error) {
	if err := r.db.Save(store).Error; err != nil {
		return nil, err
	}
	return store, nil
}

func (r *storeRepository) DeleteStore(storeID uint, hardDelete bool) error {
	if hardDelete {
		if err := r.db.Delete(&data.Store{}, storeID).Error; err != nil {
			return err
		}
	} else {
		if err := r.db.Model(&data.Store{}).Where("id = ?", storeID).Update("status", "inactive").Error; err != nil {
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
	if err := r.db.Where("address = ?", address).First(&facilityAddress).Error; err != nil {
		return nil, err
	}
	return &facilityAddress, nil
}
