package additives

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type AdditiveRepository interface {
	GetAdditivesByProductSizeIDs(productSizeIDs []uint) ([]data.ProductSizeAdditive, error)
	CheckAdditiveExists(additiveName string) (bool, error)
	GetAdditiveByID(additiveID uint) (*data.Additive, error)
	GetAdditivesByIDs(additiveIDs []uint) ([]data.Additive, error)
	GetAdditives(filter *types.AdditiveFilterQuery) ([]data.Additive, error)
	CreateAdditive(additive *data.Additive) (uint, error)
	UpdateAdditiveWithAssociations(additiveID uint, updateModels *types.AdditiveModels) error
	DeleteAdditive(additiveID uint) (*data.Additive, error)

	GetAdditiveCategories(filter *types.AdditiveCategoriesFilterQuery) ([]data.AdditiveCategory, error)
	CreateAdditiveCategory(category *data.AdditiveCategory) (uint, error)
	UpdateAdditiveCategory(category *data.AdditiveCategory) error
	DeleteAdditiveCategory(categoryID uint) error
	GetAdditiveCategoryByID(categoryID uint) (*data.AdditiveCategory, error)
}

type additiveRepository struct {
	db *gorm.DB
}

func NewAdditiveRepository(db *gorm.DB) AdditiveRepository {
	return &additiveRepository{db: db}
}

func (r *additiveRepository) GetAdditivesByProductSizeIDs(productSizeIDs []uint) ([]data.ProductSizeAdditive, error) {
	var additives []data.ProductSizeAdditive

	if len(productSizeIDs) == 0 {
		return additives, nil
	}

	err := r.db.Model(&data.ProductSizeAdditive{}).
		Where("product_size_id IN (?)", productSizeIDs).
		Find(&additives).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return additives, moduleErrors.ErrNotFound
		}
		return additives, err
	}

	return additives, nil
}

