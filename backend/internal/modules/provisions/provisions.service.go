package provisions

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type ProvisionService interface {
	GetProvisionByID(provisionID uint) (*types.ProvisionDetailsDTO, error)
	GetProvisions(filter *types.ProvisionFilterDTO) ([]types.ProvisionDTO, error)
	CreateProvision(dto *types.CreateProvisionDTO) (uint, error)
	UpdateProvision(provisionID uint, dto *types.UpdateProvisionDTO) (*types.ProvisionDTO, error)
	DeleteProvision(provisionID uint) (*data.Provision, error)
}

type provisionService struct {
	repo   ProvisionRepository
	logger *zap.SugaredLogger
}

func NewProvisionService(repo ProvisionRepository, logger *zap.SugaredLogger) ProvisionService {
	return &provisionService{
		repo:   repo,
		logger: logger,
	}
}

func (s *provisionService) GetProvisionByID(provisionID uint) (*types.ProvisionDetailsDTO, error) {
	provision, err := s.repo.GetProvisionWithDetailsByID(provisionID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve provision ID %d: %w", provisionID, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	return types.MapToProvisionDetailsDTO(provision), nil
}

func (s *provisionService) GetProvisions(filter *types.ProvisionFilterDTO) ([]types.ProvisionDTO, error) {
	provisions, err := s.repo.GetProvisions(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve provisions", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.ProvisionDTO, len(provisions))
	for i, provision := range provisions {
		dtos[i] = *types.MapToProvisionDTO(&provision)
	}

	return dtos, nil
}

func (s *provisionService) CreateProvision(dto *types.CreateProvisionDTO) (uint, error) {
	exists, err := s.repo.CheckProvisionExists(dto.Name)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to check provision existence: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	if exists {
		wrappedErr := fmt.Errorf("%w: provision with name %s already exists", types.ErrProvisionUniqueName, dto.Name)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	provision := types.CreateToProvisionModel(dto)
	id, err := s.repo.CreateProvision(provision)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create provision: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *provisionService) UpdateProvision(provisionID uint, dto *types.UpdateProvisionDTO) (*types.ProvisionDTO, error) {
	existingProvision, err := s.repo.GetProvisionByID(provisionID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch provision: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	oldProvision := *existingProvision
	updateModels, err := types.UpdateToProvisionModels(existingProvision, dto)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update provision: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	err = s.repo.SaveProvisionWithAssociations(updateModels)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update provision and its associations: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return types.MapToProvisionDTO(&oldProvision), nil
}

func (s *provisionService) DeleteProvision(provisionID uint) (*data.Provision, error) {
	provision, err := s.repo.DeleteProvision(provisionID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete provision: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	return provision, nil
}
