package employees

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type EmployeeService interface {
	CreateEmployee(input types.CreateEmployeeDTO) (*types.EmployeeDTO, error)
	GetEmployees(query types.GetEmployeesQuery) ([]types.EmployeeDTO, error)
	GetEmployeeByID(employeeID uint) (*types.EmployeeDTO, error)
	UpdateEmployee(employeeID uint, input types.UpdateEmployeeDTO) error
	DeleteEmployee(employeeID uint) error
	UpdatePassword(employeeID uint, input types.UpdatePasswordDTO) error
	GetAllRoles() ([]types.RoleDTO, error)
	EmployeeLogin(email, password string) (string, error)
}

type employeeService struct {
	repo EmployeeRepository
}

func NewEmployeeService(repo EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

func (s *employeeService) CreateEmployee(input types.CreateEmployeeDTO) (*types.EmployeeDTO, error) {
	if err := types.ValidateEmployee(input); err != nil {
		return nil, err
	}

	existingEmployee, err := s.repo.GetEmployeeByEmailOrPhone(input.Email, input.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking employee uniqueness: %v", err)
	}
	if existingEmployee != nil {
		return nil, errors.New("an employee with the same email or phone already exists")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	employee := &data.Employee{
		Name:           input.Name,
		Phone:          input.Phone,
		Email:          input.Email,
		Role:           input.Role,
		Type:           types.ToString(input.Type),
		HashedPassword: hashedPassword,
		IsActive:       true,
	}

	if input.Type == types.StoreEmployee && input.StoreDetails != nil {
		employee.StoreEmployee = &data.StoreEmployee{
			StoreID:     input.StoreDetails.StoreID,
			IsFranchise: input.StoreDetails.IsFranchise,
		}
	} else if input.Type == types.WarehouseEmployee && input.WarehouseDetails != nil {
		employee.WarehouseEmployee = &data.WarehouseEmployee{
			WarehouseID: input.WarehouseDetails.WarehouseID,
		}
	}

	if err := s.repo.CreateEmployee(employee); err != nil {
		return nil, fmt.Errorf("failed to create employee: %v", err)
	}

	return mapToEmployeeDTO(employee), nil
}

func (s *employeeService) GetEmployees(query types.GetEmployeesQuery) ([]types.EmployeeDTO, error) {
	employees, err := s.repo.GetEmployees(query)
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
	updateFields, err := types.PrepareUpdateFields(input)
	if err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	if err := s.repo.PartialUpdateEmployee(employeeID, updateFields); err != nil {
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

func (s *employeeService) UpdatePassword(employeeID uint, input types.UpdatePasswordDTO) error {
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

	if err := utils.ComparePassword(employee.HashedPassword, input.OldPassword); err != nil {
		return errors.New("incorrect old password")
	}

	if err := utils.IsValidPassword(input.NewPassword); err != nil {
		return fmt.Errorf("password validation failed: %v", err)
	}

	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	employee.HashedPassword = hashedPassword
	if err := s.repo.UpdateEmployee(employee); err != nil {
		return fmt.Errorf("failed to update password: %v", err)
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
			Name: string(role),
		}
	}
	return roleDTOs, nil
}

func (s *employeeService) EmployeeLogin(email, password string) (string, error) {
	employee, err := s.repo.GetEmployeeByEmailOrPhone(email, "")
	if err != nil {
		return "", fmt.Errorf("invalid credentials: %v", err)
	}

	if employee == nil {
		return "", errors.New("this employee is not registered")
	}

	if err := utils.ComparePassword(employee.HashedPassword, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	var workplaceID *uint
	var workplaceType string

	if employee.StoreEmployee != nil {
		workplaceID = &employee.StoreEmployee.StoreID
		workplaceType = "Store"
	} else if employee.WarehouseEmployee != nil {
		workplaceID = &employee.WarehouseEmployee.WarehouseID
		workplaceType = "Warehouse"
	}

	claims := utils.EmployeeClaims{
		BaseClaims: utils.BaseClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		Role:          employee.Role,
		WorkplaceID:   workplaceID,
		WorkplaceType: workplaceType,
	}

	return utils.GenerateJWT(claims, 24*time.Hour)
}

func mapToEmployeeDTO(employee *data.Employee) *types.EmployeeDTO {
	dto := &types.EmployeeDTO{
		ID:       employee.ID,
		Name:     employee.Name,
		Phone:    employee.Phone,
		Email:    employee.Email,
		Role:     employee.Role,
		IsActive: employee.IsActive,
		Type:     types.ToEmployeeType(employee.Type),
	}

	if employee.StoreEmployee != nil {
		dto.StoreDetails = &types.StoreEmployeeDetailsDTO{
			StoreID:     employee.StoreEmployee.StoreID,
			IsFranchise: employee.StoreEmployee.IsFranchise,
		}
	}

	if employee.WarehouseEmployee != nil {
		dto.WarehouseDetails = &types.WarehouseEmployeeDetailsDTO{
			WarehouseID: employee.WarehouseEmployee.WarehouseID,
		}
	}

	return dto
}
