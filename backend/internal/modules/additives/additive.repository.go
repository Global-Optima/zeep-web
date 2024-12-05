package additives

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"gorm.io/gorm"
)

type AdditiveRepository interface {
	GetAdditiveCategoriesByProductSize(productSizeID uint) ([]types.AdditiveCategoryDTO, error)
	GetAdditiveByID(additiveID uint) (*data.Additive, error)
}

type additiveRepository struct {
	db *gorm.DB
}

func NewAdditiveRepository(db *gorm.DB) AdditiveRepository {
	return &additiveRepository{db: db}
}

func (r *additiveRepository) GetAdditiveCategoriesByProductSize(productSizeID uint) ([]types.AdditiveCategoryDTO, error) {
	var productAdditives []data.ProductAdditive

	err := r.db.
		Preload("Additive").
		Where("product_size_id = ?", productSizeID).
		Find(&productAdditives).
		Error

	if err != nil {
		return nil, err
	}

	categoryMap := make(map[uint]*types.AdditiveCategoryDTO)

	for _, pa := range productAdditives {
		additive := pa.Additive

		categoryID := additive.AdditiveCategoryID

		if _, exists := categoryMap[categoryID]; !exists {
			category := data.AdditiveCategory{}
			err := r.db.
				Where("id = ?", categoryID).
				First(&category).
				Error
			if err != nil {
				return nil, err
			}

			categoryMap[categoryID] = &types.AdditiveCategoryDTO{
				ID:               category.ID,
				Name:             category.Name,
				IsMultipleSelect: category.IsMultipleSelect,
				Additives:        []types.AdditiveDTO{},
			}
		}

		categoryMap[categoryID].Additives = append(categoryMap[categoryID].Additives, types.AdditiveDTO{
			ID:          additive.ID,
			Name:        additive.Name,
			Description: additive.Description,
			Price:       pa.Additive.StorePrice,
			ImageURL:    additive.ImageURL,
			CategoryId:  categoryID,
		})
	}

	var additiveCategories []types.AdditiveCategoryDTO
	for _, category := range categoryMap {
		additiveCategories = append(additiveCategories, *category)
	}

	return additiveCategories, nil
}

func (r *additiveRepository) GetAdditiveByID(additiveID uint) (*data.Additive, error) {
	var additive data.Additive
	err := r.db.Where("id = ?", additiveID).First(&additive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("additive with ID %d not found", additiveID)
		}
		return nil, fmt.Errorf("failed to fetch additive with ID %d: %w", additiveID, err)
	}
	return &additive, nil
}
