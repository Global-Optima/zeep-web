package ingredientCategories

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories/types"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/translations"
	"gorm.io/gorm"
)

type TransactionManager interface {
	UpsertIngredientCategoryTranslations(ingredientCategoryID uint, dto *types.IngredientCategoryTranslationDTO) error
}

type transactionManager struct {
	db                     *gorm.DB
	ingredientCategoryRepo IngredientCategoryRepository
	translationsManager    translations.TranslationManager
}

func NewTransactionManager(
	db *gorm.DB,
	ingredientCategoryRepo IngredientCategoryRepository,
	translationsManager translations.TranslationManager,

) TransactionManager {
	return &transactionManager{
		db:                     db,
		ingredientCategoryRepo: ingredientCategoryRepo,
		translationsManager:    translationsManager,
	}
}

func (m *transactionManager) UpsertIngredientCategoryTranslations(categoryID uint, dto *types.IngredientCategoryTranslationDTO) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		repoTx := m.ingredientCategoryRepo.CloneWithTransaction(tx)

		var cat data.IngredientCategory
		if err := repoTx.FindRawIngredientCategoryByID(categoryID, &cat); err != nil {
			return fmt.Errorf("load ingredient category: %w", err)
		}

		trx := m.translationsManager.CloneWithTransaction(tx)

		nameID, err := trx.UpsertGroup(cat.NameTranslationID, dto.Name)
		if err != nil {
			return err
		}

		descID, err := trx.UpsertGroup(cat.DescriptionTranslationID, dto.Description)
		if err != nil {
			return err
		}

		return trx.UpdateIngredientCategoryTranslationIDs(categoryID, nameID, descID)
	})
}
