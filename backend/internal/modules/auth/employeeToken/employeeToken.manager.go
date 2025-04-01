package employeeToken

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type EmployeeTokenManager interface {
	CreateToken(token *data.EmployeeToken) error
	GetTokenByEmployeeID(employeeID uint) (*data.EmployeeToken, error)
	DeleteToken(token *data.EmployeeToken) error
	DeleteTokenByEmployeeID(employeeID uint) error
	DeleteTokenByStoreEmployeeID(storeEmployeeID uint) error
	DeleteTokenByWarehouseEmployeeID(warehouseEmployeeID uint) error
	DeleteTokenByRegionEmployeeID(regionEmployeeID uint) error
	DeleteTokenByFranchiseeEmployeeID(franchiseeEmployeeID uint) error
}

type employeeTokenManager struct {
	db *gorm.DB
}

func NewEmployeeTokenManager(db *gorm.DB) EmployeeTokenManager {
	return &employeeTokenManager{db: db}
}

func (r *employeeTokenManager) CreateToken(token *data.EmployeeToken) error {
	return r.db.Create(token).Error
}

func (r *employeeTokenManager) GetTokenByEmployeeID(employeeID uint) (*data.EmployeeToken, error) {
	var token data.EmployeeToken
	err := r.db.Model(&data.EmployeeToken{}).
		Preload("Employee", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).
		Preload("Employee.StoreEmployee", func(db *gorm.DB) *gorm.DB {
			return db.Select("employee_id, store_id, role")
		}).
		Preload("Employee.WarehouseEmployee", func(db *gorm.DB) *gorm.DB {
			return db.Select("employee_id, warehouse_id, role")
		}).
		Preload("Employee.RegionEmployee", func(db *gorm.DB) *gorm.DB {
			return db.Select("employee_id, region_id, role")
		}).
		Preload("Employee.FranchiseeEmployee", func(db *gorm.DB) *gorm.DB {
			return db.Select("employee_id, franchisee_id, role")
		}).
		Preload("Employee.AdminEmployee", func(db *gorm.DB) *gorm.DB {
			return db.Select("employee_id, role")
		}).
		Where("employee_id = ?", employeeID).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

func (r *employeeTokenManager) DeleteToken(token *data.EmployeeToken) error {
	return r.db.Unscoped().Delete(token).Error
}

func (r *employeeTokenManager) DeleteTokenByEmployeeID(employeeID uint) error {
	return r.db.Unscoped().Where("employee_id = ?", employeeID).Delete(&data.EmployeeToken{}).Error
}

func (r *employeeTokenManager) DeleteTokenByStoreEmployeeID(storeEmployeeID uint) error {
	var storeEmp data.StoreEmployee
	if err := r.db.
		Where("id = ?", storeEmployeeID).
		Preload("Employee").
		First(&storeEmp).Error; err != nil {
		return err
	}
	return r.DeleteTokenByEmployeeID(storeEmp.EmployeeID)
}

func (r *employeeTokenManager) DeleteTokenByWarehouseEmployeeID(warehouseEmployeeID uint) error {
	var whEmp data.WarehouseEmployee
	if err := r.db.
		Where("id = ?", warehouseEmployeeID).
		Preload("Employee").
		First(&whEmp).Error; err != nil {
		return err
	}
	return r.DeleteTokenByEmployeeID(whEmp.EmployeeID)
}

func (r *employeeTokenManager) DeleteTokenByRegionEmployeeID(regionEmployeeID uint) error {
	var regionEmp data.RegionEmployee
	if err := r.db.
		Where("id = ?", regionEmployeeID).
		Preload("Employee").
		First(&regionEmp).Error; err != nil {
		return err
	}
	return r.DeleteTokenByEmployeeID(regionEmp.EmployeeID)
}

func (r *employeeTokenManager) DeleteTokenByFranchiseeEmployeeID(franchiseeEmployeeID uint) error {
	var franchiseeEmp data.FranchiseeEmployee
	if err := r.db.
		Where("id = ?", franchiseeEmployeeID).
		Preload("Employee").
		First(&franchiseeEmp).Error; err != nil {
		return err
	}
	return r.DeleteTokenByEmployeeID(franchiseeEmp.EmployeeID)
}
