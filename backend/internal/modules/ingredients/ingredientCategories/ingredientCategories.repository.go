package ingredientCategories

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type IngredientCategoryRepository interface {
	Create(category *data.IngredientCategory) error
	GetByID(id uint) (*data.IngredientCategory, error)
	GetTranslatedByID(locale data.LanguageCode, id uint) (*data.IngredientCategory, error)
	Update(id uint, updates *data.IngredientCategory) error
	Delete(id uint) error
	GetAll(locale data.LanguageCode, filter *types.IngredientCategoryFilter) ([]data.IngredientCategory, error)

	CloneWithTransaction(tx *gorm.DB) IngredientCategoryRepository
	FindRawIngredientCategoryByID(id uint, category *data.IngredientCategory) error
}

type ingredientCategoryRepository struct {
	db *gorm.DB
}

func NewIngredientCategoryRepository(db *gorm.DB) IngredientCategoryRepository {
	return &ingredientCategoryRepository{db: db}
}

func (r *ingredientCategoryRepository) CloneWithTransaction(tx *gorm.DB) IngredientCategoryRepository {
	return &ingredientCategoryRepository{db: tx}
}
func (r *ingredientCategoryRepository) FindRawIngredientCategoryByID(id uint, category *data.IngredientCategory) error {
	err := r.db.Where("id = ?", id).First(category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.ErrIngredientCategoryNotFound
		}
		return fmt.Errorf("failed to find ingredient category by ID: %w", err)
	}
	return nil
}

func (r *ingredientCategoryRepository) Create(category *data.IngredientCategory) error {
	return r.db.Create(category).Error
}

func (r *ingredientCategoryRepository) GetByID(id uint) (*data.IngredientCategory, error) {
	var category data.IngredientCategory
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrIngredientCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *ingredientCategoryRepository) GetTranslatedByID(locale data.LanguageCode, id uint) (*data.IngredientCategory, error) {
	var ingredientCategory data.IngredientCategory

	q := r.db.Model(&data.ProductCategory{}).Where("id = ?", id)

	q = utils.ApplyLocalizedPreloads(q, locale, types.IngredientCategoryPreloadMap)

	if err := q.First(&ingredientCategory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrIngredientCategoryNotFound
		}
		return nil, err
	}
	return &ingredientCategory, nil
}

func (r *ingredientCategoryRepository) Update(id uint, updates *data.IngredientCategory) error {
	return r.db.Model(&data.IngredientCategory{}).Where("id = ?", id).Save(updates).Error
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

func (r *ingredientCategoryRepository) GetAll(locale data.LanguageCode, filter *types.IngredientCategoryFilter) ([]data.IngredientCategory, error) {
	var categories []data.IngredientCategory

	q := r.db.Model(&data.IngredientCategory{})

	if filter.Search != nil && *filter.Search != "" {
		pattern := "%" + *filter.Search + "%"
		q = q.Where("name ILIKE ? OR description ILIKE ?", pattern, pattern)
	}

	paged, err := utils.ApplySortedPaginationForModel(
		q,
		filter.Pagination,
		filter.Sort,
		&data.IngredientCategory{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	paged = utils.ApplyLocalizedPreloads(
		paged,
		locale,
		types.IngredientCategoryPreloadMap,
	)

	if err := paged.Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve ingredient categories: %w", err)
	}
	if categories == nil {
		categories = []data.IngredientCategory{}
	}
	return categories, nil
}
