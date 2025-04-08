package technicalMap

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/technicalMap/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type TechnicalMapService interface {
	GetAdditiveTechnicalMapByID(AdditiveID uint) (*types.AdditiveTechnicalMap, error)
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
func (s *technicalMapService) GetAdditiveTechnicalMapByID(AdditiveID uint) (*types.AdditiveTechnicalMap, error) {
	AdditiveIngredients, err := s.repository.GetAdditiveTechnicalMapByID(AdditiveID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get technical map for additive: ", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	technicalMap := types.ConvertToAdditiveTechnicalMapDTO(AdditiveIngredients)
	return &technicalMap, nil
}
