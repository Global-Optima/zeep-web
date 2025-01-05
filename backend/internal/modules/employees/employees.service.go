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

	CreateEmployeeWorkDay(dto *types.CreateEmployeeWorkdayDTO) (uint, error)
	GetEmployeeWorkday(workdayID uint) (*types.EmployeeWorkdayDTO, error)
	GetEmployeeWorkdays(employeeID uint) ([]types.EmployeeWorkdayDTO, error)
	UpdateEmployeeWorkday(workdayID uint, dto *types.UpdateEmployeeWorkdayDTO) error
	DeleteEmployeeWorkday(workdayID uint) error
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
	employee, err := types.CreateToStoreEmployee(input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create store employee", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	existingEmployee, err := s.repo.GetEmployeeByEmailOrPhone(employee.Email, employee.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr := utils.WrapError("error checking employee uniqueness", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	if existingEmployee != nil {
		return nil, types.ErrEmployeeAlreadyExists
	}

	if err := s.repo.CreateEmployee(employee); err != nil {
		wrappedErr := utils.WrapError("failed to create store employee", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return types.MapToStoreEmployeeDTO(employee), nil
}

func (s *employeeService) CreateWarehouseEmployee(input *types.CreateWarehouseEmployeeDTO) (*types.WarehouseEmployeeDTO, error) {
	employee, err := types.CreateToWarehouseEmployee(input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create warehouse employee", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	existingEmployee, err := s.repo.GetEmployeeByEmailOrPhone(employee.Email, employee.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr := utils.WrapError("error checking employee uniqueness", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	if existingEmployee != nil {
		return nil, types.ErrEmployeeAlreadyExists
	}

	if err := s.repo.CreateEmployee(employee); err != nil {
		wrappedErr := utils.WrapError("failed to create warehouse employee", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return types.MapToWarehouseEmployeeDTO(employee), nil
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

func (s *employeeService) CreateEmployeeWorkDay(dto *types.CreateEmployeeWorkdayDTO) (uint, error) {
	errMsg := "failed to create employee work day"

	workday, err := types.ValidateEmployeeWorkday(dto)
	if err != nil {
		wrappedErr := utils.WrapError(errMsg, err)
		s.logger.Error(wrappedErr)
		return 0, err
	}

	existingWorkday, err := s.repo.GetEmployeeWorkdayByEmployeeAndDay(workday.EmployeeID, workday.Day)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr := utils.WrapError(errMsg, err)
		s.logger.Error(wrappedErr)
		return 0, err
	}
	if existingWorkday != nil {
		err = fmt.Errorf("%w: not unique workday for employeeID %d in %v ", types.ErrWorkdayAlreadyExists, dto.EmployeeID, dto.Day)
		wrappedErr := utils.WrapError(errMsg, err)
		s.logger.Error(wrappedErr)
		return 0, err
	}

	id, err := s.repo.CreateEmployeeWorkday(workday)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create employee workday", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	return id, nil
}

func (s *employeeService) GetEmployeeWorkday(workdayID uint) (*types.EmployeeWorkdayDTO, error) {
	workday, err := s.repo.GetEmployeeWorkdayByID(workdayID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve employee workday", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dto := types.MapToEmployeeWorkdayDTO(workday)
	return dto, nil
}

func (s *employeeService) GetEmployeeWorkdays(employeeID uint) ([]types.EmployeeWorkdayDTO, error) {
	workdays, err := s.repo.GetEmployeeWorkdaysByEmployeeID(employeeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve employee workdays", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	dtos := make([]types.EmployeeWorkdayDTO, len(workdays))
	for i, workday := range workdays {
		dtos[i] = *types.MapToEmployeeWorkdayDTO(&workday)
	}

	return dtos, nil
}

func (s *employeeService) UpdateEmployeeWorkday(workdayID uint, dto *types.UpdateEmployeeWorkdayDTO) error {
	workday, err := types.WorkdaysUpdateFields(dto)
	if err != nil {
		wrappedErr := utils.WrapError("failed to update employee workday", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	err = s.repo.UpdateEmployeeWorkdayById(workdayID, workday)
	if err != nil {
		wrappedErr := utils.WrapError("failed to update employee workday", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *employeeService) DeleteEmployeeWorkday(workdayID uint) error {
	err := s.repo.DeleteEmployeeWorkday(workdayID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to delete employee workday", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}
