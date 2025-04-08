package technicalMap

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/technicalMap/types"
	"gorm.io/gorm"
)

type TechnicalMapRepository interface {
	GetProvisionTechnicalMapByID(ProvisionID uint) ([]data.ProvisionIngredient, error)
}

type technicalMapRepository struct {
	db *gorm.DB
}

func NewTechnicalMapRepository(db *gorm.DB) TechnicalMapRepository {
	return &technicalMapRepository{db: db}
}

func (r *technicalMapRepository) GetProvisionTechnicalMapByID(provisionID uint) ([]data.ProvisionIngredient, error) {
	var ProvisionIngredients []data.ProvisionIngredient
	if err := r.db.Where("provision_id = ?", provisionID).
		Preload("Ingredient").
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory").
		Find(&ProvisionIngredients).Error; err != nil {
		return nil, err
	}

	if len(ProvisionIngredients) == 0 {
		return nil, types.ErrorTechnicalMapNotFound
	}

	return ProvisionIngredients, nil
}
