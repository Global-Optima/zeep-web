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
	GetAll(filter *types.UnitFilter) ([]data.Unit, error)
	GetByID(id uint) (*data.Unit, error)
	Update(id uint, updates data.Unit) error
	Delete(id uint) error
}

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) Create(unit *data.Unit) error {
	return r.db.Create(unit).Error
}

func (r *unitRepository) GetAll(filter *types.UnitFilter) ([]data.Unit, error) {
	var units []data.Unit
	query := r.db.Model(&data.Unit{})

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &units)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch units: %w", err)
	}

	if filter.Search != nil {
		query = query.Where("name ILIKE ?", "%"+*filter.Search+"%")
	}

	if err := query.Find(&units).Error; err != nil {
		return nil, err
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
