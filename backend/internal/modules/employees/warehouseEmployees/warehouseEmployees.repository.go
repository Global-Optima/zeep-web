package employees

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/warehouseEmployees/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type WarehouseEmployeeRepository interface {
	GetWarehouseEmployees(warehouseID uint, filter *employeesTypes.EmployeesFilter) ([]data.WarehouseEmployee, error)
	GetWarehouseEmployeeByID(id, warehouseID uint) (*data.WarehouseEmployee, error)
	UpdateWarehouseEmployee(id uint, warehouseID uint, update *types.UpdateWarehouseEmployeeModels) error
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
		Where("warehouse_id = ?", warehouseID).Preload("Employee")

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Role != nil {
		query = query.Where("role = ?", *filter.Role)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR phone ILIKE ? OR email ILIKE ?",
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

func (r *warehouseEmployeeRepository) GetWarehouseEmployeeByID(id, warehouseID uint) (*data.WarehouseEmployee, error) {
	var warehouseEmployee data.WarehouseEmployee
	err := r.db.Model(&data.WarehouseEmployee{}).
		Preload("Employee.Workdays").
		Preload("Warehouse.FacilityAddress").
		Where("id = ? AND warehouse_id = ?", id, warehouseID).
		First(&warehouseEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve warehouse employee by ID: %w", err)
	}
	return &warehouseEmployee, nil
}

func (r *warehouseEmployeeRepository) UpdateWarehouseEmployee(id uint, warehouseID uint, updateModels *types.UpdateWarehouseEmployeeModels) error {
	if updateModels == nil {
		return employeesTypes.ErrNothingToUpdate
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existingWarehouseEmployee data.WarehouseEmployee
		r.db.Model(&data.WarehouseEmployee{}).
			Where("id = ? AND warehouse_id = ?", id, warehouseID).
			First(&existingWarehouseEmployee)

		if updateModels.WarehouseEmployee != nil && !utils.IsEmpty(updateModels.WarehouseEmployee) {
			err := tx.Model(&data.WarehouseEmployee{}).
				Where("id = ? AND warehouse_id = ?", id, warehouseID).
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