func (r *additiveRepository) CheckAdditiveExists(additiveName string) (bool, error) {
	var addtive data.Additive

	err := r.db.Model(&data.Additive{}).
		Where(&data.Additive{Name: additiveName}).
		First(&addtive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *additiveRepository) GetAdditiveCategories(filter *types.AdditiveCategoriesFilterQuery) ([]data.AdditiveCategory, error) {
	var categories []data.AdditiveCategory

	query := r.db.Model(&data.AdditiveCategory{}).
		Preload("Additives").
		Preload("Additives.Unit")

	hasAdditivesCondition := "EXISTS (SELECT 1 FROM additives WHERE additives.additive_category_id = additive_categories.id)"

	query = query.Joins("LEFT JOIN additives ON additives.additive_category_id = additive_categories.id").
		Joins("LEFT JOIN store_additives ON store_additives.additive_id = additives.id")

	if filter.ProductSizeId != nil {
		query = query.Where(
			"EXISTS (SELECT 1 FROM product_size_additives WHERE product_size_additives.additive_id = additives.id AND product_size_additives.product_size_id = ?)",
			*filter.ProductSizeId,
		)
	}

	if filter.IsMultipleSelect != nil {
		query = query.Where("additives.is_multiple_select = ?", *filter.IsMultipleSelect)
	}

	if filter.IncludeEmpty != nil && *filter.IncludeEmpty {
		query = query.Where(hasAdditivesCondition + " OR NOT " + hasAdditivesCondition)
	} else {
		query = query.Where(hasAdditivesCondition)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"additive_categories.name ILIKE ? OR additive_categories.description ILIKE ?",
			searchTerm, searchTerm,
		)
	}

	query = query.Group("additive_categories.id")

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.AdditiveCategory{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *additiveRepository) GetAdditives(filter *types.AdditiveFilterQuery) ([]data.Additive, error) {
	var additives []data.Additive

	query := r.db.
		Preload("Category").
		Preload("Unit").
		Joins("JOIN additive_categories ON additives.additive_category_id = additive_categories.id")

	var err error

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("additives.name ILIKE ? OR additives.description ILIKE ?", searchTerm, searchTerm)
	}

	if filter.MinPrice != nil {
		query = query.Where("additives.base_price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("additives.base_price <= ?", *filter.MaxPrice)
	}

	if filter.CategoryID != nil {
		query = query.Where("additive_categories.id = ?", *filter.CategoryID)
	}
	if filter.ProductSizeID != nil {
		query = query.Where("product_additives.product_size_id = ?", *filter.ProductSizeID)
	}

	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Additive{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&additives).Error; err != nil {
		return nil, err
	}

	return additives, nil
}

func (r *additiveRepository) GetAdditiveByID(additiveID uint) (*data.Additive, error) {
	var additive data.Additive
	err := r.db.Model(&data.Additive{}).
		Preload("Category").
		Where("id = ?", additiveID).
		Preload("Unit").
		Preload("Ingredients.Ingredient.Unit").
		Preload("Ingredients.Ingredient.IngredientCategory").
		First(&additive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, moduleErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch additive with ID %d: %w", additiveID, err)
	}

	return &additive, nil
}

func (r *additiveRepository) GetAdditivesByIDs(additiveIDs []uint) ([]data.Additive, error) {
	var additives []data.Additive

	query := r.db.
		Preload("Category").
		Preload("Unit").
		Where("id IN (?)", additiveIDs)

	if err := query.Find(&additives).Error; err != nil {
		return nil, err
	}

	return additives, nil
}

func (r *additiveRepository) CreateAdditive(additive *data.Additive) (uint, error) {
	err := r.db.Create(additive).Error
	if err != nil {
		return 0, err
	}
	return additive.ID, nil
}

func (r *additiveRepository) UpdateAdditiveWithAssociations(additiveID uint, updateModels *types.AdditiveModels) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels != nil {
			if err := tx.Model(&data.Additive{}).
				Where("id = ?", additiveID).
				Updates(updateModels.Additive).Error; err != nil {
				return fmt.Errorf("failed to update additive: %w", err)
			}
		}

		if updateModels.Ingredients != nil {
			// Remove existing ingredients
			if err := tx.Where("additive_id = ?", additiveID).Delete(&data.AdditiveIngredient{}).Error; err != nil {
				return fmt.Errorf("failed to delete ingredients: %w", err)
			}

			// Add new ingredients
			ingredients := make([]data.AdditiveIngredient, len(updateModels.Ingredients))
			for i, ingredient := range updateModels.Ingredients {
				ingredients[i] = data.AdditiveIngredient{
					AdditiveID:   additiveID,
					IngredientID: ingredient.IngredientID,
					Quantity:     ingredient.Quantity,
				}
			}
			if err := tx.Create(ingredients).Error; err != nil {
				return fmt.Errorf("failed to create ingredients: %w", err)
			}
		}

		return nil
	})
}

func (r *additiveRepository) DeleteAdditive(additiveID uint) (*data.Additive, error) {
	var additive data.Additive
	if err := r.db.First(&additive, additiveID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrAdditiveNotFound
		}
		return nil, err
	}

	if err := r.db.Where("id = ?", additiveID).Delete(&data.Additive{}).Error; err != nil {
		return nil, err
	}

	return &additive, nil
}

func (r *additiveRepository) CreateAdditiveCategory(category *data.AdditiveCategory) (uint, error) {
	err := r.db.Create(category).Error
	if err != nil {
		return 0, err
	}
	return category.ID, nil
}

func (r *additiveRepository) UpdateAdditiveCategory(category *data.AdditiveCategory) error {
	return r.db.Save(category).Error
}

func (r *additiveRepository) DeleteAdditiveCategory(categoryID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := checkAdditiveCategoryReferences(tx, categoryID); err != nil {
			return err
		}

		if err := tx.Delete(&data.AdditiveCategory{}, categoryID).Error; err != nil {
			return err
		}

		return nil
	})
}

func checkAdditiveCategoryReferences(db *gorm.DB, additiveCategoryID uint) error {
	var additiveCategory data.AdditiveCategory

	err := db.
		Preload("Additives", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Where(&data.AdditiveCategory{BaseEntity: data.BaseEntity{ID: additiveCategoryID}}).
		First(&additiveCategory).Error
	if err != nil {
		return err
	}

	if len(additiveCategory.Additives) > 0 {
		return types.ErrAdditiveCategoryIsInUse
	}

	return nil
}

func (r *additiveRepository) GetAdditiveCategoryByID(categoryID uint) (*data.AdditiveCategory, error) {
	var category data.AdditiveCategory
	err := r.db.Preload("Additives").First(&category, categoryID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, moduleErrors.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}
