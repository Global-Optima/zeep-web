package translations

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type TranslationManager interface {
	UpsertGroup(current *uint, loc FieldLocale) (uint, error)

	CreateTranslations(translations []data.AppTranslations) error
	FindTranslation(groupID uint, language data.LanguageCode) (data.AppTranslations, error)
	UpdateTranslation(t *data.AppTranslations) error
	CreateTranslation(t *data.AppTranslations) error
	DeleteObsoleteTranslations(groupID uint, validLangs []data.LanguageCode) error

	UpdateAdditiveTranslationIDs(additiveID, nameGroupID, descGroupID uint) error
	UpdateAdditiveCategoryTranslationIDs(additiveCategoryID, nameGroupID, descGroupID uint) error
	UpdateProductTranslationIDs(productID, nameGroupID, descGroupID uint) error
	UpdateProductCategoryTranslationIDs(productCategoryID, nameGroupID, descGroupID uint) error
	UpdateIngredientTranslationIDs(ingredientID, nameGroupID uint) error
	UpdateIngredientCategoryTranslationIDs(ingredientCategoryID, nameGroupID, descGroupID uint) error
	UpdateUnitTranslationIDs(unitID, nameGroupID uint) error

	CloneWithTransaction(tx *gorm.DB) TranslationManager
}

type translationManager struct {
	db *gorm.DB
}

func NewTranslationManager(db *gorm.DB) TranslationManager {
	return &translationManager{
		db: db,
	}
}

func (r *translationManager) CloneWithTransaction(tx *gorm.DB) TranslationManager {
	return &translationManager{
		db: tx,
	}
}

func (m *translationManager) UpsertGroup(
	current *uint,
	loc FieldLocale,
) (uint, error) {
	entries := m.buildEntries(loc)
	if len(entries) == 0 { // nothing to do
		if current != nil {
			return *current, nil
		}
		return 0, nil
	}

	if current == nil || *current == 0 { // no group yet
		return m.createGroup(entries)
	}
	return m.updateGroup(*current, entries)
}

type langEntry struct {
	lang data.LanguageCode
	text string
}

func (m *translationManager) buildEntries(loc FieldLocale) []langEntry {
	entries := make([]langEntry, 0, 3)
	if loc.En != "" {
		entries = append(entries, langEntry{data.LanguageCodeEN, loc.En})
	}
	if loc.Ru != "" {
		entries = append(entries, langEntry{data.LanguageCodeRU, loc.Ru})
	}
	if loc.Kk != "" {
		entries = append(entries, langEntry{data.LanguageCodeKK, loc.Kk})
	}
	return entries
}

func (m *translationManager) createGroup(entries []langEntry) (uint, error) {
	first := entries[0]
	rec := data.AppTranslations{
		LanguageCode:   first.lang,
		TranslatedText: first.text,
	}
	if err := m.CreateTranslation(&rec); err != nil {
		return 0, fmt.Errorf("create first translation: %w", err)
	}
	groupID := rec.ID

	// back‑patch its own group‑id and insert the rest
	rec.TranslationID = groupID
	if err := m.UpdateTranslation(&rec); err != nil {
		return 0, fmt.Errorf("patch group id: %w", err)
	}
	for _, e := range entries[1:] {
		if err := m.CreateTranslation(&data.AppTranslations{
			TranslationID:  groupID,
			LanguageCode:   e.lang,
			TranslatedText: e.text,
		}); err != nil {
			return 0, fmt.Errorf("insert %s: %w", e.lang, err)
		}
	}
	return groupID, nil
}

func (m *translationManager) updateGroup(
	groupID uint, entries []langEntry,
) (uint, error) {
	for _, e := range entries {
		if err := m.upsertSingle(groupID, e); err != nil {
			return 0, err
		}
	}

	valid := make([]data.LanguageCode, len(entries))
	for i, e := range entries {
		valid[i] = e.lang
	}

	if err := m.DeleteObsoleteTranslations(groupID, valid); err != nil {
		return 0, fmt.Errorf("purge obsolete: %w", err)
	}
	return groupID, nil
}

