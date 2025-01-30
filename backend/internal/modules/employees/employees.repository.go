package employees

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	CreateEmployee(employee *data.Employee) (uint, error)

	GetEmployeeByID(employeeID uint) (*data.Employee, error)
	UpdateEmployee(id uint, employee *data.Employee) error
	ReassignEmployeeType(employeeID uint, newType data.EmployeeType, workplaceID uint) error

	DeleteTypedEmployeeById(employeeID, workplaceID uint, employeeType data.EmployeeType) error
	GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error)

	CreateEmployeeWorkday(employee *data.EmployeeWorkday) (uint, error)
	GetEmployeeWorkdayByEmployeeAndDay(employeeID uint, day data.Weekday) (*data.EmployeeWorkday, error)
	GetEmployeeWorkdayByID(workdayID uint) (*data.EmployeeWorkday, error)
	GetEmployeeWorkdaysByEmployeeID(employeeID uint) ([]data.EmployeeWorkday, error)
	UpdateEmployeeWorkdayById(workdayID uint, workday *data.EmployeeWorkday) error
	DeleteEmployeeWorkday(workdayID uint) error

	GetAllStoreEmployees(storeID uint) ([]data.Employee, error)
	GetAllRegionEmployees(regionID uint) ([]data.Employee, error)
	GetAllFranchiseeEmployees(franchiseeID uint) ([]data.Employee, error)
	GetAllWarehouseEmployees(warehouseID uint) ([]data.Employee, error)
	GetAllAdminEmployees() ([]data.Employee, error)
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) CreateEmployee(employee *data.Employee) (uint, error) {
	err := r.db.Create(employee).Error
	if err != nil {
		return 0, err
	}
	return employee.ID, nil
}

func (r *employeeRepository) GetEmployeesByRoles(roles []data.EmployeeRole) ([]data.Employee, error) {
	var employees []data.Employee
	err := r.db.Where("role IN ?", roles).
		Joins("JOIN store_employees ON employees.id = store_employees.employee_id").
		Where("store_employees.role IN (?)", roles).
		Joins("JOIN warehouse_employees ON employees.id = warehouse_employees.employee_id").
		Where("warehouse_employees.role IN (?)", roles).
		Joins("JOIN franchisee_employees ON employees.id = franchisee_employees.employee_id").
		Where("warehouse_employees.role IN (?)", roles).
		Joins("JOIN region_employees ON employees.id = region_employees.employee_id").
		Where("region_employees.role IN (?)", roles).
		Joins("JOIN admin_employees ON employees.id = admin_employees.employee_id").
		Where("admin_employees.role IN (?)", roles).
		Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, err
}

func (r *employeeRepository) GetStoreEmployees(storeID uint, filter *types.EmployeesFilter) ([]data.StoreEmployee, error) {
	var employees []data.StoreEmployee
	query := r.db.Model(&data.StoreEmployee{}).
		Where("store_id = ?", storeID).Preload("Employee")

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

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreEmployee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve store employees: %w", err)
	}
	return employees, nil
}

func (r *employeeRepository) GetWarehouseEmployees(warehouseID uint, filter *types.EmployeesFilter) ([]data.WarehouseEmployee, error) {
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

func (r *employeeRepository) GetFranchiseeEmployees(franchiseeID uint, filter *types.EmployeesFilter) ([]data.FranchiseeEmployee, error) {
	var employees []data.FranchiseeEmployee
	query := r.db.Model(&data.FranchiseeEmployee{}).
		Where("franchisee_id = ?", franchiseeID).Preload("Employee")

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

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.FranchiseeEmployee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve franchisee employees: %w", err)
	}
	return employees, nil
}

func (r *employeeRepository) GetRegionEmployees(regionID uint, filter *types.EmployeesFilter) ([]data.RegionEmployee, error) {
	var employees []data.RegionEmployee
	query := r.db.Model(&data.RegionEmployee{}).
		Where("region_id = ?", regionID).Preload("Employee")

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

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.RegionEmployee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve region managers: %w", err)
	}
	return employees, nil
}

