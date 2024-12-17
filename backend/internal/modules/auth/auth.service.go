package auth

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/customers"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthenticationService interface {
	EmployeeLogin(email, password string) (*types.TokenPair, error)
	EmployeeRefreshAccessToken(refreshToken string) (string, error)
	CustomerRegister(input *types.CustomerRegisterDTO) (uint, error)
	CustomerLogin(email, password string) (*types.TokenPair, error)
	CustomerRefreshTokens(refreshToken string) (*types.TokenPair, error)
}

type authenticationService struct {
	repo          AuthenticationRepository
	customersRepo customers.CustomerRepository
	employeesRepo employees.EmployeeRepository
	logger        *zap.SugaredLogger
}

func NewAuthenticationService(
	repo AuthenticationRepository,
	customersRepo customers.CustomerRepository,
	employeesRepo employees.EmployeeRepository,
	logger *zap.SugaredLogger,
) AuthenticationService {
	return &authenticationService{
		repo:          repo,
		customersRepo: customersRepo,
		employeesRepo: employeesRepo,
		logger:        logger,
	}
}

func (s *authenticationService) EmployeeLogin(email, password string) (*types.TokenPair, error) {
	employee, err := s.employeesRepo.GetEmployeeByEmailOrPhone(email, "")
	if err != nil {
		wrappedErr := utils.WrapError("invalid credentials", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("this employee is not registered")
	}

	if err := utils.ComparePassword(employee.HashedPassword, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	employeeData := types.MapEmployeeToClaimsData(employee)

	accessToken, err := types.GenerateEmployeeJWT(employeeData, types.TokenAccess)
	if err != nil {
		return nil, utils.WrapError("failed to generate access token", err)
	}
	refreshToken, err := types.GenerateEmployeeJWT(employeeData, types.TokenRefresh)
	if err != nil {
		return nil, utils.WrapError("failed to generate refresh token", err)
	}

	tokenPair := types.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &tokenPair, nil
}

func (s *authenticationService) EmployeeRefreshAccessToken(refreshToken string) (string, error) {
	claims := &types.EmployeeClaims{}
	err := types.ValidateEmployeeJWT(refreshToken, claims, types.TokenRefresh)
	if err != nil {
		wrappedErr := utils.WrapError("failed to validate refresh token", err)
		return "", wrappedErr
	}

	if claims.EmployeeClaimsData.ID == 0 {
		wrappedErr := utils.WrapError("invalid refresh token payload", errors.New("id cannot be 0"))
		return "", wrappedErr
	}

	employee, err := s.employeesRepo.GetEmployeeByID(claims.EmployeeClaimsData.ID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve Customer", err)
		s.logger.Error(wrappedErr)
		return "", wrappedErr
	}

	employeeData := types.MapEmployeeToClaimsData(employee)

	newAccessToken, err := types.GenerateEmployeeJWT(employeeData, types.TokenAccess)
	if err != nil {
		wrappedErr := utils.WrapError("failed to generate access token", err)
		s.logger.Error(wrappedErr)
		return "", wrappedErr
	}

	return newAccessToken, nil
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
		Name:       input.Name,
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

func (s *authenticationService) CustomerLogin(phone, password string) (*types.TokenPair, error) {
	customer, err := s.repo.GetCustomerByPhone(phone)
	if err != nil {
		wrappedErr := utils.WrapError("error retrieving customer", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if customer == nil {
		return nil, errors.New("this Customer is not registered")
	}

	if err := utils.ComparePassword(customer.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	employeeData := types.MapCustomerToClaimsData(customer)

	accessToken, err := types.GenerateCustomerJWT(employeeData, types.TokenAccess)
	if err != nil {
		return nil, utils.WrapError("failed to generate access token", err)
	}
	refreshToken, err := types.GenerateCustomerJWT(employeeData, types.TokenRefresh)
	if err != nil {
		return nil, utils.WrapError("failed to generate refresh token", err)
	}

	tokenPair := types.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &tokenPair, nil
}

func (s *authenticationService) CustomerRefreshTokens(refreshToken string) (*types.TokenPair, error) {
	claims := &types.CustomerClaims{}
	err := types.ValidateCustomerJWT(refreshToken, claims, types.TokenRefresh)
	if err != nil {
		wrappedErr := utils.WrapError("failed to validate refresh token", err)
		return nil, wrappedErr
	}

	if claims.CustomerClaimsData.ID == 0 {
		wrappedErr := utils.WrapError("invalid refresh token payload", errors.New("id cannot be 0"))
		return nil, wrappedErr
	}

	customer, err := s.customersRepo.GetCustomerByID(claims.CustomerClaimsData.ID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve Customer", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	accessClaims := types.MapCustomerToClaimsData(customer)

	accessToken, err := types.GenerateCustomerJWT(accessClaims, types.TokenAccess)
	if err != nil {
		wrappedErr := utils.WrapError("failed to generate access token", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	customerData := types.MapCustomerToClaimsData(customer)

	refreshToken, err = types.GenerateCustomerJWT(customerData, types.TokenRefresh)
	if err != nil {
		wrappedErr := utils.WrapError("failed to generate access token", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	tokenPair := &types.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenPair, nil
}