func (m *translationManager) upsertSingle(
	groupID uint, e langEntry,
) error {
	ex, err := m.FindTranslation(groupID, e.lang)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("read %s: %w", e.lang, err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return m.CreateTranslation(&data.AppTranslations{
			TranslationID:  groupID,
			LanguageCode:   e.lang,
			TranslatedText: e.text,
		})
	}

	if ex.TranslatedText != e.text {
		ex.TranslatedText = e.text
		return m.UpdateTranslation(&ex)
	}
	return nil
}

func (r *translationManager) CreateTranslations(translations []data.AppTranslations) error {
	return r.db.Create(&translations).Error
}

func (r *translationManager) FindTranslation(groupID uint, language data.LanguageCode) (data.AppTranslations, error) {
	var t data.AppTranslations
	err := r.db.
		Where("translation_id = ? AND language_code = ?", groupID, language).
		First(&t).Error
	return t, err
}

func (r *translationManager) UpdateTranslation(t *data.AppTranslations) error {
	return r.db.Model(t).Updates(t).Error
}

func (r *translationManager) CreateTranslation(t *data.AppTranslations) error {
	return r.db.Create(t).Error
}

func (r *translationManager) DeleteObsoleteTranslations(groupID uint, validLangs []data.LanguageCode) error {
	q := r.db.Unscoped().Where("translation_id = ?", groupID)

	if len(validLangs) > 0 {
		q = q.Where("language_code NOT IN ?", validLangs)
	}
	return q.Delete(&data.AppTranslations{}).Error
}

func (r *translationManager) UpdateAdditiveTranslationIDs(additiveID, nameGroupID, descGroupID uint) error {
	updates := &data.Additive{
		NameTranslationID:        &nameGroupID,
		DescriptionTranslationID: &descGroupID,
	}
	if err := r.db.Model(&data.Additive{}).Where("id = ?", additiveID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update additive translation IDs: %w", err)
	}
	return nil
}

func (r *translationManager) UpdateAdditiveCategoryTranslationIDs(additiveCategoryID, nameGroupID, descGroupID uint) error {
	updates := &data.AdditiveCategory{
		NameTranslationID:        &nameGroupID,
		DescriptionTranslationID: &descGroupID,
	}
	if err := r.db.Model(&data.AdditiveCategory{}).Where("id = ?", additiveCategoryID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update additive category translation IDs: %w", err)
	}
	return nil
}

func (r *translationManager) UpdateProductTranslationIDs(productID, nameGroupID, descGroupID uint) error {
	updates := &data.Product{
		NameTranslationID:        &nameGroupID,
		DescriptionTranslationID: &descGroupID,
	}
	if err := r.db.Model(&data.Product{}).Where("id = ?", productID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update product translation IDs: %w", err)
	}
	return nil
}

func (r *translationManager) UpdateProductCategoryTranslationIDs(productCategoryID, nameGroupID, descGroupID uint) error {
	updates := &data.ProductCategory{
		NameTranslationID:        &nameGroupID,
		DescriptionTranslationID: &descGroupID,
	}
	if err := r.db.Model(&data.ProductCategory{}).Where("id = ?", productCategoryID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update product category translation IDs: %w", err)
	}
	return nil
}

func (r *translationManager) UpdateIngredientTranslationIDs(ingredientID, nameGroupID uint) error {
	updates := &data.Ingredient{
		NameTranslationID: &nameGroupID,
	}
	if err := r.db.Model(&data.Ingredient{}).Where("id = ?", ingredientID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update ingredient translation IDs: %w", err)
	}
	return nil
}

func (r *translationManager) UpdateIngredientCategoryTranslationIDs(ingredientCategoryID, nameGroupID, descGroupID uint) error {
	updates := &data.IngredientCategory{
		NameTranslationID:        &nameGroupID,
		DescriptionTranslationID: &descGroupID,
	}
	if err := r.db.Model(&data.IngredientCategory{}).Where("id = ?", ingredientCategoryID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update ingredient category translation IDs: %w", err)
	}
	return nil
}

func (r *translationManager) UpdateUnitTranslationIDs(unitID, nameGroupID uint) error {
	updates := &data.Unit{
		NameTranslationID: &nameGroupID,
	}
	if err := r.db.Model(&data.Unit{}).Where("id = ?", unitID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update unit translation IDs: %w", err)
	}
	return nil
}
