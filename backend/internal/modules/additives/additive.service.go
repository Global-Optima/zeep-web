package additives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type AdditiveService interface {
	GetAdditivesByStoreAndProductSize(productID uint) ([]types.AdditiveCategoryDTO, error)
}

type additiveService struct {
	repo   AdditiveRepository
	logger *zap.SugaredLogger
}

func NewAdditiveService(repo AdditiveRepository, logger *zap.SugaredLogger) AdditiveService {
	return &additiveService{
		repo:   repo,
		logger: logger,
	}
}

func (s *additiveService) GetAdditivesByStoreAndProductSize(productSizeID uint) ([]types.AdditiveCategoryDTO, error) {
	additiveCategories, err := s.repo.GetAdditiveCategoriesByProductSize(productSizeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve additives", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if additiveCategories == nil {
		return []types.AdditiveCategoryDTO{}, nil
	}

	return additiveCategories, nil
}
