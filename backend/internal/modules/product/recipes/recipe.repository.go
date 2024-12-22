package recipes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	CreateRecipeStep(recipeStep *data.RecipeStep) (uint, error)
	UpdateRecipeStep(id uint, recipeStep *data.RecipeStep) error
	GetRecipeStepByID(id uint) (*data.RecipeStep, error)
	GetRecipeStepsByProductID(productID uint) ([]data.RecipeStep, error)
	DeleteRecipeStep(recipeStep *data.RecipeStep) error
}

type recipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) CreateRecipeStep(recipeStep *data.RecipeStep) (uint, error) {
	if err := r.db.Create(&recipeStep).Error; err != nil {
		return 0, err
	}
	return recipeStep.ID, nil
}

func (r *recipeRepository) UpdateRecipeStep(id uint, recipeStep *data.RecipeStep) error {
	err := r.db.Model(&data.RecipeStep{}).
		Where("id = ?", id).
		Updates(&recipeStep).Error

	if err != nil {
		return err
	}

	return nil
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

func (r *recipeRepository) DeleteRecipeStep(recipeStep *data.RecipeStep) error {
	if err := r.db.Delete(recipeStep).Error; err != nil {
		return err
	}
	return nil
}
