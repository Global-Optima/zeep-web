package technicalMap

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/technicalMap/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type TechnicalMapService interface {
	GetProductSizeTechnicalMapByID(productSizeID uint) (*types.ProductSizeTechnicalMap, error)
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

func (s *technicalMapService) GetProductSizeTechnicalMapByID(productSizeID uint) (*types.ProductSizeTechnicalMap, error) {
	productSizeIngredients, err := s.repository.GetProductSizeTechnicalMapByID(productSizeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get technical map for product size: ", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	technicalMap := types.ConvertToProductSizeTechnicalMapDTO(productSizeIngredients)
	return &technicalMap, nil
}
