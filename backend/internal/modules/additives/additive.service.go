package additives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type AdditiveService interface {
	GetAdditiveCategories(filter types.AdditiveCategoriesFilterQuery) ([]types.AdditiveCategoryDTO, error)
	GetAdditives(filter types.AdditiveFilterQuery) ([]types.AdditiveDTO, error)
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

func (s *additiveService) GetAdditiveCategories(filter types.AdditiveCategoriesFilterQuery) ([]types.AdditiveCategoryDTO, error) {
	// Fetch raw data from the repository
	categories, err := s.repo.GetAdditiveCategories(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve additives", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	// Handle case where no categories are found
	if len(categories) == 0 {
		return []types.AdditiveCategoryDTO{}, nil
	}

	// Convert raw data into DTOs
	var categoryDTOs []types.AdditiveCategoryDTO
	for _, category := range categories {
		// Initialize additives list explicitly as an empty slice
		additives := make([]types.AdditiveCategoryItemDTO, 0)

		// Populate additives if present
		for _, additive := range category.Additives {
			additives = append(additives, types.AdditiveCategoryItemDTO{
				ID:          additive.ID,
				Name:        additive.Name,
				Description: additive.Description,
				Price:       additive.BasePrice,
				ImageURL:    additive.ImageURL,
				Size:        additive.Size,
				CategoryID:  category.ID,
			})
		}

		// Append the category DTO with additives list (empty or populated)
		categoryDTOs = append(categoryDTOs, types.AdditiveCategoryDTO{
			ID:               category.ID,
			Name:             category.Name,
			Description:      category.Description,
			IsMultipleSelect: category.IsMultipleSelect,
			Additives:        additives, // Always initialized as a slice
		})
	}

	return categoryDTOs, nil
}

func (s *additiveService) GetAdditives(filter types.AdditiveFilterQuery) ([]types.AdditiveDTO, error) {
	additives, err := s.repo.GetAdditives(filter)
	if err != nil {
		return nil, err
	}

	var additiveDTOs []types.AdditiveDTO
	for _, additive := range additives {
		additiveDTOs = append(additiveDTOs, types.AdditiveDTO{
			ID:          additive.ID,
			Name:        additive.Name,
			Description: additive.Description,
			Price:       additive.BasePrice,
			ImageURL:    additive.ImageURL,
			Size:        additive.Size,
			Category: struct {
				ID               uint   `json:"id"`
				Name             string `json:"name"`
				IsMultipleSelect bool   `json:"isMultipleSelect"`
			}{
				ID:               additive.Category.ID,
				Name:             additive.Category.Name,
				IsMultipleSelect: additive.Category.IsMultipleSelect,
			},
		})
	}

	return additiveDTOs, nil
}
