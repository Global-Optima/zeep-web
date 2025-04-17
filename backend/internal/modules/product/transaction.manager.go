package product

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/translations"
	"gorm.io/gorm"
)

type TransactionManager interface {
	UpsertProductTranslations(productID uint, dto *types.ProductTranslationsDTO) error
}

type transactionManager struct {
	db                  *gorm.DB
	productRepo         ProductRepository
	translationsManager translations.TranslationManager
}

func NewTransactionManager(
	db *gorm.DB,
	productRepo ProductRepository,
	translationsManager translations.TranslationManager,

) TransactionManager {
	return &transactionManager{
		db:                  db,
		productRepo:         productRepo,
		translationsManager: translationsManager,
	}
}

func (m *transactionManager) UpsertProductTranslations(
	productID uint, dto *types.ProductTranslationsDTO) error {

	return m.db.Transaction(func(tx *gorm.DB) error {
		repoTx := m.productRepo.CloneWithTransaction(tx)

		var prod data.Product
		if err := repoTx.FindRawProductByID(productID, &prod); err != nil {
			return fmt.Errorf("load product: %w", err)
		}

		trx := m.translationsManager.CloneWithTransaction(tx)

		nameID, err := trx.UpsertGroup(prod.NameTranslationID, dto.Name)
		if err != nil {
			return err
		}

		descID, err := trx.UpsertGroup(prod.DescriptionTranslationID, dto.Description)
		if err != nil {
			return err
		}

		return trx.UpdateProductTranslationIDs(productID, nameID, descID)
	})
}
