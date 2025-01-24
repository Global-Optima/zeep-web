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
	CreateEmployee(employee *data.Employee) error
	GetStoreEmployees(storeID uint, filter *types.EmployeesFilter) ([]data.Employee, error)
	GetWarehouseEmployees(warehouseID uint, filter *types.EmployeesFilter) ([]data.Employee, error)
	GetFranchiseeEmployees(franchiseeID uint, filter *types.EmployeesFilter) ([]data.Employee, error)
	GetRegionManagers(regionID uint, filter *types.EmployeesFilter) ([]data.Employee, error)
	GetTypedEmployeeByID(employeeID uint, employeeType data.EmployeeType) (*data.Employee, error)
	GetEmployeeByID(employeeID uint) (*data.Employee, error)
	UpdateEmployee(employeeType data.EmployeeType, employee *data.Employee) error

	DeleteEmployeeById(employeeID uint, employeeType data.EmployeeType) error
	GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error)
	GetAllRoles() ([]data.EmployeeRole, error)

	CreateEmployeeWorkday(employee *data.EmployeeWorkday) (uint, error)
	GetEmployeeWorkdayByEmployeeAndDay(employeeID uint, day data.Weekday) (*data.EmployeeWorkday, error)
	GetEmployeeWorkdayByID(workdayID uint) (*data.EmployeeWorkday, error)
	GetEmployeeWorkdaysByEmployeeID(employeeID uint) ([]data.EmployeeWorkday, error)
	UpdateEmployeeWorkdayById(workdayID uint, workday *data.EmployeeWorkday) error
	DeleteEmployeeWorkday(workdayID uint) error

	GetAllStoreEmployees(storeID uint) ([]data.Employee, error)
	GetAllWarehouseEmployees(warehouseID uint) ([]data.Employee, error)
	GetAllAdminEmployees() ([]data.Employee, error)

	UpdateFranchiseeEmployee(employeeID uint, franchiseeID uint, updateModels *types.FranchiseeEmployeeUpdateModels) error
	UpdateRegionManager(employeeID uint, updateModels *types.RegionManagerEmployeeUpdateModels) error
	UpdateStoreEmployee(employeeID uint, storeID uint, updateModels *types.StoreEmployeeUpdateModels) error
	UpdateWarehouseEmployee(employeeID uint, warehouseID uint, updateModels *types.WarehouseEmployeeUpdateModels) error
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) CreateEmployee(employee *data.Employee) error {
	return r.db.Create(employee).Error
}

