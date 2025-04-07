package provisions

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ProvisionRepository interface {
	CheckProvisionExists(provisionName string) (bool, error)
	CreateProvision(provision *data.Provision) (uint, error)
	GetProvisions(filter *types.ProvisionFilterDTO) ([]data.Provision, error)
	GetProvisionByID(provisionID uint) (*data.Provision, error)
	GetProvisionWithDetailsByID(provisionID uint) (*data.Provision, error)
	GetProvisionIngredientIDs(provisionID uint) ([]uint, error)
	SaveProvisionWithAssociations(updateModels *types.ProvisionModels) error
	DeleteProvision(provisionID uint) (*data.Provision, error)
}

type provisionRepository struct {
	db *gorm.DB
}

func NewProvisionRepository(db *gorm.DB) ProvisionRepository {
	return &provisionRepository{db: db}
}

func (r *provisionRepository) CheckProvisionExists(provisionName string) (bool, error) {
	var provision data.Provision
	err := r.db.Where("name = ?", provisionName).First(&provision).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *provisionRepository) CreateProvision(provision *data.Provision) (uint, error) {
	if err := r.db.Create(provision).Error; err != nil {
		if strings.Contains(err.Error(), "23505") { // unique constraint
			return 0, types.ErrProvisionUniqueName
		}
		return 0, err
	}
	return provision.ID, nil
}

func (r *provisionRepository) GetProvisions(filter *types.ProvisionFilterDTO) ([]data.Provision, error) {
	var provisions []data.Provision

	query := r.db.Model(&data.Provision{}).
		Preload("Unit")

	if filter.Search != nil {
		pattern := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ?", pattern)
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Provision{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	err = query.Find(&provisions).Error
	if err != nil {
		return nil, err
	}

	if provisions == nil {
		provisions = []data.Provision{}
	}

	return provisions, nil
}

func (r *provisionRepository) GetProvisionByID(provisionID uint) (*data.Provision, error) {
	var provision data.Provision
	err := r.db.Where("id = ?", provisionID).First(&provision).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrProvisionNotFound
		}
		return nil, err
	}
	return &provision, nil
}

func (r *provisionRepository) GetProvisionWithDetailsByID(provisionID uint) (*data.Provision, error) {
	var provision data.Provision
	err := r.db.
		Preload("Unit").
		Preload("ProvisionIngredients.Ingredient.Unit").
		Preload("ProvisionIngredients.Ingredient.IngredientCategory").
		Where("id = ?", provisionID).
		First(&provision).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrProvisionNotFound
		}
		return nil, err
	}
	return &provision, nil
}

func (r *provisionRepository) GetProvisionIngredientIDs(provisionID uint) ([]uint, error) {
	var provisionIngredients []uint

	err := r.db.Model(&data.ProvisionIngredient{}).
		Distinct("ingredient_id").
		Where("provision_id = ?", provisionID).
		Pluck("ingredient_id", &provisionIngredients).Error
	if err != nil {
		return nil, err
	}

	return provisionIngredients, nil
}

func (r *provisionRepository) SaveProvisionWithAssociations(updateModels *types.ProvisionModels) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels == nil || updateModels.Provision == nil {
			return fmt.Errorf("nothing to update")
		}

		ingredientsUpdated := false

		if updateModels.Ingredients != nil {
			if err := r.updateProvisionIngredients(tx, updateModels.Provision.ID, updateModels.Ingredients); err != nil {
				return err
			}
			ingredientsUpdated = true
		}

		if err := r.saveProvision(tx, updateModels.Provision, ingredientsUpdated); err != nil {
			return err
		}

		productSizeIDs, err := r.getProductSizeIDsByProvision(tx, updateModels.Provision.ID)
		if err != nil {
			return err
		}

		additiveIDs, err := r.getProductSizeIDsByProvision(tx, updateModels.Provision.ID)
		if err != nil {
			return err
		}

		if err := r.updateProductSizeProvisionsUpdatedAt(tx, productSizeIDs); err != nil {
			return err
		}

		if err := r.updateAdditiveProvisionsUpdatedAt(tx, additiveIDs); err != nil {
			return err
		}

		return nil
	})
}

func (r *provisionRepository) getProductSizeIDsByProvision(tx *gorm.DB, provisionID uint) ([]uint, error) {
	var productSizeIDs []uint

	err := tx.Model(&data.ProductSizeProvision{}).
		Where("provision_id = ?", provisionID).
		Pluck("product_size_id", &productSizeIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get productSizeIDs by proviosionID %d: %w", provisionID, err)
	}

	return productSizeIDs, nil
}

