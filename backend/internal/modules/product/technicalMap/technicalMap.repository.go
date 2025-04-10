package technicalMap

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/technicalMap/types"
	"gorm.io/gorm"
)

type TechnicalMapRepository interface {
	GetProductSizeTechnicalMapByID(productSizeID uint) ([]data.ProductSizeIngredient, error)
}

type technicalMapRepository struct {
	db *gorm.DB
}

func NewTechnicalMapRepository(db *gorm.DB) TechnicalMapRepository {
	return &technicalMapRepository{db: db}
}

func (r *technicalMapRepository) GetProductSizeTechnicalMapByID(productSizeID uint) ([]data.ProductSizeIngredient, error) {
	var productSizeIngredients []data.ProductSizeIngredient
	if err := r.db.Where("product_size_id = ?", productSizeID).
		Preload("Ingredient").
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory").
		Find(&productSizeIngredients).Error; err != nil {
		return nil, err
	}

	if len(productSizeIngredients) == 0 {
		return nil, types.ErrorTechnicalMapNotFound
	}

	return productSizeIngredients, nil
}