func (r *employeeRepository) GetAdminEmployees(filter *types.EmployeesFilter) ([]data.AdminEmployee, error) {
	var employees []data.AdminEmployee
	query := r.db.Model(&data.AdminEmployee{}).Preload("Employee")

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

func (r *employeeRepository) GetStoreEmployeeByID(id, storeID uint) (*data.StoreEmployee, error) {
	var storeEmployee data.StoreEmployee
	err := r.db.Model(&data.StoreEmployee{}).
		Preload("Employee").
		Where("id = ? AND store_id = ?", id, storeID).
		First(&storeEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve store employee by ID: %w", err)
	}
	return &storeEmployee, nil
}

func (r *employeeRepository) GetWarehouseEmployeeByID(id, warehouseID uint) (*data.WarehouseEmployee, error) {
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

func (r *employeeRepository) GetFranchiseeEmployeeByID(id, franchiseeID uint) (*data.FranchiseeEmployee, error) {
	var franchiseeEmployee data.FranchiseeEmployee
	err := r.db.Model(&data.FranchiseeEmployee{}).
		Preload("Employee").
		Where("id = ? AND franchisee_id = ?", id, franchiseeID).
		First(&franchiseeEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve franchisee employee by ID: %w", err)
	}
	return &franchiseeEmployee, nil
}

func (r *employeeRepository) GetRegionEmployeeByID(id, regionID uint) (*data.RegionEmployee, error) {
	var regionEmployee data.RegionEmployee
	err := r.db.Model(&data.RegionEmployee{}).
		Preload("Employee").
		Where("id = ? AND region_id = ?", id, regionID).
		First(&regionEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve region manager by ID: %w", err)
	}
	return &regionEmployee, nil
}

func (r *employeeRepository) GetAdminEmployeeByID(id uint) (*data.AdminEmployee, error) {
	var adminEmployee data.AdminEmployee
	err := r.db.Model(&data.AdminEmployee{}).
		Preload("Employee").
		Where("id = ?", id).
		First(&adminEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin employee by ID: %w", err)
	}
	return &adminEmployee, nil
}

func (r *employeeRepository) GetEmployeeByID(employeeID uint) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Model(&data.Employee{}).
		Preload("StoreEmployee").
		Preload("WarehouseEmployee").
		Preload("RegionEmployee").
		Preload("FranchiseeEmployee").
		Preload("AdminEmployee").
		First(&employee, employeeID).Error
	return &employee, err
}

func (r *employeeRepository) GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Model(&data.Employee{}).
		Preload("StoreEmployee").
		Preload("WarehouseEmployee").
		Preload("FranchiseeEmployee").
		Preload("RegionEmployee").
		Preload("AdminEmployee").
		Where("email = ? OR phone = ?", email, phone).
		First(&employee).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &employee, err
}

