package technicalMap

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/technicalMap/types"
)

type TechnicalMapService interface {
	GetAdditiveTechnicalMapByID(AdditiveID uint) (*types.AdditiveTechnicalMap, error)
}

type technicalMapService struct {
	repository TechnicalMapRepository
}

func NewTechnicalMapService(repository TechnicalMapRepository) TechnicalMapService {
	return &technicalMapService{repository: repository}
}

func (s *technicalMapService) GetAdditiveTechnicalMapByID(AdditiveID uint) (*types.AdditiveTechnicalMap, error) {
	AdditiveIngredients, err := s.repository.GetAdditiveTechnicalMapByID(AdditiveID)
	if err != nil {
		return nil, err
	}

	technicalMap := types.ConvertToAdditiveTechnicalMapDTO(AdditiveIngredients)
	return &technicalMap, nil
}
