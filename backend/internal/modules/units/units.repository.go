package units

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UnitRepository interface {
	Create(unit *data.Unit) error
	GetAll(locale data.LanguageCode, filter *types.UnitFilter) ([]data.Unit, error)
	GetByID(id uint) (*data.Unit, error)
	GetTranslatedByID(locale data.LanguageCode, id uint) (*data.Unit, error)
	Update(id uint, updates data.Unit) error
	Delete(id uint) error
	FindRawUnitByID(id uint, unit *data.Unit) error

	CloneWithTransaction(tx *gorm.DB) UnitRepository
}

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) CloneWithTransaction(tx *gorm.DB) UnitRepository {
	return &unitRepository{db: tx}
}

func (r *unitRepository) FindRawUnitByID(id uint, unit *data.Unit) error {
	err := r.db.
		Where("id = ?", id).
		First(unit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.ErrUnitNotFound
		}
		return fmt.Errorf("failed to find unit by ID: %w", err)
	}
	return nil
}

func (r *unitRepository) Create(unit *data.Unit) error {
	return r.db.Create(unit).Error
}

func (r *unitRepository) GetAll(locale data.LanguageCode, filter *types.UnitFilter) ([]data.Unit, error) {
	var units []data.Unit

	q := r.db.Model(&data.Unit{})

	if filter.Search != nil && *filter.Search != "" {
		q = q.Where("name ILIKE ?", "%"+*filter.Search+"%")
	}

	paged, err := utils.ApplySortedPaginationForModel(
		q,
		filter.Pagination,
		filter.Sort,
		&data.Unit{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch units: %w", err)
	}

	paged = utils.ApplyLocalizedPreloads(
		paged,
		locale,
		types.UnitPreloadMap,
	)

	if err := paged.Find(&units).Error; err != nil {
		return nil, err
	}
	if units == nil {
		units = []data.Unit{}
	}
	return units, nil
}

func (r *unitRepository) GetByID(id uint) (*data.Unit, error) {
	var unit data.Unit
	err := r.db.First(&unit, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrUnitNotFound
		}
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepository) GetTranslatedByID(locale data.LanguageCode, id uint) (*data.Unit, error) {
	var unit data.Unit

	q := r.db.Model(&data.Unit{}).
		Where("id = ?", id)

	q = utils.ApplyLocalizedPreloads(
		q, locale, types.UnitPreloadMap)

	if err := q.First(&unit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrUnitNotFound
		}
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepository) Update(id uint, updates data.Unit) error {
	return r.db.Model(&data.Unit{}).Where("id = ?", id).Updates(updates).Error
}

func (r *unitRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := checkUnitReferences(tx, id); err != nil {
			return err
		}

		if err := tx.Delete(&data.Unit{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func checkUnitReferences(db *gorm.DB, unitID uint) error {
	var unit data.Unit

	err := db.
		Preload("StockMaterials", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Preload("Additives", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Preload("ProductSizes", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Where(&data.Unit{BaseEntity: data.BaseEntity{ID: unitID}}).
		First(&unit).Error
	if err != nil {
		return err
	}

	if len(unit.StockMaterials) > 0 || len(unit.Additives) > 0 ||
		len(unit.ProductSizes) > 0 || len(unit.Ingredients) > 0 {
		return types.ErrUnitIsInUse
	}

	return nil
}