func (r *provisionRepository) getAdditivesIDsByProvision(tx *gorm.DB, provisionID uint) ([]uint, error) {
	var additiveIDs []uint

	err := tx.Model(&data.AdditiveProvision{}).
		Where("provision_id = ?", provisionID).
		Pluck("additive_id", &additiveIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get additiveIDs by proviosionID %d: %w", provisionID, err)
	}

	return additiveIDs, nil
}

func (r *provisionRepository) updateProductSizeProvisionsUpdatedAt(tx *gorm.DB, productSizeIDs []uint) error {
	if len(productSizeIDs) == 0 {
		return nil
	}

	err := tx.Model(&data.ProductSize{}).
		Where("id IN ?", productSizeIDs).
		Update("provisions_updated_at", time.Now().UTC()).Error
	if err != nil {
		return fmt.Errorf("failed to update provisions_updated_at for productSizeIDs: %w", err)
	}

	return nil
}

func (r *provisionRepository) updateAdditiveProvisionsUpdatedAt(tx *gorm.DB, additiveIDs []uint) error {
	if len(additiveIDs) == 0 {
		return nil
	}

	err := tx.Model(&data.Additive{}).
		Where("id IN ?", additiveIDs).
		Update("provisions_updated_at", time.Now().UTC()).Error
	if err != nil {
		return fmt.Errorf("failed to update provisions_updated_at for additiveIDs: %w", err)
	}

	return nil
}

func (r *provisionRepository) saveProvision(tx *gorm.DB, provision *data.Provision, ingredientsUpdated bool) error {
	if ingredientsUpdated {
		provision.IngredientsUpdatedAt = time.Now().UTC()
	}

	err := tx.Save(provision).Error
	if err != nil {
		return fmt.Errorf("failed to update provision: %w", err)
	}
	return nil
}

func (r *provisionRepository) updateProvisionIngredients(tx *gorm.DB, provisionID uint, ingredients []data.ProvisionIngredient) error {
	var existingIngredients []data.ProvisionIngredient
	if err := tx.Where("provision_id = ?", provisionID).Find(&existingIngredients).Error; err != nil {
		return fmt.Errorf("failed to fetch existing provision ingredients: %w", err)
	}

	existingMap := make(map[uint]data.ProvisionIngredient)
	for _, ing := range existingIngredients {
		existingMap[ing.IngredientID] = ing
	}

	var toInsert []data.ProvisionIngredient
	var toUpdate []data.ProvisionIngredient
	existingIDs := make(map[uint]struct{})

	for _, ing := range ingredients {
		existing, found := existingMap[ing.IngredientID]
		if found {
			if existing.Quantity != ing.Quantity {
				existing.Quantity = ing.Quantity
				toUpdate = append(toUpdate, existing)
			}
			existingIDs[ing.IngredientID] = struct{}{}
		} else {
			toInsert = append(toInsert, data.ProvisionIngredient{
				ProvisionID:  provisionID,
				IngredientID: ing.IngredientID,
				Quantity:     ing.Quantity,
			})
		}
	}

	var toDeleteIDs []uint
	for id := range existingMap {
		if _, ok := existingIDs[id]; !ok {
			toDeleteIDs = append(toDeleteIDs, id)
		}
	}

	if len(toUpdate) > 0 {
		if err := tx.Save(&toUpdate).Error; err != nil {
			return fmt.Errorf("failed to update provision ingredients: %w", err)
		}
	}

	if len(toDeleteIDs) > 0 {
		if err := tx.Unscoped().
			Where("provision_id = ? AND ingredient_id IN ?", provisionID, toDeleteIDs).
			Delete(&data.ProvisionIngredient{}).Error; err != nil {
			return fmt.Errorf("failed to delete old provision ingredients: %w", err)
		}
	}

	if len(toInsert) > 0 {
		if err := tx.Create(&toInsert).Error; err != nil {
			return fmt.Errorf("failed to insert new provision ingredients: %w", err)
		}
	}

	return nil
}

func (r *provisionRepository) DeleteProvision(provisionID uint) (*data.Provision, error) {
	var provision data.Provision

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", provisionID).
			Delete(&provision).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().
			Where("provision_id = ?", provisionID).
			Delete(&data.ProvisionIngredient{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().
			Where("provision_id = ?", provisionID).
			Delete(&data.AdditiveProvision{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().
			Where("provision_id = ?", provisionID).
			Delete(&data.ProductSizeProvision{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &provision, nil
}
