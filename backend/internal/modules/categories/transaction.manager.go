package categories

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/translations"
	"gorm.io/gorm"
)

type TransactionManager interface {
	UpsertProductCategoryTranslations(productCategoryID uint, dto *types.ProductCategoryTranslationsDTO) error
}

type transactionManager struct {
	db                  *gorm.DB
	productCategoryRepo CategoryRepository
	translationsManager translations.TranslationManager
}

func NewTransactionManager(
	db *gorm.DB,
	productCategoryRepo CategoryRepository,
	translationsManager translations.TranslationManager,
) TransactionManager {
	return &transactionManager{
		db:                  db,
		productCategoryRepo: productCategoryRepo,
		translationsManager: translationsManager,
	}
}

func (m *transactionManager) UpsertProductCategoryTranslations(
	categoryID uint, dto *types.ProductCategoryTranslationsDTO,
) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		repoTx := m.productCategoryRepo.CloneWithTransaction(tx)

		var cat data.ProductCategory
		if err := repoTx.FindRawProductCategoryByID(categoryID, &cat); err != nil {
			return fmt.Errorf("load product category: %w", err)
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

		return trx.UpdateProductCategoryTranslationIDs(categoryID, nameID, descID)
	})
}
