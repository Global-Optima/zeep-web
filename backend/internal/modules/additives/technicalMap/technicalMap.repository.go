package technicalMap

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/technicalMap/types"
	"gorm.io/gorm"
)

type TechnicalMapRepository interface {
	GetAdditiveTechnicalMapByID(AdditiveID uint) ([]data.AdditiveIngredient, error)
}

type technicalMapRepository struct {
	db *gorm.DB
}

func NewTechnicalMapRepository(db *gorm.DB) TechnicalMapRepository {
	return &technicalMapRepository{db: db}
}

func (r *technicalMapRepository) GetAdditiveTechnicalMapByID(AdditiveID uint) ([]data.AdditiveIngredient, error) {
	var AdditiveIngredients []data.AdditiveIngredient
	if err := r.db.Where("additive_id = ?", AdditiveID).
		Preload("Ingredient").
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory").
		Find(&AdditiveIngredients).Error; err != nil {
		return nil, err
	}

	if len(AdditiveIngredients) == 0 {
		return nil, types.ErrorTechnicalMapNotFound
	}

	return AdditiveIngredients, nil
}
