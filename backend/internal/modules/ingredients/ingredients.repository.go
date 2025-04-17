package ingredients

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type IngredientRepository interface {
	CreateIngredient(ingredient *data.Ingredient) (uint, error)
	SaveIngredient(ingredientID uint, ingredient *data.Ingredient) error
	DeleteIngredient(ingredientID uint) error
	GetIngredientByID(ingredientID uint) (*data.Ingredient, error)
	GetTranslatedIngredientByID(locale data.LanguageCode, ingredientID uint) (*data.Ingredient, error)
	GetRawIngredientByID(ingredientID uint) (*data.Ingredient, error)
	GetIngredientsWithDetailsByIDs(ingredientIDs []uint) ([]data.Ingredient, error)
	GetIngredients(locale data.LanguageCode, filter *types.IngredientFilter) ([]data.Ingredient, error)
	GetIngredientsForProductSizes(productSizeIDs []uint) ([]data.Ingredient, error)
	GetIngredientsForAdditives(additiveIDs []uint) ([]data.Ingredient, error)
	GetIngredientsForProvisions(provisionIDs []uint) ([]data.Ingredient, error)
	GetIngredientIDsForProvisions(provisionIDs []uint) ([]uint, error)
	FindRawIngredientByID(ingredientID uint, ingredient *data.Ingredient) error

	CloneWithTransaction(tx *gorm.DB) IngredientRepository
}

type ingredientRepository struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) IngredientRepository {
	return &ingredientRepository{db: db}
}

func (r *ingredientRepository) CloneWithTransaction(tx *gorm.DB) IngredientRepository {
	return &ingredientRepository{
		db: tx,
	}
}

func (r *ingredientRepository) FindRawIngredientByID(ingredientID uint, ingredient *data.Ingredient) error {
	return r.db.
		Where("id = ?", ingredientID).
		First(ingredient).Error
}

func (r *ingredientRepository) CreateIngredient(ingredient *data.Ingredient) (uint, error) {
	err := r.db.Create(ingredient).Error
	if err != nil {
		return 0, err
	}
	return ingredient.ID, err
}

func (r *ingredientRepository) SaveIngredient(ingredientID uint, ingredient *data.Ingredient) error {
	return r.db.Model(&data.Ingredient{}).
		Where("id = ?", ingredientID).
		Save(ingredient).Error
}

func (r *ingredientRepository) GetRawIngredientByID(ingredientID uint) (*data.Ingredient, error) {
	var ingredient data.Ingredient
	err := r.db.First(&ingredient, ingredientID).Error
	if err != nil {
		return nil, err
	}

	return &ingredient, nil
}

func (r *ingredientRepository) GetIngredientByID(ingredientID uint) (*data.Ingredient, error) {
	var ingredient data.Ingredient
	err := r.db.Preload("Unit").
		Preload("IngredientCategory").
		First(&ingredient, ingredientID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrIngredientNotFound
		}
		return nil, err
	}

	return &ingredient, nil
}

func (r *ingredientRepository) GetTranslatedIngredientByID(locale data.LanguageCode, ingredientID uint) (*data.Ingredient, error) {
	var ingredient data.Ingredient

	q := r.db.Model(&data.Ingredient{}).Where("id = ?", ingredientID)

	q = utils.ApplyLocalizedPreloads(q, locale, types.IngredientPreloadMap)

	if err := q.First(&ingredient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrIngredientNotFound
		}
		return nil, err
	}
	return &ingredient, nil
}

func (r *ingredientRepository) GetIngredientsWithDetailsByIDs(ingredientIDs []uint) ([]data.Ingredient, error) {
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

func (r *ingredientRepository) GetIngredients(locale data.LanguageCode, filter *types.IngredientFilter) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient

	q := r.db.Model(&data.Ingredient{})

	if filter.ProductSizeID != nil {
		q = q.Joins(`
			JOIN product_size_ingredients psi
			  ON psi.ingredient_id = ingredients.id
		`).Where("psi.product_size_id = ?", *filter.ProductSizeID)
	}

	if filter.Name != nil && *filter.Name != "" {
		q = q.Where("ingredients.name ILIKE ?", "%"+*filter.Name+"%")
	}
	if filter.MinCalories != nil {
		q = q.Where("calories >= ?", *filter.MinCalories)
	}
	if filter.MaxCalories != nil {
		q = q.Where("calories <= ?", *filter.MaxCalories)
	}
	if filter.IsAllergen != nil {
		q = q.Where("is_allergen = ?", *filter.IsAllergen)
	}

	paged, err := utils.ApplySortedPaginationForModel(
		q, filter.Pagination, filter.Sort, &data.Ingredient{},
	)
	if err != nil {
		return nil, err
	}

	paged = utils.ApplyLocalizedPreloads(
		paged,
		locale,
		types.IngredientPreloadMap,
	)

	if err := paged.Find(&ingredients).Error; err != nil {
		return nil, err
	}
	if ingredients == nil {
		ingredients = []data.Ingredient{}
	}
	return ingredients, nil
}

