package ingredients

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/translations"
	"gorm.io/gorm"
)

type TransactionManager interface {
	UpsertIngredientTranslations(ingredientID uint, dto *types.IngredientTranslationsDTO) error
}

type transactionManager struct {
	db                  *gorm.DB
	ingredientRepo      IngredientRepository
	translationsManager translations.TranslationManager
}

func NewTransactionManager(
	db *gorm.DB,
	ingredientRepo IngredientRepository,
	translationsManager translations.TranslationManager,

) TransactionManager {
	return &transactionManager{
		db:                  db,
		ingredientRepo:      ingredientRepo,
		translationsManager: translationsManager,
	}
}

func (m *transactionManager) UpsertIngredientTranslations(ingredientID uint, dto *types.IngredientTranslationsDTO) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		repoTx := m.ingredientRepo.CloneWithTransaction(tx)

		var ingredient data.Ingredient
		if err := repoTx.FindRawIngredientByID(ingredientID, &ingredient); err != nil {
			return fmt.Errorf("load ingredient category: %w", err)
		}

		trx := m.translationsManager.CloneWithTransaction(tx)

		nameID, err := trx.UpsertGroup(ingredient.NameTranslationID, dto.Name)
		if err != nil {
			return err
		}

		return trx.UpdateIngredientTranslationIDs(ingredientID, nameID)
	})
}
