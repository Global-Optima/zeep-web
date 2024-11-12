package additives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"gorm.io/gorm"
)

type AdditiveRepository interface {
	GetAdditivesByStoreAndProduct(storeID uint, productID uint) ([]types.AdditiveCategoryDTO, error)
}

type additiveRepository struct {
	db *gorm.DB
}

func NewAdditiveRepository(db *gorm.DB) AdditiveRepository {
	return &additiveRepository{db: db}
}

func (r *additiveRepository) GetAdditivesByStoreAndProduct(storeID uint, productID uint) ([]types.AdditiveCategoryDTO, error) {
	var categories []data.AdditiveCategory

	err := r.db.Preload("Additives", func(db *gorm.DB) *gorm.DB {
		return db.Joins("LEFT JOIN store_additives sa ON sa.additive_id = additives.id AND sa.store_id = ?", storeID).
			Select("additives.*, COALESCE(NULLIF(sa.price, 0), additives.base_price) AS store_price")
	}).Find(&categories).Error

	if err != nil {
		return nil, err
	}

	var categoryDTOs []types.AdditiveCategoryDTO
	for _, category := range categories {
		categoryDTO := types.AdditiveCategoryDTO{
			ID:   category.ID,
			Name: category.Name,
		}

		for _, additive := range category.Additives {
			additiveDTO := types.AdditiveDTO{
				ID:          additive.ID,
				Name:        additive.Name,
				Description: additive.Description,
				Price:       additive.StorePrice,
				ImageURL:    additive.ImageURL,
			}
			categoryDTO.Additives = append(categoryDTO.Additives, additiveDTO)
		}

		if len(categoryDTO.Additives) > 0 {
			categoryDTOs = append(categoryDTOs, categoryDTO)
		}
	}

	return categoryDTOs, nil
}
