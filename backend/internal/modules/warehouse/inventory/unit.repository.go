package inventory

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type UnitRepository interface {
	GetAllUnits() ([]data.Unit, error)
	GetUnitByID(unitID uint) (*data.Unit, error)
	CreateUnit(unit *data.Unit) error
	UpdateUnit(unit *data.Unit) error
	DeleteUnit(unitID uint) error
}

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) GetAllUnits() ([]data.Unit, error) {
	var units []data.Unit
	if err := r.db.Find(&units).Error; err != nil {
		return nil, err
	}
	return units, nil
}

func (r *unitRepository) GetUnitByID(unitID uint) (*data.Unit, error) {
	var unit data.Unit
	if err := r.db.First(&unit, unitID).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepository) CreateUnit(unit *data.Unit) error {
	return r.db.Create(unit).Error
}

func (r *unitRepository) UpdateUnit(unit *data.Unit) error {
	return r.db.Save(unit).Error
}

func (r *unitRepository) DeleteUnit(unitID uint) error {
	return r.db.Delete(&data.Unit{}, unitID).Error
}
