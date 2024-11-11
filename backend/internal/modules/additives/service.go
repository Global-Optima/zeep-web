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

	var additiveCategoryDTOs []types.AdditiveCategoryDTO
	for _, category := range additiveCategories {
		additives := make([]types.AdditiveDTO, 0)
		for _, additive := range category.Additives {
			additives = append(additives, types.AdditiveDTO{
				ID:          additive.ID,
				Name:        additive.Name,
				Description: additive.Description,
				Price:       additive.StorePrice,
				ImageURL:    additive.ImageURL,
			})
		}
		if len(additives) > 0 {
			additiveCategoryDTOs = append(additiveCategoryDTOs, types.AdditiveCategoryDTO{
				ID:        category.ID,
				Name:      category.Name,
				Additives: additives,
			})
		}
	}

	return additiveCategoryDTOs, nil
}
