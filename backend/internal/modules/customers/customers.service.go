package customers

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/customers/types"
	"go.uber.org/zap"
)

type CustomerService interface {
	GetCustomerById(id uint) (*types.CustomerAdminDTO, error)
}

type customerService struct {
	repo   CustomerRepository
	logger *zap.SugaredLogger
}

func NewCustomerService(repo CustomerRepository, logger *zap.SugaredLogger) CustomerService {
	return &customerService{
		repo:   repo,
		logger: logger,
	}
}

func (s *customerService) GetCustomerById(customerID uint) (*types.CustomerAdminDTO, error) {
	customer, err := s.repo.GetCustomerByID(customerID)
	if err != nil {
		wrappedErr := fmt.Errorf("error retrieving customer with ID = %d: %w", customerID, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return types.MapToCustomerDTO(customer), nil
}
