package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
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
	EmployeeFaceRecognitionPass(c *gin.Context, image *multipart.FileHeader) (*types.Token, error)

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

func (s *authenticationService) EmployeeFaceRecognitionPass(c *gin.Context, image *multipart.FileHeader) (*types.Token, error) {
	claims, _, err := types.ExtractEmployeeSessionTokenAndValidate(c)
	if err != nil {
		return nil, err
	}

	embedding, err := s.repo.GetEmployeeEmbedding(claims.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get embedding: %w", err)
	}
	if embedding == "" {
		return nil, errors.New("no face embedding found for this user")
	}

	fileBytes, err := readFileHeader(image)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 4) Отправить запрос POST /compare к face-сервису
	match, _, err := callFaceServiceCompare(fileBytes, embedding)
	if err != nil {
		return nil, fmt.Errorf("face compare request failed: %w", err)
	}

	// 5) Если лицо совпало — выдаём новый токен
	if match {
		tokenStr, err := types.GenerateEmployeeJWT(claims.EmployeeID, true)
		if err != nil {
			return nil, err
		}
		return &types.Token{SessionToken: tokenStr}, nil
	}

	return nil, fmt.Errorf("MFA failed")
}

// Вспомогательный метод чтения image
func readFileHeader(fileHeader *multipart.FileHeader) ([]byte, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

// Вызов face_service по /compare
func callFaceServiceCompare(fileBytes []byte, embedding string) (bool, float64, error) {
	// 1) Создаём multipart/form-data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Поле "image"
	w, err := writer.CreateFormFile("image", "face.jpg")
	if err != nil {
		return false, 0.0, err
	}
	_, err = w.Write(fileBytes)
	if err != nil {
		return false, 0.0, err
	}

	// Поле "embedding" (JSON-строка)
	err = writer.WriteField("embedding", embedding)
	if err != nil {
		return false, 0.0, err
	}

	writer.Close()

	// 2) Отправляем HTTP-запрос
	faceServiceURL := "http://localhost:8000/compare" // <-- адрес в Docker-сети
	req, err := http.NewRequest(http.MethodPost, faceServiceURL, body)
	if err != nil {
		return false, 0.0, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, 0.0, err
	}
	defer resp.Body.Close()

	// 3) Разбираем ответ
	var result struct {
		Match    bool    `json:"match"`
		Distance float64 `json:"distance"`
		Error    string  `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, 0.0, err
	}

	logrus.Info(result)
	if resp.StatusCode != http.StatusOK && !result.Match {
		// Например, 400/401
		return false, result.Distance, fmt.Errorf("face compare failed: %s", result.Error)
	}

	return result.Match, result.Distance, nil
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

	if existingToken == nil {
		return s.createEmployeeToken(employeeID, false)
	}

	if err := s.employeeTokenManager.DeleteTokenByEmployeeID(employeeID); err != nil {
		return "", fmt.Errorf("failed to delete old/expired token: %w", err)
	}

	return s.createEmployeeToken(employeeID, false)
}

func (s *authenticationService) createEmployeeToken(employeeID uint, mfa bool) (string, error) {
	sessionToken, err := types.GenerateEmployeeJWT(employeeID, mfa)
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
