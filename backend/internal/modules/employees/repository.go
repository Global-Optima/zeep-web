package employees

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	CreateEmployee(employee *data.Employee, auth *data.EmployeeAuth) error
	GetEmployeesByStore(storeID uint, roleID *uint, limit, offset int) ([]data.Employee, error)
	GetEmployeeByID(employeeID uint) (*data.Employee, error)
	UpdateEmployee(employee *data.Employee) error
	DeleteEmployee(employeeID uint) error

	CreateEmployeeAuth(auth *data.EmployeeAuth) error
	UpdateEmployeeAuth(auth *data.EmployeeAuth) error
	GetEmployeeAuthByUsername(username string) (*data.EmployeeAuth, error)
	GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error)
	GetAllRoles() ([]data.EmployeeRole, error)
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) CreateEmployee(employee *data.Employee, auth *data.EmployeeAuth) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(employee).Error; err != nil {
			return err
		}

		auth.EmployeeID = employee.ID
		if err := tx.Create(auth).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *employeeRepository) GetEmployeesByStore(storeID uint, roleID *uint, limit, offset int) ([]data.Employee, error) {
	var employees []data.Employee
	query := r.db.Where("store_id = ? AND is_active = TRUE", storeID)
	if roleID != nil {
		query = query.Where("role_id = ?", *roleID)
	}
	err := query.Limit(limit).Offset(offset).Preload("Role").Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) GetEmployeeByID(employeeID uint) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Preload("Role").First(&employee, employeeID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &employee, err
}

func (r *employeeRepository) GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Where("email = ? OR phone = ?", email, phone).Preload("Role").First(&employee).Error
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

func (r *employeeRepository) CreateEmployeeAuth(auth *data.EmployeeAuth) error {
	return r.db.Create(auth).Error
}

func (r *employeeRepository) UpdateEmployeeAuth(auth *data.EmployeeAuth) error {
	return r.db.Save(auth).Error
}

func (r *employeeRepository) GetEmployeeAuthByUsername(username string) (*data.EmployeeAuth, error) {
	var auth data.EmployeeAuth
	err := r.db.Where("username = ?", username).First(&auth).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &auth, err
}

func (r *employeeRepository) GetAllRoles() ([]data.EmployeeRole, error) {
	var roles []data.EmployeeRole
	err := r.db.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
