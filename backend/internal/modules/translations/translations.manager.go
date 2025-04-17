package translations

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type TranslationManager interface {
	CreateTranslations(translations []data.AppTranslations) error
	FindTranslation(groupID uint, language data.LanguageCode) (data.AppTranslations, error)
	UpdateTranslation(t *data.AppTranslations) error
	CreateTranslation(t *data.AppTranslations) error
	DeleteObsoleteTranslations(groupID uint, validLangs []data.LanguageCode) error

	UpdateAdditiveTranslationIDs(additiveID, nameGroupID, descGroupID uint) error
	UpdateAdditiveCategoryTranslationIDs(additiveCategoryID, nameGroupID, descGroupID uint) error
	UpdateProductTranslationIDs(productID, nameGroupID, descGroupID uint) error
	UpdateProductCategoryTranslationIDs(productCategoryID, nameGroupID, descGroupID uint) error
	UpdateIngredientTranslationIDs(ingredientID, nameGroupID, descGroupID uint) error
	UpdateIngredientCategoryTranslationIDs(ingredientCategoryID, nameGroupID, descGroupID uint) error

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

func (r *translationManager) UpdateIngredientTranslationIDs(ingredientID, nameGroupID, descGroupID uint) error {
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
