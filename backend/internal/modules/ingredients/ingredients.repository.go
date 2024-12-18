package ingredients

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type IngredientRepository interface {
	CreateIngredient(ingredient *data.Ingredient) error
	UpdateIngredient(ingredient *data.Ingredient) error
	DeleteIngredient(ingredientID uint) error
	GetIngredientByID(ingredientID uint) (*data.Ingredient, error)
	GetIngredients(filter types.IngredientFilter, pagination *utils.Pagination) ([]data.Ingredient, error)
}

type ingredientRepository struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) IngredientRepository {
	return &ingredientRepository{db: db}
}

func (r *ingredientRepository) CreateIngredient(ingredient *data.Ingredient) error {
	return r.db.Create(ingredient).Error
}

func (r *ingredientRepository) UpdateIngredient(ingredient *data.Ingredient) error {
	return r.db.Save(ingredient).Error
}

func (r *ingredientRepository) DeleteIngredient(ingredientID uint) error {
	return r.db.Where("id = ?", ingredientID).Delete(&data.Ingredient{}).Error
}

func (r *ingredientRepository) GetIngredientByID(ingredientID uint) (*data.Ingredient, error) {
	var ingredient data.Ingredient
	if err := r.db.First(&ingredient, ingredientID).Error; err != nil {
		return nil, err
	}
	return &ingredient, nil
}

func (r *ingredientRepository) GetIngredients(filter types.IngredientFilter, pagination *utils.Pagination) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient
	query := r.db.Model(&data.Ingredient{})

	// Apply filtering
	if filter.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
	}
	if filter.MinCalories != nil {
		query = query.Where("calories >= ?", *filter.MinCalories)
	}
	if filter.MaxCalories != nil {
		query = query.Where("calories <= ?", *filter.MaxCalories)
	}

	// Apply pagination
	query, err := utils.ApplyPagination(query, pagination, data.Ingredient{})
	if err != nil {
		return nil, err
	}

	// Execute query
	if err := query.Find(&ingredients).Error; err != nil {
		return nil, err
	}

	return ingredients, nil
}
