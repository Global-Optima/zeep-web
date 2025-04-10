package technicalMap

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/technicalMap/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type TechnicalMapService interface {
	GetProvisionTechnicalMapByID(ProvisionID uint) (*types.ProvisionTechnicalMap, error)
}

type technicalMapService struct {
	repository TechnicalMapRepository
	logger     *zap.SugaredLogger
}

func NewTechnicalMapService(repository TechnicalMapRepository, logger *zap.SugaredLogger) TechnicalMapService {
	return &technicalMapService{
		repository: repository,
		logger:     logger,
	}
}

func (s *technicalMapService) GetProvisionTechnicalMapByID(ProvisionID uint) (*types.ProvisionTechnicalMap, error) {
	ProvisionIngredients, err := s.repository.GetProvisionTechnicalMapByID(ProvisionID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get technical map for provision: ", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	technicalMap := types.ConvertToProvisionTechnicalMapDTO(ProvisionIngredients)
	return &technicalMap, nil
}
