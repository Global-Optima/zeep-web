package recipes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/recipes/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	CreateOrReplaceRecipeStepsByProductID(productID uint, recipeSteps []data.RecipeStep) ([]uint, error)
	GetRecipeStepByID(id uint) (*data.RecipeStep, error)
	GetRecipeStepsByProductID(productID uint) ([]data.RecipeStep, error)
	DeleteRecipeStepsByProductID(productID uint) error
	HasRecipes(productID uint) (bool, error)
}

type recipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) CreateOrReplaceRecipeStepsByProductID(productID uint, recipeSteps []data.RecipeStep) ([]uint, error) {
	if len(recipeSteps) == 0 {
		return nil, types.ErrNothingToUpdate
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&data.RecipeStep{}).
			Where("product_id = ?", productID).
			Delete(&recipeSteps).Error

		if err != nil {
			return err
		}

		if err := tx.Create(&recipeSteps).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(recipeSteps))
	for i, recipeStep := range recipeSteps {
		ids[i] = recipeStep.ID
	}
	return ids, nil
}

func (r *recipeRepository) GetRecipeStepByID(id uint) (*data.RecipeStep, error) {
	var recipeStep data.RecipeStep
	err := r.db.Model(&data.RecipeStep{}).
		Where("id = ?", id).
		First(&recipeStep).Error

	if err != nil {
		return nil, err
	}

	return &recipeStep, nil
}

func (r *recipeRepository) GetRecipeStepsByProductID(productID uint) ([]data.RecipeStep, error) {
	var recipeSteps []data.RecipeStep
	if err := r.db.Where("product_id = ?", productID).Find(&recipeSteps).Error; err != nil {
		return nil, err
	}

	return recipeSteps, nil
}

func (r *recipeRepository) DeleteRecipeStepsByProductID(productID uint) error {
	if err := r.db.Where("product_id = ?", productID).
		Delete(&data.RecipeStep{}).
		Error; err != nil {
		return err
	}
	return nil
}

func (r *recipeRepository) HasRecipes(productID uint) (bool, error) {
	var recipeStep data.RecipeStep
	err := r.db.Model(&data.RecipeStep{}).
		Where("product_id = ?", productID).
		First(&recipeStep).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
