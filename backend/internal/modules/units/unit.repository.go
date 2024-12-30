package units

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type UnitRepository interface {
	Create(unit *data.Unit) error
	GetAll() ([]data.Unit, error)
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

func (r *unitRepository) GetAll() ([]data.Unit, error) {
	var units []data.Unit
	err := r.db.Find(&units).Error
	return units, err
}

func (r *unitRepository) GetByID(id uint) (*data.Unit, error) {
	var unit data.Unit
	err := r.db.First(&unit, id).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepository) Update(id uint, updates data.Unit) error {
	return r.db.Model(&data.Unit{}).Where("id = ?", id).Updates(updates).Error
}

func (r *unitRepository) Delete(id uint) error {
	return r.db.Delete(&data.Unit{}, id).Error
}
