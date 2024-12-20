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
	CreateStoreEmployee(employee *data.Employee) error
	CreateWarehouseEmployee(employee *data.Employee) error
	//GetEmployees(filter types.GetEmployeesFilter) ([]data.Employee, error)
	GetStoreEmployees(filter *types.GetStoreEmployeesFilter) ([]data.Employee, error)
	GetWarehouseEmployees(filter *types.GetWarehouseEmployeesFilter) ([]data.Employee, error)
	GetTypedEmployeeByID(employeeID uint, employeeType data.EmployeeType) (*data.Employee, error)
	GetEmployeeByID(employeeID uint) (*data.Employee, error)
	UpdateEmployee(employeeType data.EmployeeType, employee *data.Employee) error

	PartialUpdateEmployee(employeeID uint, employeeType data.EmployeeType, employee *data.Employee) error
	DeleteEmployeeById(employeeID uint, employeeType data.EmployeeType) error
	GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error)
	GetAllRoles() ([]data.EmployeeRole, error)

	CreateEmployeeWorkday(employee *data.EmployeeWorkday) (uint, error)
	GetEmployeeWorkdayByEmployeeAndDay(employeeID uint, day data.Weekday) (*data.EmployeeWorkday, error)
	GetEmployeeWorkdayByID(workdayID uint) (*data.EmployeeWorkday, error)
	GetEmployeeWorkdaysByEmployeeID(employeeID uint) ([]data.EmployeeWorkday, error)
	UpdateEmployeeWorkdayById(workdayID uint, workday *data.EmployeeWorkday) error
	DeleteEmployeeWorkday(workdayID uint) error
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) CreateStoreEmployee(employee *data.Employee) error {
	return r.db.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Create(employee).Error; err != nil {
				return err
			}

			storeEmployee := &data.StoreEmployee{
				EmployeeID:  employee.ID,
				StoreID:     employee.StoreEmployee.StoreID,
				IsFranchise: employee.StoreEmployee.IsFranchise,
			}

			if err := tx.Create(storeEmployee).Error; err != nil {
				return err
			}

			return nil
		},
	)
}

func (r *employeeRepository) CreateWarehouseEmployee(employee *data.Employee) error {
	return r.db.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Create(employee).Error; err != nil {
				return err
			}

			warehouseEmployee := &data.WarehouseEmployee{
				EmployeeID:  employee.ID,
				WarehouseID: employee.WarehouseEmployee.WarehouseID,
			}
			if err := tx.Create(warehouseEmployee).Error; err != nil {
				return err
			}

			return nil
		},
	)
}

/*func (r *employeeRepository) GetEmployees(filter types.GetEmployeesFilter) ([]data.Employee, error) {
	var employees []data.Employee
	dbQuery := r.db.Where("is_active = TRUE")

	if filter.Type != nil {
		dbQuery = dbQuery.Where("LOWER(type) = ?", strings.ToLower(*filter.Type))
	}
	if filter.StoreID != nil {
		dbQuery = dbQuery.Joins("LEFT JOIN store_employees ON employees.id = store_employees.employee_id").
			Where("store_employees.store_id = ?", *filter.StoreID)
	}
	if filter.WarehouseID != nil {
		dbQuery = dbQuery.Joins("LEFT JOIN warehouse_employees ON employees.id = warehouse_employees.employee_id").
			Where("warehouse_employees.warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.Role != nil {
		dbQuery = dbQuery.Where("LOWER(role) = ?", strings.ToLower(*filter.Role))
	}

	var err error
	dbQuery, err = utils.ApplyPagination(dbQuery, filter.Pagination, &data.Employee{})
	if err != nil {
		return nil, err
	}

	dbQuery = dbQuery.Preload("StoreEmployee").Preload("WarehouseEmployee")

	err = dbQuery.Find(&employees).Error
	if err != nil {
		return nil, err
	}

	return employees, err
}*/

func (r *employeeRepository) GetStoreEmployees(filter *types.GetStoreEmployeesFilter) ([]data.Employee, error) {
	var employees []data.Employee
	query := r.db.Preload("StoreEmployee").Where("is_active = TRUE").
		Joins("JOIN store_employees ON employees.id = store_employees.employee_id").
		Where("store_employees.store_id = ?", filter.StoreID)

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

func (r *employeeRepository) GetWarehouseEmployees(filter *types.GetWarehouseEmployeesFilter) ([]data.Employee, error) {
	var employees []data.Employee
	query := r.db.Preload("WarehouseEmployee").Where("is_active = TRUE").
		Joins("JOIN warehouse_employees ON employees.id = warehouse_employees.employee_id").
		Where("warehouse_employees.warehouse_id = ?", filter.WarehouseID).
		Where("is_active = ?", filter.IsActive != nil && *filter.IsActive)

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.phone ILIKE ? OR employees.email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	if filter.Role != nil {
		query = query.Where("employees.role = ?", *filter.Role)
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

func (r *employeeRepository) GetTypedEmployeeByID(employeeID uint, employeeType data.EmployeeType) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Preload("StoreEmployee").
		Preload("WarehouseEmployee").
		Where("type = ?", employeeType).
		First(&employee, employeeID).Error
	return &employee, err
}

func (r *employeeRepository) GetEmployeeByID(employeeID uint) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Preload("StoreEmployee").
		Preload("WarehouseEmployee").
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
			if err := tx.Save(employee.StoreEmployee).Error; err != nil {
				return err
			}
		} else if employeeType == data.WarehouseEmployeeType && employee.WarehouseEmployee != nil {
			if err := tx.Save(employee.WarehouseEmployee).Error; err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported employee type: %v", employeeType)
		}

		return nil
	})
}

// TODO remove redundant?
func (r *employeeRepository) PartialUpdateEmployee(employeeID uint, employeeType data.EmployeeType, employee *data.Employee) error {
	if employeeType == data.StoreEmployeeType {
		return r.db.Model(&data.Employee{}).Where("id = ? AND type = ?", employeeID, employeeType).Updates(employee).Error
	}

	if employeeType == data.WarehouseEmployeeType {
		return r.db.Model(&data.Employee{}).Where("id = ? AND type = ?", employeeID, employeeType).Updates(employee).Error
	}

	return fmt.Errorf("unsupported employee type: %v", employeeType)
	//return r.db.Model(&data.Employee{}).Where("id = ?", employeeID).Updates(fields).Error
}

func (r *employeeRepository) DeleteEmployeeById(employeeID uint, employeeType data.EmployeeType) error {
	return r.db.Model(&data.Employee{}).Where("id = ? AND type = ?", employeeID, employeeType).Update("is_active", false).Error
}

func (r *employeeRepository) GetAllRoles() ([]data.EmployeeRole, error) {
	roles := []data.EmployeeRole{
		data.RoleAdmin,
		data.RoleManager,
		data.RoleBarista,
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

	existingWorkday, err := r.GetEmployeeWorkdayByEmployeeAndDay(workday.EmployeeID, workday.Day)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if existingWorkday != nil {
		return 0, fmt.Errorf("workday for employee %d in %v already exists", workday.EmployeeID, workday.Day)
	}

	err = r.db.Create(workday).Error
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
