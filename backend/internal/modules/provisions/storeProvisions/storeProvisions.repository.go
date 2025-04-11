package storeProvisions

import (
	"fmt"
	"slices"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreProvisionRepository interface {
	CreateStoreProvision(storeProvision *data.StoreProvision) (uint, error)
	GetStoreProvisions(storeID uint, filter *types.StoreProvisionFilterDTO) ([]data.StoreProvision, error)
	GetStoreProvisionByID(storeID uint, storeProvisionID uint) (*data.StoreProvision, error)
	GetStoreProvisionWithDetailsByID(storeID, storeProvisionID uint) (*data.StoreProvision, error)
	GetAllStoreProvisionList(storeID uint) ([]data.StoreProvision, error)
	SaveStoreProvisionWithAssociations(updateModels *types.StoreProvisionModels) error
	HardDeleteStoreProvision(storeProvisionID uint) error
	CountStoreProvisionsToday(storeID, provisionID uint) (uint, error)
	SaveStoreProvision(storeProvision *data.StoreProvision) error

	CloneWithTransaction(tx *gorm.DB) StoreProvisionRepository
}

type storeProvisionRepository struct {
	db *gorm.DB
}

func NewStoreProvisionRepository(db *gorm.DB) StoreProvisionRepository {
	return &storeProvisionRepository{db: db}
}

func (r *storeProvisionRepository) CloneWithTransaction(tx *gorm.DB) StoreProvisionRepository {
	return &storeProvisionRepository{
		db: tx,
	}
}

func (r *storeProvisionRepository) CreateStoreProvision(storeProvision *data.StoreProvision) (uint, error) {
	err := r.db.Create(storeProvision).Error
	if err != nil {
		return 0, err
	}
	return storeProvision.ID, nil
}

func (r *storeProvisionRepository) GetStoreProvisions(storeID uint, filter *types.StoreProvisionFilterDTO) ([]data.StoreProvision, error) {
	var provisions []data.StoreProvision

	query := r.db.Model(&data.StoreProvision{}).
		Preload("Provision.Unit").
		Where("store_id = ?", storeID)

	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	if len(filter.Statuses) > 0 {
		if slices.Contains(filter.Statuses, data.STORE_PROVISION_VISUAL_STATUS_EXPIRED) {
			query = query.Where("status IN (?)", filter.Statuses)
		} else {
			query = query.Where("status IN (?) AND (expires_at > ? OR expires_at IS NULL)", filter.Statuses, time.Now().UTC())
		}
	}

	if filter.MinCompletedAt != nil {
		query = query.Where("completed_at >= ?", *filter.MinCompletedAt)
	}

	if filter.MaxCompletedAt != nil {
		query = query.Where("completed_at <= ?", *filter.MaxCompletedAt)
	}

	if filter.Search != nil {
		pattern := "%" + *filter.Search + "%"
		query = query.Joins("JOIN provisions ON provisions.id = store_provisions.provision_id").
			Where("provisions.name ILIKE ?", pattern)
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreProvision{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	err = query.Find(&provisions).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStoreProvisionNotFound
		}
		return nil, err
	}

	return provisions, nil
}

func (r *storeProvisionRepository) GetStoreProvisionByID(storeID, storeProvisionID uint) (*data.StoreProvision, error) {
	var sp data.StoreProvision
	err := r.db.Model(&data.StoreProvision{}).
		Where("id = ? AND store_id = ?", storeProvisionID, storeID).
		First(&sp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStoreProvisionNotFound
		}
		return nil, err
	}
	return &sp, nil
}

func (r *storeProvisionRepository) GetStoreProvisionWithDetailsByID(storeID, storeProvisionID uint) (*data.StoreProvision, error) {
	var sp data.StoreProvision
	err := r.db.Preload("Provision.Unit").
		Preload("StoreProvisionIngredients.Ingredient.Unit").
		Preload("StoreProvisionIngredients.Ingredient.IngredientCategory").
		Where("id = ? AND store_id = ?", storeProvisionID, storeID).
		First(&sp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStoreProvisionNotFound
		}
		return nil, err
	}
	return &sp, nil
}

func (r *storeProvisionRepository) SaveStoreProvisionWithAssociations(updateModels *types.StoreProvisionModels) error {
	if updateModels == nil {
		return fmt.Errorf("nothing to update")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		sp := updateModels.StoreProvision
		multiplier := updateModels.StoreProvisionIngredientsMultiplier

		if sp != nil && sp.ID != 0 && sp.StoreID != 0 {
			if err := tx.Save(sp).Error; err != nil {
				return fmt.Errorf("failed to save store provision: %w", err)
			}

			if multiplier != nil {
				if err := r.updateStoreProvisionIngredients(tx, sp.ID, *multiplier); err != nil {
					return fmt.Errorf("failed to update ingredients: %w", err)
				}
			}
		}

		return nil
	})
}

func (r *storeProvisionRepository) updateStoreProvisionIngredients(tx *gorm.DB, storeProvisionID uint, multiplier float64) error {
	return tx.Model(&data.StoreProvisionIngredient{}).
		Where("store_provision_id = ?", storeProvisionID).
		Update("quantity", gorm.Expr("initial_quantity * ?", multiplier)).Error
}

func (r *storeProvisionRepository) SaveStoreProvision(storeProvision *data.StoreProvision) error {
	if storeProvision == nil || storeProvision.ID == 0 {
		return fmt.Errorf("invalid store provision model")
	}

	err := r.db.Save(storeProvision).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.ErrStoreProvisionNotFound
		}
		return err
	}

	return nil
}

func (r *storeProvisionRepository) HardDeleteStoreProvision(storeProvisionID uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().
			Where("store_provision_id = ?", storeProvisionID).
			Delete(&data.StoreProvisionIngredient{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().
			Where("id = ?", storeProvisionID).
			Delete(&data.StoreProvision{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *storeProvisionRepository) CountStoreProvisionsToday(storeID, provisionID uint) (uint, error) {
	if provisionID == 0 {
		return 0, fmt.Errorf("provision ID cannot be zero")
	}

	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var count int64
	err := r.db.Model(&data.StoreProvision{}).
		Joins("JOIN provisions ON provisions.id = store_provisions.provision_id").
		Where("store_provisions.provision_id = ? AND store_provisions.store_id = ? AND store_provisions.created_at >= ? AND store_provisions.created_at < ?",
			provisionID, storeID, startOfDay, endOfDay).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count today's store provisions: %w", err)
	}

	return uint(count), nil
}

func (r *storeProvisionRepository) GetAllStoreProvisionList(storeID uint) ([]data.StoreProvision, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	var storeProvisionList []data.StoreProvision

	query := r.db.Model(&data.StoreProvision{}).
		Preload("Provision.Unit").
		Preload("Store").
		Where("store_id = ?", storeID)

	err := query.Find(&storeProvisionList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve store provison list: %w", err)
	}

	return storeProvisionList, nil
}
