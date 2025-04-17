package additives

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/translations"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TransactionManager interface {
	UpsertAdditiveTranslations(additiveID uint, dto *types.AdditiveTranslationsDTO) error
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

func (m *transactionManager) UpsertAdditiveTranslations(additiveID uint, dto *types.AdditiveTranslationsDTO) error {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		repoTx := m.additiveRepo.CloneWithTransaction(tx)

		var additive data.Additive
		if err := repoTx.FindRawAdditiveByID(additiveID, &additive); err != nil {
			return fmt.Errorf("failed to load additive: %w", err)
		}

		translationTx := m.translationsManager.CloneWithTransaction(tx)

		nameGroupID, err := m.upsertFieldTranslations(translationTx, additive.NameTranslationID, dto.Name)
		if err != nil {
			return fmt.Errorf("failed upserting name translations: %w", err)
		}

		descGroupID, err := m.upsertFieldTranslations(translationTx, additive.DescriptionTranslationID, dto.Description)
		if err != nil {
			return fmt.Errorf("failed upserting description translations: %w", err)
		}

		if err := translationTx.UpdateAdditiveTranslationIDs(additiveID, nameGroupID, descGroupID); err != nil {
			return fmt.Errorf("failed to update additive with translation group IDs: %w", err)
		}

		return nil
	})

	return err
}

func (m *transactionManager) upsertFieldTranslations(tx translations.TranslationManager, currentGroupID *uint, locale types.FieldLocale) (uint, error) {
	var entries []struct {
		Language data.LanguageCode
		Text     string
	}
	if locale.En != "" {
		entries = append(entries, struct {
			Language data.LanguageCode
			Text     string
		}{"en", locale.En})
	}
	if locale.Ru != "" {
		entries = append(entries, struct {
			Language data.LanguageCode
			Text     string
		}{"ru", locale.Ru})
	}
	if locale.Kk != "" {
		entries = append(entries, struct {
			Language data.LanguageCode
			Text     string
		}{"kk", locale.Kk})
	}

	// If there are no provided translations, simply return the current group ID (or 0 if none exists).
	if len(entries) == 0 {
		if currentGroupID != nil {
			return *currentGroupID, nil
		}
		return 0, nil
	}

	// If no translation group exists, create one.
	if currentGroupID == nil || *currentGroupID == 0 {
		firstEntry := entries[0]
		firstTranslation := data.AppTranslations{
			TranslationID:  0, // temporary placeholder
			LanguageCode:   firstEntry.Language,
			TranslatedText: firstEntry.Text,
		}
		if err := tx.CreateTranslation(&firstTranslation); err != nil {
			return 0, fmt.Errorf("failed to create first translation: %w", err)
		}
		groupID := firstTranslation.ID
		firstTranslation.TranslationID = groupID
		if err := tx.UpdateTranslation(&firstTranslation); err != nil {
			return 0, fmt.Errorf("failed to update first translation with group id: %w", err)
		}

		for i := 1; i < len(entries); i++ {
			entry := entries[i]
			newRec := data.AppTranslations{
				TranslationID:  groupID,
				LanguageCode:   entry.Language,
				TranslatedText: entry.Text,
			}
			if err := tx.CreateTranslation(&newRec); err != nil {
				return 0, fmt.Errorf("failed creating translation for language %s: %w", entry.Language, err)
			}
		}
		return groupID, nil
	}

	// Otherwise, the translation group already exists.
	groupID := *currentGroupID
	for _, entry := range entries {
		lang := entry.Language
		existing, err := tx.FindTranslation(groupID, lang)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Insert a new translation record if none exists.
				newRec := data.AppTranslations{
					TranslationID:  groupID,
					LanguageCode:   lang,
					TranslatedText: entry.Text,
				}
				if err := tx.CreateTranslation(&newRec); err != nil {
					return 0, fmt.Errorf("failed creating translation for language %s: %w", entry.Language, err)
				}
			} else {
				return 0, fmt.Errorf("failed retrieving translation for language %s: %w", entry.Language, err)
			}
		} else {
			// Update the translation if the text has changed.
			if existing.TranslatedText != entry.Text {
				existing.TranslatedText = entry.Text
				if err := tx.UpdateTranslation(&existing); err != nil {
					return 0, fmt.Errorf("failed updating translation for language %s: %w", entry.Language, err)
				}
			}
		}
	}

	// Delete obsolete translations in this group that are not among our fixed set.
	providedLangs := make([]data.LanguageCode, len(entries))
	for i, e := range entries {
		providedLangs[i] = e.Language
	}
	if err := tx.DeleteObsoleteTranslations(groupID, providedLangs); err != nil {
		return 0, fmt.Errorf("failed to delete obsolete translations: %w", err)
	}

	return groupID, nil
}