func (r *employeeRepository) GetStoreEmployees(storeID uint, filter *types.EmployeesFilter) ([]data.Employee, error) {
	var employees []data.Employee
	query := r.db.Model(&data.Employee{}).
		Joins("JOIN store_employees ON employees.id = store_employees.employee_id").
		Preload("StoreEmployee")

	if storeID != 0 {
		query = query.Where("store_employees.store_id = ?", storeID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Role != nil {
		query = query.Where("employees.role = ?", *filter.Role)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.phone ILIKE ? OR employees.email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Employee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve store employees: %w", err)
	}
	return employees, err
}

// TODO search with data.WarehouseEmployee
func (r *employeeRepository) GetWarehouseEmployees(warehouseID uint, filter *types.EmployeesFilter) ([]data.Employee, error) {
	var employees []data.Employee
	query := r.db.Model(&data.Employee{}).
		Joins("JOIN warehouse_employees ON warehouse_employees.employee_id = employees.id").
		Preload("WarehouseEmployee")

	if warehouseID != 0 {
		query = query.Where("warehouse_employees.warehouse_id = ?", warehouseID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.phone ILIKE ? OR employees.email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	if filter.Role != nil {
		query = query.Where("employees.role = ?", *filter.Role)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Employee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve warehouse employees: %w", err)
	}

	return employees, err
}

func (r *employeeRepository) GetFranchiseeEmployees(franchiseeID uint, filter *types.EmployeesFilter) ([]data.Employee, error) {
	var employees []data.Employee
	query := r.db.Model(&data.Employee{}).
		Joins("JOIN franchisee_employees ON employees.id = franchisee_employees.employee_id").
		Preload("FranchiseeEmployee")

	if franchiseeID != 0 {
		query = query.Where("franchisee_employees.franchisee_id = ?", franchiseeID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Role != nil {
		query = query.Where("employees.role = ?", *filter.Role)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.phone ILIKE ? OR employees.email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Employee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve franchisee employees: %w", err)
	}
	return employees, nil
}

func (r *employeeRepository) GetRegionManagers(regionID uint, filter *types.EmployeesFilter) ([]data.Employee, error) {
	var employees []data.Employee
	query := r.db.Model(&data.Employee{}).
		Joins("JOIN region_managers ON employees.id = region_managers.employee_id").
		Preload("RegionManager")

	if regionID != 0 {
		query = query.Where("region_managers.region_id = ?", regionID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Role != nil {
		query = query.Where("employees.role = ?", *filter.Role)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.phone ILIKE ? OR employees.email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Employee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve region managers: %w", err)
	}
	return employees, nil
}

func (r *employeeRepository) GetTypedEmployeeByID(employeeID uint, employeeType data.EmployeeType) (*data.Employee, error) {
	var employee data.Employee
	var err error

	query := r.db.Model(&data.Employee{}).
		Where("type = ?", employeeType)

	switch employeeType {
	case data.StoreEmployeeType:
		query = r.db.Preload("StoreEmployee")
	case data.WarehouseEmployeeType:
		query = r.db.Preload("WarehouseEmployee")
	case data.FranchiseeEmployeeType:
		query = r.db.Preload("FranchiseeEmployee")
	case data.WarehouseRegionManagerEmployeeType:
		query = r.db.Preload("RegionManager")
	case data.AdminEmployeeType:
		query = r.db.Preload("StoreEmployee")
	default:
		return nil, types.ErrUnsupportedEmployeeType
	}

	err = query.First(&employee, employeeID).Error

	return &employee, err
}

func (r *employeeRepository) GetEmployeeByID(employeeID uint) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Preload("StoreEmployee").
		Preload("WarehouseEmployee").
		Preload("RegionManager").
		Preload("FranchiseeEmployee").
		First(&employee, employeeID).Error
	return &employee, err
}

func (r *employeeRepository) GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Preload("StoreEmployee").Preload("WarehouseEmployee").
		Where("email = ? OR phone = ?", email, phone).
		First(&employee).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &employee, err
}

func (r *employeeRepository) UpdateEmployee(employeeType data.EmployeeType, employee *data.Employee) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(employee).Error; err != nil {
			return err
		}

		if employeeType == data.StoreEmployeeType && employee.StoreEmployee != nil {
			if err := tx.Updates(employee.StoreEmployee).Error; err != nil {
				return err
			}
		} else if employeeType == data.WarehouseEmployeeType && employee.WarehouseEmployee != nil {
			if err := tx.Updates(employee.WarehouseEmployee).Error; err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported employee type: %v", employeeType)
		}

		return nil
	})
}

func (r *employeeRepository) UpdateFranchiseeEmployee(employeeID uint, franchiseeID uint, updateModels *types.FranchiseeEmployeeUpdateModels) error {
	if updateModels == nil {
		return types.ErrNothingToUpdate
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels.Employee != nil {
			if !data.IsAllowableRole(data.FranchiseeEmployeeType, updateModels.Employee.Role) {
				return types.ErrEmployeeTypeAndRoleMismatch
			}

			if err := tx.Model(&data.Employee{}).
				Where("id = ?", employeeID).
				Updates(updateModels.Employee).Error; err != nil {
				return err
			}
		}

		if updateModels.FranchiseeEmployee != nil {
			if err := tx.Model(&data.FranchiseeEmployee{}).
				Where("employee_id = ? AND franchisee_id = ?", employeeID, franchiseeID).
				Updates(updateModels.FranchiseeEmployee).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *employeeRepository) UpdateRegionManager(employeeID uint, updateModels *types.RegionManagerEmployeeUpdateModels) error {
	if updateModels == nil {
		return types.ErrNothingToUpdate
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels.Employee != nil {
			if !data.IsAllowableRole(data.WarehouseRegionManagerEmployeeType, updateModels.Employee.Role) {
				return types.ErrEmployeeTypeAndRoleMismatch
			}

			if err := tx.Model(&data.Employee{}).
				Where("id = ?", employeeID).
				Updates(updateModels.Employee).Error; err != nil {
				return err
			}
		}

		if updateModels.RegionManager != nil {
			if err := tx.Model(&data.RegionManager{}).
				Where("employee_id = ?", employeeID).
				Updates(updateModels.RegionManager).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *employeeRepository) UpdateStoreEmployee(employeeID uint, storeID uint, updateModels *types.StoreEmployeeUpdateModels) error {
	if updateModels == nil {
		return types.ErrNothingToUpdate
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels.Employee != nil {
			if !data.IsAllowableRole(data.StoreEmployeeType, updateModels.Employee.Role) {
				return types.ErrEmployeeTypeAndRoleMismatch
			}

			if err := tx.Model(&data.Employee{}).
				Where("id = ?", employeeID).
				Updates(updateModels.Employee).Error; err != nil {
				return err
			}
		}

		if updateModels.StoreEmployee != nil {
			if err := tx.Model(&data.StoreEmployee{}).
				Where("employee_id = ? AND store_id = ?", employeeID, storeID).
				Updates(updateModels.StoreEmployee).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *employeeRepository) UpdateWarehouseEmployee(employeeID uint, warehouseID uint, updateModels *types.WarehouseEmployeeUpdateModels) error {
	if updateModels == nil {
		return types.ErrNothingToUpdate
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels.Employee != nil {
			if !data.IsAllowableRole(data.WarehouseEmployeeType, updateModels.Employee.Role) {
				return types.ErrEmployeeTypeAndRoleMismatch
			}

			if err := tx.Model(&data.Employee{}).
				Where("id = ?", employeeID).
				Updates(updateModels.Employee).Error; err != nil {
				return err
			}
		}

		if updateModels.WarehouseEmployee != nil {
			if err := tx.Model(&data.WarehouseEmployee{}).
				Where("employee_id = ? AND warehouse_id = ?", employeeID, warehouseID).
				Updates(updateModels.WarehouseEmployee).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *employeeRepository) DeleteEmployeeById(employeeID uint, employeeType data.EmployeeType) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND type = ?", employeeID, employeeType).Delete(&data.Employee{}).Error; err != nil {
			return err
		}

		switch employeeType {
		case data.StoreEmployeeType:
			if err := tx.Where("employee_id = ?", employeeID).Delete(&data.StoreEmployee{}).Error; err != nil {
				return err
			}
		case data.WarehouseEmployeeType:
			if err := tx.Where("employee_id = ?", employeeID).Delete(&data.WarehouseEmployee{}).Error; err != nil {
				return err
			}
		case data.FranchiseeEmployeeType:
			if err := tx.Where("employee_id = ?", employeeID).Delete(&data.FranchiseeEmployee{}).Error; err != nil {
				return err
			}
		case data.WarehouseRegionManagerEmployeeType:
			if err := tx.Where("employee_id = ?", employeeID).Delete(&data.RegionManager{}).Error; err != nil {
				return err
			}
		case data.AdminEmployeeType:

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

func (r *employeeRepository) GetAllRoles() ([]data.EmployeeRole, error) {
	roles := []data.EmployeeRole{
		data.RoleAdmin,
		data.RoleOwner,
		data.RoleFranchiseOwner,
		data.RoleFranchiseManager,
		data.RoleStoreManager,
		data.RoleWarehouseRegionManager,
		data.RoleWarehouseManager,
		data.RoleBarista,
		data.RoleWarehouseEmployee,
	}
	return roles, nil
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
	adminRoles := []data.EmployeeRole{data.RoleAdmin}
	err := r.db.Model(&data.Employee{}).
		Where("role IN (?)", adminRoles).
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}
