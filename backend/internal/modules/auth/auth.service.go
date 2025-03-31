package auth

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/employeeToken"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/customers"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthenticationService interface {
	EmployeeLogin(email, password string) (*types.Token, error)

	CustomerRegister(input *types.CustomerRegisterDTO) (uint, error)
	CustomerLogin(email, password string) (*types.Token, error)
}

type authenticationService struct {
	repo                 AuthenticationRepository
	customersRepo        customers.CustomerRepository
	employeesRepo        employees.EmployeeRepository
	employeeTokenManager employeeToken.EmployeeTokenManager
	logger               *zap.SugaredLogger
}

func NewAuthenticationService(
	repo AuthenticationRepository,
	customersRepo customers.CustomerRepository,
	employeesRepo employees.EmployeeRepository,
	employeeTokenManager employeeToken.EmployeeTokenManager,
	logger *zap.SugaredLogger,
) AuthenticationService {
	return &authenticationService{
		repo:                 repo,
		customersRepo:        customersRepo,
		employeesRepo:        employeesRepo,
		employeeTokenManager: employeeTokenManager,
		logger:               logger,
	}
}

func (s *authenticationService) EmployeeLogin(email, password string) (*types.Token, error) {
	employee, err := s.checkEmployeeCredentials(email, password)
	if err != nil {
		return nil, err
	}

	sessionToken, err := s.handleEmployeeToken(employee.ID)
	if err != nil {
		return nil, err
	}

	return &types.Token{SessionToken: sessionToken}, nil
}

func (s *authenticationService) CustomerRegister(input *types.CustomerRegisterDTO) (uint, error) {
	if err := types.ValidateCustomer(*input); err != nil {
		return 0, err
	}

	existingCustomer, err := s.repo.GetCustomerByPhone(input.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr := utils.WrapError("error checking customer uniqueness", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	if existingCustomer != nil {
		return 0, errors.New("an customer with the same phone already exists")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %v", err)
	}

	customer := &data.Customer{
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Password:   hashedPassword,
		Phone:      input.Phone,
		IsVerified: false,
		IsBanned:   false,
	}

	id, err := s.repo.CreateCustomer(customer)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create Customer", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	s.logger.Infof("customer with ID=%d CREATED successfully", id)

	return id, nil
}

func (s *authenticationService) CustomerLogin(phone, password string) (*types.Token, error) {
	customer, err := s.repo.GetCustomerByPhone(phone)
	if err != nil {
		wrappedErr := utils.WrapError("error retrieving customer", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if customer == nil {
		return nil, errors.New("this Customer is not registered")
	}

	if customer.IsBanned {
		return nil, types.ErrBannedCustomer
	}

	if err := utils.ComparePassword(customer.Password, password); err != nil {
		return nil, types.ErrInvalidCredentials
	}

	sessionToken, err := types.GenerateCustomerJWT(customer.ID)
	if err != nil {
		return nil, utils.WrapError("failed to generate access token", err)
	}

	token := types.Token{
		SessionToken: sessionToken,
	}

	return &token, nil
}

func (s *authenticationService) checkEmployeeCredentials(email, password string) (*data.Employee, error) {
	employee, err := s.employeesRepo.GetEmployeeByEmailOrPhone(email, "")
	if err != nil {
		s.logger.Error(err)
		return nil, types.ErrInvalidCredentials
	}
	if employee == nil {
		return nil, errors.New("this employee is not registered")
	}
	if !employee.IsActive {
		return nil, types.ErrInactiveEmployee
	}

	if err := utils.ComparePassword(employee.HashedPassword, password); err != nil {
		return nil, types.ErrInvalidCredentials
	}
	return employee, nil
}

func (s *authenticationService) handleEmployeeToken(employeeID uint) (string, error) {
	existingToken, err := s.employeeTokenManager.GetTokenByEmployeeID(employeeID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("failed to retrieve existing token: %w", err)
	}

	if existingToken != nil && existingToken.ExpiresAt.After(time.Now()) {
		return existingToken.Token, nil
	}

	if err := s.employeeTokenManager.DeleteTokenByEmployeeID(employeeID); existingToken != nil && err != nil {
		return "", fmt.Errorf("failed to delete expired token: %w", err)
	}

	sessionToken, err := types.GenerateEmployeeJWT(employeeID)
	if err != nil {
		return "", utils.WrapError("failed to generate session token", err)
	}

	if err := s.saveEmployeeToken(employeeID, sessionToken); err != nil {
		return "", fmt.Errorf("failed to save new token: %w", err)
	}

	return sessionToken, nil
}

func (s *authenticationService) saveEmployeeToken(employeeID uint, token string) error {
	cfg := config.GetConfig()
	now := time.Now()
	expirationTime := now.Add(cfg.JWT.EmployeeTokenTTL)

	employeeToken := &data.EmployeeToken{
		EmployeeID: employeeID,
		Token:      token,
		ExpiresAt:  expirationTime,
	}
	if err := s.employeeTokenManager.CreateToken(employeeToken); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}