func (r *employeeRepository) UpdateEmployee(id uint, employee *data.Employee) error {
	err := r.db.Model(&data.Employee{}).Where("id = ?", id).Updates(employee).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *employeeRepository) ReassignEmployeeType(employeeID uint, newType data.EmployeeType, workplaceID uint) error {
	employee, err := r.GetEmployeeByID(employeeID)
	if err != nil {
		return err
	}

	typeMappings := map[data.EmployeeType]struct {
		deleteModel interface{}
		createModel interface{}
	}{
		data.StoreEmployeeType: {
			deleteModel: &data.StoreEmployee{},
			createModel: &data.StoreEmployee{EmployeeID: employeeID, StoreID: workplaceID},
		},
		data.WarehouseEmployeeType: {
			deleteModel: &data.WarehouseEmployee{},
			createModel: &data.WarehouseEmployee{EmployeeID: employeeID, WarehouseID: workplaceID},
		},
		data.RegionEmployeeType: {
			deleteModel: &data.RegionEmployee{},
			createModel: &data.RegionEmployee{EmployeeID: employeeID, RegionID: workplaceID},
		},
		data.FranchiseeEmployeeType: {
			deleteModel: &data.FranchiseeEmployee{},
			createModel: &data.FranchiseeEmployee{EmployeeID: employeeID, FranchiseeID: workplaceID},
		},
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if deleteModel, ok := typeMappings[employee.GetType()]; ok {
			if err := tx.Where("employee_id = ?", employeeID).Delete(deleteModel.deleteModel).Error; err != nil {
				return err
			}
		} else {
			return types.ErrUnsupportedEmployeeType
		}

		if err := tx.Model(&data.Employee{}).Where("id = ?", employeeID).Update("type", newType).Error; err != nil {
			return err
		}

		if newTypeMapping, ok := typeMappings[newType]; ok {
			if err := tx.Create(newTypeMapping.createModel).Error; err != nil {
				return err
			}
		} else {
			return types.ErrUnsupportedEmployeeType
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *employeeRepository) UpdateStoreEmployee(id uint, storeID uint, storeEmployee *data.StoreEmployee) error {
	if storeEmployee == nil {
		return types.ErrNothingToUpdate
	}
	return r.db.Model(&data.StoreEmployee{}).
		Where("id = ? AND store_id = ?", id, storeID).
		Updates(storeEmployee).Error
}

func (r *employeeRepository) UpdateWarehouseEmployee(id uint, warehouseID uint, warehouseEmployee *data.WarehouseEmployee) error {
	if warehouseEmployee == nil {
		return types.ErrNothingToUpdate
	}
	return r.db.Model(&data.WarehouseEmployee{}).
		Where("id = ? AND warehouse_id = ?", id, warehouseID).
		Updates(warehouseEmployee).Error
}

func (r *employeeRepository) UpdateFranchiseeEmployee(id uint, franchiseeID uint, franchiseeEmployee *data.FranchiseeEmployee) error {
	if franchiseeEmployee == nil {
		return types.ErrNothingToUpdate
	}
	return r.db.Model(&data.FranchiseeEmployee{}).
		Where("id = ? AND franchisee_id = ?", id, franchiseeID).
		Updates(franchiseeEmployee).Error
}

func (r *employeeRepository) UpdateRegionEmployee(id uint, regionID uint, regionEmployee *data.RegionEmployee) error {
	if regionEmployee == nil {
		return types.ErrNothingToUpdate
	}
	return r.db.Model(&data.RegionEmployee{}).
		Where("id = ? AND region_id = ?", id, regionID).
		Updates(regionEmployee).Error
}

func (r *employeeRepository) UpdateAdminEmployee(id uint, adminEmployee *data.AdminEmployee) error {
	if adminEmployee == nil {
		return types.ErrNothingToUpdate
	}
	return r.db.Model(&data.AdminEmployee{}).
		Where("id = ?", id).
		Updates(adminEmployee).Error
}

func (r *employeeRepository) DeleteTypedEmployeeById(employeeID, workplaceID uint, employeeType data.EmployeeType) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND type = ?", employeeID, employeeType).Delete(&data.Employee{}).Error; err != nil {
			return err
		}

		switch employeeType {
		case data.StoreEmployeeType:
			if err := tx.Where("id = ? AND store_id", employeeID, workplaceID).Delete(&data.StoreEmployee{}).Error; err != nil {
				return err
			}
		case data.WarehouseEmployeeType:
			if err := tx.Where("id = ? AND warehouse_id", employeeID, workplaceID).Delete(&data.WarehouseEmployee{}).Error; err != nil {
				return err
			}
		case data.FranchiseeEmployeeType:
			if err := tx.Where("id = ? AND franchisee_id", employeeID, workplaceID).Delete(&data.FranchiseeEmployee{}).Error; err != nil {
				return err
			}
		case data.RegionEmployeeType:
			if err := tx.Where("id = ? AND region_id", employeeID, workplaceID).Delete(&data.RegionEmployee{}).Error; err != nil {
				return err
			}
		default:
			return types.ErrUnsupportedEmployeeType
		}

		if err := tx.Where("employee_id = ?", employeeID).Delete(&data.EmployeeWorkday{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *employeeRepository) GetEmployeeWorkdayByEmployeeAndDay(employeeID uint, day data.Weekday) (*data.EmployeeWorkday, error) {
	var workday data.EmployeeWorkday

	if employeeID == 0 {
		return nil, fmt.Errorf("employeeID cannot be zero")
	}

	err := r.db.Where("employee_id = ? AND day = ?", employeeID, day).First(&workday).Error
	if err != nil {
		return nil, err
	}
	return &workday, err
}

func (r *employeeRepository) CreateEmployeeWorkday(workday *data.EmployeeWorkday) (uint, error) {
	if workday == nil {
		return 0, fmt.Errorf("workday is nil")
	}

	err := r.db.Create(workday).Error
	if err != nil {
		return 0, err
	}
	return workday.ID, nil
}

func (r *employeeRepository) GetEmployeeWorkdayByID(workdayID uint) (*data.EmployeeWorkday, error) {
	var workday data.EmployeeWorkday
	err := r.db.
		Preload("Employee").
		First(&workday, workdayID).Error

	if err != nil {
		return nil, err
	}
	return &workday, nil
}

func (r *employeeRepository) GetEmployeeWorkdaysByEmployeeID(employeeID uint) ([]data.EmployeeWorkday, error) {
	var workdays []data.EmployeeWorkday
	err := r.db.Where("employee_id = ?", employeeID).Find(&workdays).Error
	if err != nil {
		return nil, err
	}
	return workdays, nil
}

func (r *employeeRepository) UpdateEmployeeWorkdayById(workdayID uint, workday *data.EmployeeWorkday) error {
	return r.db.Model(&data.EmployeeWorkday{}).Where("id = ?", workdayID).Updates(workday).Error
}

func (r *employeeRepository) DeleteEmployeeWorkday(workdayID uint) error {
	return r.db.Where("id = ?", workdayID).Delete(&data.EmployeeWorkday{}).Error
}

func (r *employeeRepository) GetAllStoreEmployees(storeID uint) ([]data.Employee, error) {
	var employees []data.Employee

	err := r.db.Model(&data.Employee{}).
		Joins("INNER JOIN store_employees ON store_employees.employee_id = employees.id").
		Where("store_employees.store_id = ?", storeID).
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *employeeRepository) GetAllWarehouseEmployees(warehouseID uint) ([]data.Employee, error) {
	var employees []data.Employee

	err := r.db.Model(&data.Employee{}).
		Joins("INNER JOIN warehouse_employees ON warehouse_employees.employee_id = employees.id").
		Where("warehouse_employees.warehouse_id = ?", warehouseID).
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *employeeRepository) GetAllAdminEmployees() ([]data.Employee, error) {
	var employees []data.Employee

	err := r.db.Model(&data.Employee{}).
		Joins("INNER JOIN warehouse_employees ON warehouse_employees.employee_id = employees.id").
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *employeeRepository) GetAllRegionEmployees(regionID uint) ([]data.Employee, error) {
	var employees []data.Employee

	err := r.db.Model(&data.Employee{}).
		Joins("INNER JOIN region_employees ON region_employees.employee_id = employees.id").
		Where("region_employees.region_id = ?", regionID).
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *employeeRepository) GetAllFranchiseeEmployees(franchiseeID uint) ([]data.Employee, error) {
	var employees []data.Employee

	err := r.db.Model(&data.Employee{}).
		Joins("INNER JOIN franchisee_employees ON franchisee_employees.employee_id = employees.id").
		Where("franchisee_employees.franchisee_id = ?", franchiseeID).
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}
