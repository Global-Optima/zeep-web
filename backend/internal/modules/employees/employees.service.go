package employees

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EmployeeService interface {
	CreateStoreEmployee(input *types.CreateStoreEmployeeDTO) (*types.StoreEmployeeDTO, error)
	CreateWarehouseEmployee(input *types.CreateWarehouseEmployeeDTO) (*types.WarehouseEmployeeDTO, error)
	GetStoreEmployees(filter *types.GetStoreEmployeesFilter) ([]types.StoreEmployeeDTO, error)
	GetWarehouseEmployees(filter *types.GetWarehouseEmployeesFilter) ([]types.WarehouseEmployeeDTO, error)
	GetStoreEmployeeByID(employeeID uint) (*types.StoreEmployeeDTO, error)
	GetWarehouseEmployeeByID(employeeID uint) (*types.WarehouseEmployeeDTO, error)
	UpdateStoreEmployee(employeeID uint, input *types.UpdateStoreEmployeeDTO) error
	UpdateWarehouseEmployee(employeeID uint, input *types.UpdateWarehouseEmployeeDTO) error

	DeleteEmployee(employeeID uint) error
	UpdatePassword(employeeID uint, input *types.UpdatePasswordDTO) error
	GetAllRoles() ([]types.RoleDTO, error)
}

type employeeService struct {
	repo   EmployeeRepository
	logger *zap.SugaredLogger
}

func NewEmployeeService(repo EmployeeRepository, logger *zap.SugaredLogger) EmployeeService {
	return &employeeService{
		repo:   repo,
		logger: logger,
	}
}

func (s *employeeService) CreateStoreEmployee(input *types.CreateStoreEmployeeDTO) (*types.StoreEmployeeDTO, error) {
	if err := types.ValidateStoreEmployee(input); err != nil {
		return nil, err
	}

	hashedPassword, err := s.createNewEmployeePassword(&input.CreateEmployeeDTO)
	if err != nil {
		wrappedErr := utils.WrapError("error creating employee", err)
		return nil, wrappedErr
	}

	employee := &data.Employee{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Phone:          input.Phone,
		Email:          input.Email,
		Role:           input.Role,
		HashedPassword: hashedPassword,
		IsActive:       true,
		StoreEmployee: &data.StoreEmployee{
			StoreID:     input.StoreID,
			IsFranchise: input.IsFranchise,
		},
	}

	if err := s.repo.CreateStoreEmployee(employee); err != nil {
		wrappedErr := utils.WrapError("failed to create store employee", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return types.MapToStoreEmployeeDTO(employee), nil
}

func (s *employeeService) CreateWarehouseEmployee(input *types.CreateWarehouseEmployeeDTO) (*types.WarehouseEmployeeDTO, error) {
	if err := types.ValidateWarehouseEmployee(input); err != nil {
		return nil, err
	}

	hashedPassword, err := s.createNewEmployeePassword(&input.CreateEmployeeDTO)
	if err != nil {
		wrappedErr := utils.WrapError("error creating employee", err)
		return nil, wrappedErr
	}

	employee := &data.Employee{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Phone:          input.Phone,
		Email:          input.Email,
		Role:           input.Role,
		HashedPassword: hashedPassword,
		IsActive:       true,
		WarehouseEmployee: &data.WarehouseEmployee{
			WarehouseID: input.WarehouseID,
		},
	}

	if err := s.repo.CreateWarehouseEmployee(employee); err != nil {
		wrappedErr := utils.WrapError("failed to create warehouse employee", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return types.MapToWarehouseEmployeeDTO(employee), nil
}

func (s *employeeService) createNewEmployeePassword(input *types.CreateEmployeeDTO) (string, error) {
	existingEmployee, err := s.repo.GetEmployeeByEmailOrPhone(input.Email, input.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr := utils.WrapError("error checking employee uniqueness", err)
		return "", wrappedErr
	}
	if existingEmployee != nil {
		return "", errors.New("an employee with the same email or phone already exists")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return hashedPassword, nil
}

func (s *employeeService) GetStoreEmployees(filter *types.GetStoreEmployeesFilter) ([]types.StoreEmployeeDTO, error) {
	employees, err := s.repo.GetStoreEmployees(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve store employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.StoreEmployeeDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *types.MapToStoreEmployeeDTO(&employee)
	}

	return dtos, nil
}

func (s *employeeService) GetWarehouseEmployees(filter *types.GetWarehouseEmployeesFilter) ([]types.WarehouseEmployeeDTO, error) {
	employees, err := s.repo.GetWarehouseEmployees(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve store employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.WarehouseEmployeeDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *types.MapToWarehouseEmployeeDTO(&employee)
	}

	return dtos, nil
}

func (s *employeeService) GetStoreEmployeeByID(employeeID uint) (*types.StoreEmployeeDTO, error) {
	if employeeID == 0 {
		return nil, errors.New("invalid employee ID")
	}

	employee, err := s.repo.GetTypedEmployeeByID(employeeID, data.StoreEmployeeType)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve store employee for employee with ID = %d: %w", employeeID, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToStoreEmployeeDTO(employee), nil
}

func (s *employeeService) GetWarehouseEmployeeByID(employeeID uint) (*types.WarehouseEmployeeDTO, error) {
	if employeeID == 0 {
		return nil, errors.New("invalid employee ID")
	}

	employee, err := s.repo.GetTypedEmployeeByID(employeeID, data.WarehouseEmployeeType)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve warehouse employee for employee with ID = %d: %w", employeeID, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToWarehouseEmployeeDTO(employee), nil
}

func (s *employeeService) UpdateStoreEmployee(employeeID uint, input *types.UpdateStoreEmployeeDTO) error {
	employee, err := types.StoreEmployeeUpdateFields(input)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.repo.PartialUpdateEmployee(employeeID, data.StoreEmployeeType, employee); err != nil {
		wrappedErr := fmt.Errorf("failed to update store employee for employee with ID = %d: %w", employeeID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *employeeService) UpdateWarehouseEmployee(employeeID uint, input *types.UpdateWarehouseEmployeeDTO) error {
	employee, err := types.WarehouseEmployeeUpdateFields(input)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.repo.PartialUpdateEmployee(employeeID, data.WarehouseEmployeeType, employee); err != nil {
		wrappedErr := fmt.Errorf("failed to update store employee for employee with ID = %d: %w", employeeID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *employeeService) DeleteEmployee(employeeID uint) error {
	if employeeID == 0 {
		return errors.New("invalid employee ID")
	}

	employee, err := s.repo.GetEmployeeByID(employeeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve employee", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	if employee == nil {
		return errors.New("employee not found")
	}

	if err := s.repo.DeleteEmployeeById(employeeID, employee.Type); err != nil {
		wrappedErr := fmt.Errorf("failed to delete employee with ID = %d: %w", employeeID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *employeeService) UpdatePassword(employeeID uint, input *types.UpdatePasswordDTO) error {
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
	if err := s.repo.UpdateEmployee(employee.Type, employee); err != nil {
		wrappedErr := fmt.Errorf("failed to update password for employee with ID = %d: %w", employeeID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
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
