package customers

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/customers/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type CustomerService interface {
	GetCustomerById(id uint) (*types.CustomerDTO, error)
}

type customerService struct {
	repo   CustomerRepository
	logger *zap.SugaredLogger
}

func NewEmployeeService(repo CustomerRepository, logger *zap.SugaredLogger) CustomerService {
	return &customerService{
		repo:   repo,
		logger: logger,
	}
}

func (s *customerService) GetCustomerById(customerID uint) (*types.CustomerDTO, error) {
	customer, err := s.repo.GetCustomerByID(customerID)
	if err != nil {
		wrappedErr := utils.WrapError("error retrieving customer", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	
	return types.MapToCustomerDTO(customer), nil
}
