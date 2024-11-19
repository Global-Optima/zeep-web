package stores

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type StoreRepository interface {
	GetAllStores() ([]data.Store, error)
	GetStoreEmployees(storeID uint) ([]data.Employee, error)
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) GetAllStores() ([]data.Store, error) {
	var stores []data.Store
	err := r.db.Preload("FacilityAddress").Find(&stores).Error
	if err != nil {
		return nil, err
	}

	return stores, nil
}

func (r *storeRepository) GetStoreEmployees(storeID uint) ([]data.Employee, error) {
	var employees []data.Employee
	err := r.db.Where(&data.Employee{StoreID: &storeID}).Preload("Role").Find(&employees).Error
	if err != nil {
		return nil, err
	}

	return employees, nil
}
