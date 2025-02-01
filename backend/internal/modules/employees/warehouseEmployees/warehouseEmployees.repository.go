package employees

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type WarehouseEmployeeRepository interface {
	GetWarehouseEmployees(warehouseID uint, filter *types.EmployeesFilter) ([]data.WarehouseEmployee, error)
	GetWarehouseEmployeeByID(id, warehouseID uint) (*data.WarehouseEmployee, error)
	UpdateWarehouseEmployee(id uint, warehouseID uint, update *data.WarehouseEmployee) error
}

type warehouseEmployeeRepository struct {
	db *gorm.DB
}

func NewWarehouseEmployeeRepository(db *gorm.DB) WarehouseEmployeeRepository {
	return &warehouseEmployeeRepository{db: db}
}

func (r *warehouseEmployeeRepository) GetWarehouseEmployees(warehouseID uint, filter *types.EmployeesFilter) ([]data.WarehouseEmployee, error) {
	var employees []data.WarehouseEmployee
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

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve warehouse employees: %w", err)
	}
	return employees, nil
}

func (r *warehouseEmployeeRepository) GetWarehouseEmployeeByID(id, warehouseID uint) (*data.WarehouseEmployee, error) {
	var warehouseEmployee data.WarehouseEmployee
	err := r.db.Model(&data.WarehouseEmployee{}).
		Preload("Employee").
		Where("id = ? AND warehouse_id = ?", id, warehouseID).
		First(&warehouseEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve warehouse employee by ID: %w", err)
	}
	return &warehouseEmployee, nil
}

func (r *warehouseEmployeeRepository) UpdateWarehouseEmployee(id uint, warehouseID uint, warehouseEmployee *data.WarehouseEmployee) error {
	if warehouseEmployee == nil {
		return types.ErrNothingToUpdate
	}
	return r.db.Model(&data.WarehouseEmployee{}).
		Where("id = ? AND warehouse_id = ?", id, warehouseID).
		Updates(warehouseEmployee).Error
}
