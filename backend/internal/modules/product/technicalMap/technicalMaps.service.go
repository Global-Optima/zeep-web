package technicalMap

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/technicalMap/types"
)

type TechnicalMapService interface {
	GetProductSizeTechnicalMapByID(productSizeID uint) (*types.ProductSizeTechnicalMap, error)
}

type technicalMapService struct {
	repository TechnicalMapRepository
}

func NewTechnicalMapService(repository TechnicalMapRepository) TechnicalMapService {
	return &technicalMapService{repository: repository}
}

func (s *technicalMapService) GetProductSizeTechnicalMapByID(productSizeID uint) (*types.ProductSizeTechnicalMap, error) {
	productSizeIngredients, err := s.repository.GetProductSizeTechnicalMapByID(productSizeID)
	if err != nil {
		return nil, err
	}

	technicalMap := types.ConvertToProductSizeTechnicalMapDTO(productSizeIngredients)
	return &technicalMap, nil
}
