package employees

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/warehouseEmployees/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type WarehouseEmployeeRepository interface {
	GetWarehouseEmployees(warehouseID uint, filter *employeesTypes.EmployeesFilter) ([]data.WarehouseEmployee, error)
	GetWarehouseEmployeeByID(id uint, filter *contexts.WarehouseContextFilter) (*data.WarehouseEmployee, error)
	GetAllWarehouseEmployees(warehouseID uint) ([]data.WarehouseEmployee, error)
	UpdateWarehouseEmployee(id uint, filter *contexts.WarehouseContextFilter, update *types.UpdateWarehouseEmployeeModels) error
}

type warehouseEmployeeRepository struct {
	db           *gorm.DB
	employeeRepo employees.EmployeeRepository
}

func NewWarehouseEmployeeRepository(db *gorm.DB, warehouseRepo employees.EmployeeRepository) WarehouseEmployeeRepository {
	return &warehouseEmployeeRepository{
		db:           db,
		employeeRepo: warehouseRepo,
	}
}

func (r *warehouseEmployeeRepository) GetWarehouseEmployees(warehouseID uint, filter *employeesTypes.EmployeesFilter) ([]data.WarehouseEmployee, error) {
	var warehouseEmployees []data.WarehouseEmployee
	query := r.db.Model(&data.WarehouseEmployee{}).
		Where(&data.WarehouseEmployee{WarehouseID: warehouseID}).
		Preload("Employee").
		Joins("JOIN employees ON employees.id = warehouse_employees.employee_id")

	if filter.IsActive != nil {
		query = query.Where(&data.WarehouseEmployee{Employee: data.Employee{IsActive: filter.IsActive}})
	}

	if filter.Role != nil {
		query = query.Where(&data.WarehouseEmployee{Role: *filter.Role})
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.phone ILIKE ? OR employees.email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.WarehouseEmployee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&warehouseEmployees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve warehouse employees: %w", err)
	}
	return warehouseEmployees, nil
}

func (r *warehouseEmployeeRepository) GetWarehouseEmployeeByID(id uint, filter *contexts.WarehouseContextFilter) (*data.WarehouseEmployee, error) {
	var warehouseEmployee data.WarehouseEmployee
	query := r.db.Model(&data.WarehouseEmployee{}).
		Preload("Employee.Workdays").
		Preload("Warehouse.FacilityAddress").
		Preload("Warehouse.Region").
		Where(&data.WarehouseEmployee{BaseEntity: data.BaseEntity{ID: id}})

	if filter != nil {
		if filter.WarehouseID != nil {
			query.Where(&data.WarehouseEmployee{WarehouseID: *filter.WarehouseID})
		}

		if filter.RegionID != nil {
			query.Joins("JOIN warehouses ON warehouses.id = warehouse_employees.warehouse_id")
			query.Where(&data.WarehouseEmployee{Warehouse: data.Warehouse{RegionID: *filter.RegionID}})
		}
	}

	if err := query.First(&warehouseEmployee).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve warehouse employee by ID: %w", err)
	}

	return &warehouseEmployee, nil
}

func (r *warehouseEmployeeRepository) GetAllWarehouseEmployees(warehouseID uint) ([]data.WarehouseEmployee, error) {
	var warehouseEmployees []data.WarehouseEmployee

	err := r.db.Model(&data.WarehouseEmployee{}).
		Preload("Employee").
		Joins("INNER JOIN employees ON warehouse_employees.employee_id = employees.id").
		Where("warehouse_employees.warehouse_id = ? AND employees.is_active = true", warehouseID).
		Find(&warehouseEmployees).Error

	if err != nil {
		return nil, err
	}

	return warehouseEmployees, nil
}

func (r *warehouseEmployeeRepository) UpdateWarehouseEmployee(id uint, filter *contexts.WarehouseContextFilter, updateModels *types.UpdateWarehouseEmployeeModels) error {
	if updateModels == nil {
		return employeesTypes.ErrNothingToUpdate
	}

	existingWarehouseEmployee, err := r.GetWarehouseEmployeeByID(id, filter)
	if err != nil {
		return err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels.WarehouseEmployee != nil && !utils.IsEmpty(updateModels.WarehouseEmployee) {
			err := tx.Model(&data.WarehouseEmployee{}).
				Where(&data.WarehouseEmployee{BaseEntity: data.BaseEntity{ID: id}}).
				Updates(updateModels.WarehouseEmployee).Error
			if err != nil {
				return err
			}
		}

		if updateModels.UpdateEmployeeModels != nil && !utils.IsEmpty(updateModels.UpdateEmployeeModels) {
			err := r.employeeRepo.UpdateEmployeeWithAssociations(tx, existingWarehouseEmployee.EmployeeID, updateModels.UpdateEmployeeModels)
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
