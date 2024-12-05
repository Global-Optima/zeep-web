package employees

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	CreateEmployee(employee *data.Employee) error
	GetEmployeesByStore(storeID uint, role *string, limit, offset int) ([]data.Employee, error)
	GetEmployeeByID(employeeID uint) (*data.Employee, error)
	UpdateEmployee(employee *data.Employee) error
	DeleteEmployee(employeeID uint) error

	GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error)
	GetAllRoles() ([]types.EmployeeRole, error)
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

func (r *employeeRepository) GetEmployeesByStore(storeID uint, role *string, limit, offset int) ([]data.Employee, error) {
	var employees []data.Employee
	query := r.db.Where("store_id = ? AND is_active = TRUE", storeID)

	if role != nil {
		query = query.Where("role = ?", *role)
	}

	err := query.Limit(limit).Offset(offset).Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) GetEmployeeByID(employeeID uint) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.First(&employee, employeeID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &employee, err
}

func (r *employeeRepository) GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Where("email = ? OR phone = ?", email, phone).First(&employee).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &employee, err
}

func (r *employeeRepository) UpdateEmployee(employee *data.Employee) error {
	return r.db.Save(employee).Error
}

func (r *employeeRepository) DeleteEmployee(employeeID uint) error {
	return r.db.Model(&data.Employee{}).Where("id = ?", employeeID).Update("is_active", false).Error
}

func (r *employeeRepository) GetAllRoles() ([]types.EmployeeRole, error) {
	roles := []types.EmployeeRole{
		types.RoleAdmin,
		types.RoleManager,
		types.RoleBarista,
	}
	return roles, nil
}
