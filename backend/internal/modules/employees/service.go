package employees

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type EmployeeService interface {
	CreateEmployee(input types.CreateEmployeeDTO) (*types.EmployeeDTO, error)
	GetEmployeesByStore(storeID uint, roleID *uint, limit, offset int) ([]types.EmployeeDTO, error)
	GetEmployeeByID(employeeID uint) (*types.EmployeeDTO, error)
	UpdateEmployee(employeeID uint, input types.UpdateEmployeeDTO) error
	DeleteEmployee(employeeID uint) error
	GetAllRoles() ([]types.RoleDTO, error)
}

type employeeService struct {
	repo EmployeeRepository
}

func NewEmployeeService(repo EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

func (s *employeeService) CreateEmployee(input types.CreateEmployeeDTO) (*types.EmployeeDTO, error) {

	if input.Name == "" {
		return nil, errors.New("employee name cannot be empty")
	}
	if input.Email == "" {
		return nil, errors.New("employee email cannot be empty")
	}
	if input.RoleID == 0 {
		return nil, errors.New("employee role must be valid")
	}
	if input.StoreID == 0 {
		return nil, errors.New("employee must belong to a valid store")
	}
	if input.Username == "" || input.Password == "" {
		return nil, errors.New("username and password cannot be empty")
	}

	existingEmployee, err := s.repo.GetEmployeeByEmailOrPhone(input.Email, input.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking employee uniqueness: %v", err)
	}
	if existingEmployee != nil {
		return nil, errors.New("an employee with the same email or phone already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	employee := &data.Employee{
		Name:     input.Name,
		Phone:    input.Phone,
		Email:    input.Email,
		RoleID:   &input.RoleID,
		StoreID:  &input.StoreID,
		IsActive: true,
	}

	auth := &data.EmployeeAuth{
		Username:       input.Username,
		HashedPassword: string(hashedPassword),
	}

	if err := s.repo.CreateEmployee(employee, auth); err != nil {
		return nil, fmt.Errorf("failed to create employee: %v", err)
	}

	return mapToEmployeeDTO(employee), nil
}

func (s *employeeService) GetEmployeesByStore(storeID uint, roleID *uint, limit, offset int) ([]types.EmployeeDTO, error) {
	if storeID == 0 {
		return nil, errors.New("invalid store ID")
	}

	employees, err := s.repo.GetEmployeesByStore(storeID, roleID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve employees: %v", err)
	}

	dtos := make([]types.EmployeeDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *mapToEmployeeDTO(&employee)
	}
	return dtos, nil
}

func (s *employeeService) GetEmployeeByID(employeeID uint) (*types.EmployeeDTO, error) {
	if employeeID == 0 {
		return nil, errors.New("invalid employee ID")
	}

	employee, err := s.repo.GetEmployeeByID(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve employee: %v", err)
	}
	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return mapToEmployeeDTO(employee), nil
}

func (s *employeeService) UpdateEmployee(employeeID uint, input types.UpdateEmployeeDTO) error {
	if employeeID == 0 {
		return errors.New("invalid employee ID")
	}

	employee, err := s.repo.GetEmployeeByID(employeeID)
	if err != nil {
		return fmt.Errorf("failed to retrieve employee: %v", err)
	}
	if employee == nil {
		return errors.New("employee not found")
	}

	if input.Name != nil {
		employee.Name = *input.Name
	}
	if input.Phone != nil {
		employee.Phone = *input.Phone
	}
	if input.Email != nil {
		employee.Email = *input.Email
	}
	if input.RoleID != nil {
		employee.RoleID = input.RoleID
	}
	if input.StoreID != nil {
		employee.StoreID = input.StoreID
	}

	if err := s.repo.UpdateEmployee(employee); err != nil {
		return fmt.Errorf("failed to update employee: %v", err)
	}

	return nil
}

func (s *employeeService) DeleteEmployee(employeeID uint) error {
	if employeeID == 0 {
		return errors.New("invalid employee ID")
	}

	employee, err := s.repo.GetEmployeeByID(employeeID)
	if err != nil {
		return fmt.Errorf("failed to retrieve employee: %v", err)
	}
	if employee == nil {
		return errors.New("employee not found")
	}

	if err := s.repo.DeleteEmployee(employeeID); err != nil {
		return fmt.Errorf("failed to delete employee: %v", err)
	}

	return nil
}

func (s *employeeService) GetAllRoles() ([]types.RoleDTO, error) {
	roles, err := s.repo.GetAllRoles()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve roles: %v", err)
	}

	roleDTOs := make([]types.RoleDTO, len(roles))
	for i, role := range roles {
		roleDTOs[i] = types.RoleDTO{
			ID:   role.ID,
			Name: role.Name,
		}
	}
	return roleDTOs, nil
}

func mapToEmployeeDTO(employee *data.Employee) *types.EmployeeDTO {
	return &types.EmployeeDTO{
		ID:      employee.ID,
		Name:    employee.Name,
		Phone:   employee.Phone,
		Email:   employee.Email,
		Role:    employee.Role.Name,
		StoreID: *employee.StoreID,
	}
}
