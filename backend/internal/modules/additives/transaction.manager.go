package additives

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/translations"
	"gorm.io/gorm"
)

type TransactionManager interface {
	UpsertAdditiveTranslations(additiveID uint, dto *types.AdditiveTranslationsDTO) error
	UpsertAdditiveCategoryTranslations(categoryID uint, dto *types.AdditiveCategoryTranslationsDTO) error
}

type transactionManager struct {
	db                  *gorm.DB
	additiveRepo        AdditiveRepository
	translationsManager translations.TranslationManager
}

func NewTransactionManager(
	db *gorm.DB,
	additiveRepo AdditiveRepository,
	translationsManager translations.TranslationManager,
) TransactionManager {
	return &transactionManager{
		db:                  db,
		additiveRepo:        additiveRepo,
		translationsManager: translationsManager,
	}
}

func (m *transactionManager) UpsertAdditiveTranslations(
	additiveID uint, dto *types.AdditiveTranslationsDTO,
) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		repoTx := m.additiveRepo.CloneWithTransaction(tx)

		var add data.Additive
		if err := repoTx.FindRawAdditiveByID(additiveID, &add); err != nil {
			return fmt.Errorf("load additive: %w", err)
		}

		trx := m.translationsManager.CloneWithTransaction(tx)

		nameID, err := trx.UpsertGroup(add.NameTranslationID, dto.Name)
		if err != nil {
			return err
		}

		descID, err := trx.UpsertGroup(add.DescriptionTranslationID, dto.Description)
		if err != nil {
			return err
		}

		return trx.UpdateAdditiveTranslationIDs(additiveID, nameID, descID)
	})
}

func (m *transactionManager) UpsertAdditiveCategoryTranslations(
	categoryID uint, dto *types.AdditiveCategoryTranslationsDTO,
) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		repoTx := m.additiveRepo.CloneWithTransaction(tx)

		var cat data.AdditiveCategory
		if err := repoTx.FindRawAdditiveCategoryByID(categoryID, &cat); err != nil {
			return fmt.Errorf("load category: %w", err)
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

		return trx.UpdateAdditiveCategoryTranslationIDs(categoryID, nameID, descID)
	})
}
