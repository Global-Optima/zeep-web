package employees

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/storeEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type StoreEmployeeRepository interface {
	GetStoreEmployees(storeID uint, filter *employeesTypes.EmployeesFilter) ([]data.StoreEmployee, error)
	GetStoreEmployeeByID(id uint, filter *contexts.StoreContextFilter) (*data.StoreEmployee, error)
	GetAllStoreEmployees(storeID uint) ([]data.StoreEmployee, error)
	UpdateStoreEmployee(id uint, filter *contexts.StoreContextFilter, updateModels *types.UpdateStoreEmployeeModels) error
}

type storeEmployeeRepository struct {
	db           *gorm.DB
	employeeRepo employees.EmployeeRepository
}

func NewStoreEmployeeRepository(db *gorm.DB, employeeRepo employees.EmployeeRepository) StoreEmployeeRepository {
	return &storeEmployeeRepository{
		db:           db,
		employeeRepo: employeeRepo,
	}
}

func (r *storeEmployeeRepository) GetStoreEmployees(storeID uint, filter *employeesTypes.EmployeesFilter) ([]data.StoreEmployee, error) {
	var storeEmployees []data.StoreEmployee
	query := r.db.Model(&data.StoreEmployee{}).
		Where("store_id = ?", storeID).
		Preload("Employee").
		Joins("JOIN employees ON employees.id = store_employees.employee_id")

	if filter.IsActive != nil {
		query = query.Where("employees.is_active = ?", *filter.IsActive)
	}

	if filter.Role != nil {
		query = query.Where("role = ?", *filter.Role)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.phone ILIKE ? OR employees.email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreEmployee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&storeEmployees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve store employees: %w", err)
	}
	return storeEmployees, nil
}

func (r *storeEmployeeRepository) GetStoreEmployeeByID(id uint, filter *contexts.StoreContextFilter) (*data.StoreEmployee, error) {
	var storeEmployee data.StoreEmployee
	query := r.db.Model(&data.StoreEmployee{}).
		Preload("Employee.Workdays").
		Preload("Store.FacilityAddress").
		Preload("Store.Warehouse.FacilityAddress").
		Preload("Store.Warehouse.Region").
		Where(&data.StoreEmployee{BaseEntity: data.BaseEntity{ID: id}}, id)

	if filter != nil {
		if filter.StoreID != nil {
			query = query.Where(&data.StoreEmployee{StoreID: *filter.StoreID})
		}
		if filter.FranchiseeID != nil {
			query = query.Joins("JOIN stores ON stores.id = store_employees.store_id").
				Where(&data.StoreEmployee{Store: data.Store{FranchiseeID: filter.FranchiseeID}})
		}
	}

	err := query.First(&storeEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve store employee by ID: %w", err)
	}
	return &storeEmployee, nil
}

func (r *storeEmployeeRepository) UpdateStoreEmployee(id uint, filter *contexts.StoreContextFilter, updateModels *types.UpdateStoreEmployeeModels) error {
	if updateModels == nil {
		return employeesTypes.ErrNothingToUpdate
	}

	existingStoreEmployee, err := r.GetStoreEmployeeByID(id, filter)
	if err != nil {
		return err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels.StoreEmployee != nil {
			err := tx.Model(&data.StoreEmployee{}).
				Where("id = ?", id).
				Updates(updateModels.StoreEmployee).Error
			if err != nil {
				return err
			}
		}

		if updateModels.UpdateEmployeeModels != nil {
			err := r.employeeRepo.UpdateEmployeeWithAssociations(tx, existingStoreEmployee.EmployeeID, updateModels.UpdateEmployeeModels)
			if err != nil {
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

func (r *storeEmployeeRepository) GetAllStoreEmployees(storeID uint) ([]data.StoreEmployee, error) {
	var storeEmployees []data.StoreEmployee

	err := r.db.Model(&data.StoreEmployee{}).
		Preload("Employee").
		Joins("INNER JOIN employees ON store_employees.employee_id = employees.id").
		Where("store_employees.store_id = ? AND employees.is_active = true", storeID).
		Find(&storeEmployees).Error
	if err != nil {
		return nil, err
	}

	return storeEmployees, nil
}
