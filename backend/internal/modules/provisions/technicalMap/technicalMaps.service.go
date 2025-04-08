package technicalMap

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/technicalMap/types"
)

type TechnicalMapService interface {
	GetProvisionTechnicalMapByID(ProvisionID uint) (*types.ProvisionTechnicalMap, error)
}

type technicalMapService struct {
	repository TechnicalMapRepository
}

func NewTechnicalMapService(repository TechnicalMapRepository) TechnicalMapService {
	return &technicalMapService{repository: repository}
}

func (s *technicalMapService) GetProvisionTechnicalMapByID(ProvisionID uint) (*types.ProvisionTechnicalMap, error) {
	ProvisionIngredients, err := s.repository.GetProvisionTechnicalMapByID(ProvisionID)
	if err != nil {
		return nil, err
	}

	technicalMap := types.ConvertToProvisionTechnicalMapDTO(ProvisionIngredients)
	return &technicalMap, nil
}
