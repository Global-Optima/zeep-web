package storeProvisions

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type StoreProvisionRepository interface {
	CreateStoreProvision(storeProvision *data.StoreProvision) (uint, error)
	GetStoreProvisions(storeID uint, filter *types.StoreProvisionFilterDTO) ([]data.StoreProvision, error)
	GetStoreProvisionByID(storeID uint, storeProvisionID uint) (*data.StoreProvision, error)
	GetStoreProvisionWithDetailsByID(storeID, storeProvisionID uint) (*data.StoreProvision, error)
	SaveStoreProvisionWithAssociations(updateModels *types.StoreProvisionModels) error
	DeleteStoreProvision(storeID, storeProvisionID uint) (*data.StoreProvision, error)
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

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.MinCompletedAt != nil {
		query = query.Where("completed_at >= ?", *filter.MinCompletedAt)
	}

	if filter.MaxCompletedAt != nil {
		query = query.Where("completed_at <= ?", *filter.MaxCompletedAt)
	}

	if filter.IsExpired != nil {
		now := time.Now()
		if *filter.IsExpired {
			query = query.Where("expires_at <= ?", now)
		} else {
			query = query.Where("expires_at > ?", now)
		}
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
		if updateModels.StoreProvision != nil || updateModels.StoreProvision.ID != 0 {
			err := tx.Save(updateModels.StoreProvision).Error
			if err != nil {
				return err
			}
		}

		if updateModels.StoreProvisionIngredientsMultiplier != nil && *updateModels.StoreProvisionIngredientsMultiplier != 1 {
			err := r.updateStoreProvisionIngredients(
				tx, updateModels.StoreProvision.StoreID, updateModels.StoreProvision.StoreID, *updateModels.StoreProvisionIngredientsMultiplier,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *storeProvisionRepository) updateStoreProvisionIngredients(tx *gorm.DB, storeID, storeProvisionID uint, multiplier float64) error {
	return tx.Model(&data.StoreProvisionIngredient{}).
		Where("store_provision_id = ? AND store_id = ?", storeProvisionID, storeID).
		Update("volume", gorm.Expr("volume * ?", multiplier)).Error
}

func (r *storeProvisionRepository) SaveStoreProvision(storeProvision *data.StoreProvision) error {
	if storeProvision == nil || storeProvision.ID == 0 {
		return fmt.Errorf("invalid store provision model")
	}

	err := r.db.Model(&data.StoreProvision{}).
		Save(storeProvision).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.ErrStoreProvisionNotFound
		}
		return err
	}

	return nil
}

func (r *storeProvisionRepository) DeleteStoreProvision(storeID, storeProvisionID uint) (*data.StoreProvision, error) {
	var sp data.StoreProvision

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND store_id = ?", storeProvisionID, storeID).
			First(&sp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return types.ErrStoreProvisionNotFound
			}
			return err
		}

		if err := tx.Where("store_provision_id = ?", storeProvisionID).
			Delete(&data.StoreProvisionIngredient{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&sp).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &sp, nil
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
