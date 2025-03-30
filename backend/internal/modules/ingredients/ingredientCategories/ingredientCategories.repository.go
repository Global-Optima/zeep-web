package ingredientCategories

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type IngredientCategoryRepository interface {
	Create(category *data.IngredientCategory) error
	GetByID(id uint) (*data.IngredientCategory, error)
	Update(id uint, updates data.IngredientCategory) error
	Delete(id uint) error
	GetAll(filter *types.IngredientCategoryFilter) ([]data.IngredientCategory, error)
}

type ingredientCategoryRepository struct {
	db *gorm.DB
}

func NewIngredientCategoryRepository(db *gorm.DB) IngredientCategoryRepository {
	return &ingredientCategoryRepository{db: db}
}

func (r *ingredientCategoryRepository) Create(category *data.IngredientCategory) error {
	return r.db.Create(category).Error
}

func (r *ingredientCategoryRepository) GetByID(id uint) (*data.IngredientCategory, error) {
	var category data.IngredientCategory
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *ingredientCategoryRepository) Update(id uint, updates data.IngredientCategory) error {
	return r.db.Model(&data.IngredientCategory{}).Where("id = ?", id).
		Updates(&data.IngredientCategory{
			Name:        updates.Name,
			Description: updates.Description,
		}).Error
}

func (r *ingredientCategoryRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := checkIngredientCategoryReferences(tx, id); err != nil {
			return err
		}

		if err := tx.Delete(&data.IngredientCategory{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func checkIngredientCategoryReferences(db *gorm.DB, ingredientCategoryID uint) error {
	var ingredientCategory data.IngredientCategory

	err := db.
		Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Where(&data.IngredientCategory{BaseEntity: data.BaseEntity{ID: ingredientCategoryID}}).
		First(&ingredientCategory).Error
	if err != nil {
		return err
	}

	if len(ingredientCategory.Ingredients) > 0 {
		return types.ErrIngredientCategoryIsInUse
	}

	return nil
}

func (r *ingredientCategoryRepository) GetAll(filter *types.IngredientCategoryFilter) ([]data.IngredientCategory, error) {
	var categories []data.IngredientCategory

	query := r.db.Model(&data.IngredientCategory{})

	if filter.Search != nil {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+*filter.Search+"%", "%"+*filter.Search+"%")
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &categories)
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	err = query.Find(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ingredient categories: %w", err)
	}

	return categories, nil
}
