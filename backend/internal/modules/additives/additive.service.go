package additives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
)

type AdditiveService interface {
	GetAdditivesByStoreAndProductSize(productID uint) ([]types.AdditiveCategoryDTO, error)
}

type additiveService struct {
	repo AdditiveRepository
}

func NewAdditiveService(repo AdditiveRepository) AdditiveService {
	return &additiveService{repo: repo}
}

func (s *additiveService) GetAdditivesByStoreAndProductSize(productSizeID uint) ([]types.AdditiveCategoryDTO, error) {
	additiveCategories, err := s.repo.GetAdditiveCategoriesByProductSize(productSizeID)
	if err != nil {
		return nil, err
	}

	if additiveCategories == nil {
		return []types.AdditiveCategoryDTO{}, nil
	}

	return additiveCategories, nil
}