func (r *ingredientRepository) GetIngredientsForProductSizes(productSizeIDs []uint) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient

	subDirect := r.db.Model(&data.ProductSizeIngredient{}).
		Select("ingredient_id").
		Where("product_size_id IN ?", productSizeIDs)

	subAdditive := r.db.Model(&data.ProductSizeAdditive{}).
		Select("ai.ingredient_id").
		Joins("JOIN additive_ingredients ai ON ai.additive_id = product_size_additives.additive_id").
		Where("product_size_additives.product_size_id IN ?", productSizeIDs)

	subProvision := r.db.Model(&data.ProductSizeProvision{}).
		Select("pi.ingredient_id").
		Joins("JOIN provision_ingredients pi ON pi.provision_id = product_size_provisions.provision_id").
		Where("product_size_provisions.product_size_id IN ?", productSizeIDs)

	query := r.db.Model(&data.Ingredient{}).
		Select("DISTINCT ingredients.*").
		Preload("Unit").
		Preload("IngredientCategory").
		Where("ingredients.id IN (?) OR ingredients.id IN (?) OR ingredients.id IN (?)", subDirect, subAdditive, subProvision)

	if err := query.Find(&ingredients).Error; err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (r *ingredientRepository) GetIngredientsForAdditives(additiveIDs []uint) ([]data.Ingredient, error) {
	ingredientIDsDirect, err := r.getIngredientIDsFromAdditiveIngredients(additiveIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ingredients from additive_ingredients: %w", err)
	}

	ingredientIDsFromProvisions, err := r.getIngredientIDsFromAdditiveProvisions(additiveIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ingredients from additive_provisions: %w", err)
	}

	return r.fetchIngredientsByIDs(utils.UnionSlices(ingredientIDsDirect, ingredientIDsFromProvisions))
}

func (r *ingredientRepository) getIngredientIDsFromAdditiveIngredients(additiveIDs []uint) ([]uint, error) {
	var ids []uint
	err := r.db.
		Model(&data.AdditiveIngredient{}).
		Distinct("ingredient_id").
		Where("additive_id IN ?", additiveIDs).
		Pluck("ingredient_id", &ids).Error
	return ids, err
}

func (r *ingredientRepository) getIngredientIDsFromAdditiveProvisions(additiveIDs []uint) ([]uint, error) {
	var ids []uint
	err := r.db.
		Model(&data.ProvisionIngredient{}).
		Distinct("provision_ingredients.ingredient_id").
		Joins("JOIN additive_provisions ON additive_provisions.provision_id = provision_ingredients.provision_id").
		Where("additive_provisions.additive_id IN ?", additiveIDs).
		Pluck("provision_ingredients.ingredient_id", &ids).Error
	return ids, err
}

func (r *ingredientRepository) fetchIngredientsByIDs(ingredientIDs []uint) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient
	if len(ingredientIDs) == 0 {
		return ingredients, nil
	}

	err := r.db.
		Model(&data.Ingredient{}).
		Preload("Unit").
		Preload("IngredientCategory").
		Where("id IN ?", ingredientIDs).
		Find(&ingredients).Error

	return ingredients, err
}

func (r *ingredientRepository) GetIngredientsForProvisions(provisionIDs []uint) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient

	if len(provisionIDs) == 0 {
		return ingredients, nil
	}

	query := r.db.Model(&data.Ingredient{}).
		Select("DISTINCT ingredients.*").
		Preload("Unit").
		Preload("IngredientCategory").
		Joins("JOIN provision_ingredients pi ON pi.ingredient_id = ingredients.id").
		Where("pi.provision_id IN ?", provisionIDs)

	if err := query.Find(&ingredients).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch ingredients for provisions: %w", err)
	}

	return ingredients, nil
}

func (r *ingredientRepository) GetIngredientIDsForProvisions(provisionIDs []uint) ([]uint, error) {
	var ids []uint

	if len(provisionIDs) == 0 {
		return ids, nil
	}

	query := r.db.Model(&data.ProvisionIngredient{}).
		Distinct("ingredient_id").
		Where("provision_id IN ?", provisionIDs)

	if err := query.Find(&ids).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch ingredientIDs for provisions: %w", err)
	}

	return ids, nil
}

func (r *ingredientRepository) DeleteIngredient(ingredientID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := checkIngredientReferences(tx, ingredientID); err != nil {
			return err
		}

		if err := tx.Unscoped().Delete(&data.Ingredient{}, ingredientID).Error; err != nil {
			return err
		}

		return nil
	})
}

func checkIngredientReferences(db *gorm.DB, ingredientID uint) error {
	var ingredient data.Ingredient

	err := db.
		Preload("StockMaterials", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Preload("ProductSizeIngredients", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Preload("StoreStocks", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Preload("AdditiveIngredients", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Where(&data.Ingredient{BaseEntity: data.BaseEntity{ID: ingredientID}}).
		First(&ingredient).Error
	if err != nil {
		return err
	}

	if len(ingredient.StockMaterials) > 0 || len(ingredient.ProductSizeIngredients) > 0 ||
		len(ingredient.StoreStocks) > 0 || len(ingredient.AdditiveIngredients) > 0 {
		return types.ErrIngredientIsInUse
	}

	return nil
}
