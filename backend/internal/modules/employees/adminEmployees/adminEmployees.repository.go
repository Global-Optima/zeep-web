package employees

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type AdminEmployeeRepository interface {
	GetAdminEmployees(filter *types.EmployeesFilter) ([]data.AdminEmployee, error)
	GetAdminEmployeeByID(id uint) (*data.AdminEmployee, error)
	GetAllAdminEmployees() ([]data.AdminEmployee, error)
}

type adminEmployeeRepository struct {
	db *gorm.DB
}

func NewAdminEmployeeRepository(db *gorm.DB) AdminEmployeeRepository {
	return &adminEmployeeRepository{db: db}
}

func (r *adminEmployeeRepository) GetAdminEmployees(filter *types.EmployeesFilter) ([]data.AdminEmployee, error) {
	var employees []data.AdminEmployee
	query := r.db.Model(&data.AdminEmployee{}).
		Preload("Employee").
		Joins("JOIN employees ON employees.id = admin_employees.employee_id")

	if filter.IsActive != nil {
		query = query.Where(&data.AdminEmployee{
			Employee: data.Employee{IsActive: filter.IsActive},
		})
	}

	if filter.Role != nil {
		query = query.Where(&data.AdminEmployee{Role: *filter.Role})
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.phone ILIKE ? OR employees.email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.AdminEmployee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin employees: %w", err)
	}
	return employees, nil
}

func (r *adminEmployeeRepository) GetAdminEmployeeByID(id uint) (*data.AdminEmployee, error) {
	var adminEmployee data.AdminEmployee
	err := r.db.Model(&data.AdminEmployee{}).
		Preload("Employee.Workdays").
		Where(&data.AdminEmployee{
			BaseEntity: data.BaseEntity{ID: id},
		}).
		First(&adminEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin employee by ID: %w", err)
	}
	return &adminEmployee, nil
}

func (r *adminEmployeeRepository) GetAllAdminEmployees() ([]data.AdminEmployee, error) {
	var employees []data.AdminEmployee

	err := r.db.Model(&data.AdminEmployee{}).
		Preload("Employee").
		Joins("INNER JOIN employees ON admin_employees.employee_id = employees.id").
		Where("employees.is_active = true").
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}
