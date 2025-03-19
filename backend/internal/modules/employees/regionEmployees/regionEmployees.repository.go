package employees

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/regionEmployees/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type RegionEmployeeRepository interface {
	GetRegionEmployees(regionID uint, filter *employeesTypes.EmployeesFilter) ([]data.RegionEmployee, error)
	GetRegionEmployeeByID(id uint, regionID *uint) (*data.RegionEmployee, error)
	GetAllRegionEmployees(regionID uint) ([]data.RegionEmployee, error)
	UpdateRegionEmployee(id uint, regionID *uint, updateModels *types.UpdateRegionEmployeeModels) error
}

type regionEmployeeRepository struct {
	db           *gorm.DB
	employeeRepo employees.EmployeeRepository
}

func NewRegionEmployeeRepository(db *gorm.DB, employeeRepo employees.EmployeeRepository) RegionEmployeeRepository {
	return &regionEmployeeRepository{
		db:           db,
		employeeRepo: employeeRepo,
	}
}

func (r *regionEmployeeRepository) GetRegionEmployees(regionID uint, filter *employeesTypes.EmployeesFilter) ([]data.RegionEmployee, error) {
	var regionEmployees []data.RegionEmployee
	query := r.db.Model(&data.RegionEmployee{}).
		Where(&data.RegionEmployee{
			RegionID: regionID,
		}).
		Preload("Employee").
		Joins("JOIN employees ON employees.id = region_employees.employee_id")

	if filter.IsActive != nil {
		query = query.Where("employees.is_active = ?", *filter.IsActive)
	}

	if filter.Role != nil {
		query = query.Where(&data.RegionEmployee{
			Role: *filter.Role,
		})
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.phone ILIKE ? OR employees.email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.RegionEmployee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&regionEmployees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve region managers: %w", err)
	}
	return regionEmployees, nil
}

func (r *regionEmployeeRepository) GetRegionEmployeeByID(id uint, regionID *uint) (*data.RegionEmployee, error) {
	var regionEmployee data.RegionEmployee
	query := r.db.Model(&data.RegionEmployee{}).
		Preload("Employee.Workdays").
		Preload("Region").
		Where(&data.RegionEmployee{
			BaseEntity: data.BaseEntity{ID: id},
		})

	if regionID != nil {
		query.Where(&data.RegionEmployee{
			RegionID: *regionID,
		})
	}
	if err := query.First(&regionEmployee).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve region manager by ID: %w", err)
	}

	return &regionEmployee, nil
}

func (r *regionEmployeeRepository) GetAllRegionEmployees(regionID uint) ([]data.RegionEmployee, error) {
	var regionEmployees []data.RegionEmployee

	err := r.db.Model(&data.RegionEmployee{}).
		Joins("INNER JOIN employees ON region_employees.employee_id = employees.id").
		Preload("Employee").
		Where("region_employees.region_id = ? AND employees.is_active = true", regionID).
		Find(&regionEmployees).Error
	if err != nil {
		return nil, err
	}

	return regionEmployees, nil
}

func (r *regionEmployeeRepository) UpdateRegionEmployee(id uint, regionID *uint, updateModels *types.UpdateRegionEmployeeModels) error {
	if updateModels == nil {
		return employeesTypes.ErrNothingToUpdate
	}

	existingRegionEmployee, err := r.GetRegionEmployeeByID(id, regionID)
	if err != nil {
		return err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels.RegionEmployee != nil && !utils.IsEmpty(updateModels.RegionEmployee) {
			err := tx.Model(&data.RegionEmployee{}).
				Where(&data.RegionEmployee{
					BaseEntity: data.BaseEntity{ID: id},
				}).
				Updates(updateModels.RegionEmployee).Error
			if err != nil {
				return err
			}
		}

		if updateModels.UpdateEmployeeModels != nil && !utils.IsEmpty(updateModels.UpdateEmployeeModels) {
			err := r.employeeRepo.UpdateEmployeeWithAssociations(tx, existingRegionEmployee.EmployeeID, updateModels.UpdateEmployeeModels)
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
