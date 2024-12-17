package supplier

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type SupplierRepository interface {
	CreateSupplier(supplier *data.Supplier) error
	GetSupplierByID(id uint) (*data.Supplier, error)
	UpdateSupplier(id uint, fields *data.Supplier) error
	DeleteSupplier(id uint) error
	GetAllSuppliers() ([]data.Supplier, error)
}

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &supplierRepository{db}
}

func (r *supplierRepository) CreateSupplier(supplier *data.Supplier) error {
	return r.db.Create(supplier).Error
}

func (r *supplierRepository) GetSupplierByID(id uint) (*data.Supplier, error) {
	var supplier data.Supplier
	err := r.db.First(&supplier, id).Error
	return &supplier, err
}

func (r *supplierRepository) UpdateSupplier(id uint, dto *data.Supplier) error {
	if err := r.db.Model(&data.Supplier{}).Where("id = ?", id).Updates(dto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("supplier not found")
		}
		return err
	}

	return nil
}

func (r *supplierRepository) DeleteSupplier(id uint) error {
	return r.db.Delete(&data.Supplier{}, id).Error
}

func (r *supplierRepository) GetAllSuppliers() ([]data.Supplier, error) {
	var suppliers []data.Supplier
	err := r.db.Find(&suppliers).Error
	return suppliers, err
}
