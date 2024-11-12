package additives

import "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"

type AdditiveService interface {
	GetAdditivesByStoreAndProduct(storeID uint, productID uint) ([]types.AdditiveCategoryDTO, error)
}

type additiveService struct {
	repo AdditiveRepository
}

func NewAdditiveService(repo AdditiveRepository) AdditiveService {
	return &additiveService{repo: repo}
}

func (s *additiveService) GetAdditivesByStoreAndProduct(storeID uint, productID uint) ([]types.AdditiveCategoryDTO, error) {
	additiveCategories, err := s.repo.GetAdditivesByStoreAndProduct(storeID, productID)
	if err != nil {
		return nil, err
	}

	return additiveCategories, nil
}
