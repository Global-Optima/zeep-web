package ingredients

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type IngredientRepository interface {
	CreateIngredient(ingredient *data.Ingredient) (uint, error)
	UpdateIngredient(ingredientID uint, ingredient *data.Ingredient) error
	DeleteIngredient(ingredientID uint) error
	GetIngredientByID(ingredientID uint) (*data.Ingredient, error)
	GetIngredientsByIDs(ingredientIDs []uint) ([]data.Ingredient, error)
	GetIngredients(filter *types.IngredientFilter) ([]data.Ingredient, error)
	GetIngredientsForProductSizes(productSizeIDs []uint) ([]data.Ingredient, error)
	GetIngredientsForAdditives(additiveIDs []uint) ([]data.Ingredient, error)
}

type ingredientRepository struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) IngredientRepository {
	return &ingredientRepository{db: db}
}

func (r *ingredientRepository) CreateIngredient(ingredient *data.Ingredient) (uint, error) {
	err := r.db.Create(ingredient).Error
	if err != nil {
		return 0, err
	}
	return ingredient.ID, err
}

func (r *ingredientRepository) UpdateIngredient(ingredientID uint, ingredient *data.Ingredient) error {
	return r.db.Model(&data.Ingredient{}).
		Where("id = ?", ingredientID).
		Updates(ingredient).Error
}

func (r *ingredientRepository) DeleteIngredient(ingredientID uint) error {
	return r.db.Where("id = ?", ingredientID).Delete(&data.Ingredient{}).Error
}

func (r *ingredientRepository) GetIngredientByID(ingredientID uint) (*data.Ingredient, error) {
	var ingredient data.Ingredient
	err := r.db.Preload("Unit").
		Preload("IngredientCategory").
		First(&ingredient, ingredientID).Error

	if err != nil {
		return nil, err
	}

	return &ingredient, nil
}

func (r *ingredientRepository) GetIngredientsByIDs(ingredientIDs []uint) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient
	err := r.db.Model(&data.Ingredient{}).
		Preload("Unit").
		Preload("IngredientCategory").
		Where("id IN (?)", ingredientIDs).
		Find(&ingredients).Error

	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (r *ingredientRepository) GetIngredients(filter *types.IngredientFilter) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient
	query := r.db.Model(&data.Ingredient{}).Preload("Unit").Preload("IngredientCategory")

	// Apply filtering
	if filter.ProductSizeID != nil {
		query = query.Joins("JOIN product_size_ingredients psi ON psi.ingredient_id = ingredients.id").
			Where("psi.product_size_id = ?", *filter.ProductSizeID)
	}

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
	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Ingredient{})
	if err != nil {
		return nil, err
	}

	// Execute query
	if err := query.Find(&ingredients).Error; err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (r *ingredientRepository) GetIngredientsForProductSizes(productSizeIDs []uint) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient
	query := r.db.Model(&data.Ingredient{}).
		Select("DISTINCT ingredients.*").
		Preload("Unit").
		Preload("IngredientCategory").
		Joins("JOIN product_size_ingredients psi ON psi.ingredient_id = ingredients.id").
		Joins("JOIN product_sizes ps ON psi.product_size_id = ps.id").
		Where("ps.id IN ?", productSizeIDs)

	if err := query.Find(&ingredients).Error; err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (r *ingredientRepository) GetIngredientsForAdditives(additiveIDs []uint) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient
	query := r.db.Model(&data.Ingredient{}).
		Select("DISTINCT ingredients.*").
		Preload("Unit").
		Preload("IngredientCategory").
		Joins("JOIN additive_ingredients ai ON ai.ingredient_id = ingredients.id").
		Joins("JOIN additives a ON ai.additive_id = a.id").
		Where("a.id IN ?", additiveIDs)

	if err := query.Find(&ingredients).Error; err != nil {
		return nil, err
	}

	return ingredients, nil
}
