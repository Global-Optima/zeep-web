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
	CreateEmployee(input types.CreateEmployeeDTO) (*types.EmployeeDTO, error)
	GetEmployees(query types.GetEmployeesQuery) ([]types.EmployeeDTO, error)
	GetEmployeeByID(employeeID uint) (*types.EmployeeDTO, error)
	UpdateEmployee(employeeID uint, input types.UpdateEmployeeDTO) error
	DeleteEmployee(employeeID uint) error
	UpdatePassword(employeeID uint, input types.UpdatePasswordDTO) error
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

func (s *employeeService) CreateEmployee(input types.CreateEmployeeDTO) (*types.EmployeeDTO, error) {
	if err := types.ValidateEmployee(input); err != nil {
		return nil, err
	}

	existingEmployee, err := s.repo.GetEmployeeByEmailOrPhone(input.Email, input.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr := utils.WrapError("error checking employee uniqueness", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
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
		Type:           input.Type,
		HashedPassword: hashedPassword,
		IsActive:       true,
	}

	if input.Type == data.StoreEmployeeType && input.StoreDetails != nil {
		employee.StoreEmployee = &data.StoreEmployee{
			StoreID:     input.StoreDetails.StoreID,
			IsFranchise: input.StoreDetails.IsFranchise,
		}
	} else if input.Type == data.WarehouseEmployeeType && input.WarehouseDetails != nil {
		employee.WarehouseEmployee = &data.WarehouseEmployee{
			WarehouseID: input.WarehouseDetails.WarehouseID,
		}
	}

	if err := s.repo.CreateEmployee(employee); err != nil {
		wrappedErr := utils.WrapError("failed to create employee", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return mapToEmployeeDTO(employee), nil
}

func (s *employeeService) GetEmployees(query types.GetEmployeesQuery) ([]types.EmployeeDTO, error) {
	employees, err := s.repo.GetEmployees(query)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.EmployeeDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *mapToEmployeeDTO(&employee)
	}

	return dtos, nil
}

func (s *employeeService) GetEmployeeByID(employeeID uint) (*types.EmployeeDTO, error) {
	fmt.Print("employeeID", employeeID)

	if employeeID == 0 {
		return nil, errors.New("invalid employee ID")
	}

	employee, err := s.repo.GetEmployeeByID(employeeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve employee", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
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
		wrappedErr := utils.WrapError("failed to update employee", err)
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

	if err := s.repo.DeleteEmployee(employeeID); err != nil {
		wrappedErr := utils.WrapError("failed to delete employee", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
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
		wrappedErr := utils.WrapError("failed to update password", err)
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

func mapToEmployeeDTO(employee *data.Employee) *types.EmployeeDTO {
	dto := &types.EmployeeDTO{
		ID:       employee.ID,
		Name:     employee.Name,
		Phone:    employee.Phone,
		Email:    employee.Email,
		Role:     employee.Role,
		IsActive: employee.IsActive,
		Type:     employee.Type,
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
