package units

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/translations"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"gorm.io/gorm"
)

type TransactionManager interface {
	UpsertUnitTranslations(unitID uint, dto *types.UnitTranslationsDTO) error
}

type transactionManager struct {
	db                  *gorm.DB
	unitRepo            UnitRepository
	translationsManager translations.TranslationManager
}

func NewTransactionManager(
	db *gorm.DB,
	unitRepo UnitRepository,
	translationsManager translations.TranslationManager,
) TransactionManager {
	return &transactionManager{
		db:                  db,
		unitRepo:            unitRepo,
		translationsManager: translationsManager,
	}
}

func (m *transactionManager) UpsertUnitTranslations(unitID uint, dto *types.UnitTranslationsDTO) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		repoTx := m.unitRepo.CloneWithTransaction(tx)

		var unit data.Unit
		if err := repoTx.FindRawUnitByID(unitID, &unit); err != nil {
			return fmt.Errorf("load unit category: %w", err)
		}

		trx := m.translationsManager.CloneWithTransaction(tx)

		nameID, err := trx.UpsertGroup(unit.NameTranslationID, dto.Name)
		if err != nil {
			return err
		}

		return trx.UpdateUnitTranslationIDs(unitID, nameID)
	})
}
