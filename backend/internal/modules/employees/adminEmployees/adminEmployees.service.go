package employees

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/adminEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminEmployeeService interface {
	CreateAdminEmployee(dto *employeesTypes.CreateEmployeeDTO, image *multipart.FileHeader) (uint, error)
	GetAdminEmployees(filter *employeesTypes.EmployeesFilter) ([]types.AdminEmployeeDTO, error)
	GetAdminEmployeeByID(id uint) (*types.AdminEmployeeDetailsDTO, error)
	GetAllAdminEmployees() ([]employeesTypes.EmployeeAccountDTO, error)
}

type adminEmployeeService struct {
	repo         AdminEmployeeRepository
	employeeRepo employees.EmployeeRepository
	logger       *zap.SugaredLogger
}

func NewAdminEmployeeService(repo AdminEmployeeRepository, employeeRepo employees.EmployeeRepository, logger *zap.SugaredLogger) AdminEmployeeService {
	return &adminEmployeeService{
		repo:         repo,
		employeeRepo: employeeRepo,
		logger:       logger,
	}
}

func (s *adminEmployeeService) CreateAdminEmployee(input *employeesTypes.CreateEmployeeDTO, image *multipart.FileHeader) (uint, error) {

	embedding, err := s.callFaceServiceExtract(image)
	if err != nil {
		s.logger.Errorf("failed to extract face embedding: %v", err)
		return 0, err
	}

	employee, err := types.CreateToAdminEmployee(input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create admin employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	employee.FaceEmbedding = embedding

	existingEmployee, err := s.employeeRepo.GetEmployeeByEmailOrPhone(employee.Email, employee.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr := utils.WrapError("error checking employee uniqueness", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	if existingEmployee != nil {
		return 0, moduleErrors.ErrAlreadyExists
	}

	id, err := s.employeeRepo.CreateEmployee(employee)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create admin employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *adminEmployeeService) callFaceServiceExtract(fileHeader *multipart.FileHeader) (string, error) {
	// Считываем bytes из fileHeader
	fileBytes, err := readFileHeader(fileHeader)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Создаём multipart/form-data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем поле "image"
	part, err := writer.CreateFormFile("image", fileHeader.Filename)
	if err != nil {
		return "", err
	}
	_, err = part.Write(fileBytes)
	if err != nil {
		return "", err
	}

	writer.Close()

	// Отправляем POST-запрос на /extract
	// Предположим, сервис доступен внутри сети Docker по URL:
	faceServiceURL := "http://localhost:8000/extract" // см. docker-compose

	req, err := http.NewRequest(http.MethodPost, faceServiceURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Разбираем ответ
	var result struct {
		Success   bool      `json:"success"`
		Embedding []float64 `json:"embedding"` // массив чисел
		Error     string    `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if !result.Success {
		return "", fmt.Errorf("face_service error: %s", result.Error)
	}

	// сериализуем embedding в JSON (или другой формат) для хранения
	embeddingBytes, err := json.Marshal(result.Embedding)
	if err != nil {
		return "", err
	}

	return string(embeddingBytes), nil
}

// вспомогательная функция чтения *multipart.FileHeader
func readFileHeader(fileHeader *multipart.FileHeader) ([]byte, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

func (s *adminEmployeeService) GetAdminEmployees(filter *employeesTypes.EmployeesFilter) ([]types.AdminEmployeeDTO, error) {
	adminEmployees, err := s.repo.GetAdminEmployees(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve admin employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	dtos := make([]types.AdminEmployeeDTO, len(adminEmployees))
	for i, employee := range adminEmployees {
		dtos[i] = *types.MapToAdminEmployeeDTO(&employee)
	}
	return dtos, nil
}

func (s *adminEmployeeService) GetAllAdminEmployees() ([]employeesTypes.EmployeeAccountDTO, error) {
	adminEmployees, err := s.repo.GetAllAdminEmployees()
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve all admin employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]employeesTypes.EmployeeAccountDTO, len(adminEmployees))
	for i, adminEmployee := range adminEmployees {
		dtos[i] = *employeesTypes.MapToEmployeeAccountDTO(&adminEmployee.Employee)
	}

	return dtos, nil
}

func (s *adminEmployeeService) GetAdminEmployeeByID(id uint) (*types.AdminEmployeeDetailsDTO, error) {
	if id == 0 {
		return nil, errors.New("invalid admin employee ID")
	}

	employee, err := s.repo.GetAdminEmployeeByID(id)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve admin employee with ID = %d: %w", id, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToAdminEmployeeDetailsDTO(employee), nil
}
