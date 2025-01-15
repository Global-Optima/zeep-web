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
	ExistsByContactPhone(phone string) (bool, error)

	CreateSupplierMaterial(material *data.SupplierMaterial) error
	CreateSupplierPrice(price *data.SupplierPrice) error
	GetMaterialsBySupplier(supplierID uint) ([]data.SupplierMaterial, error)
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

func (r *supplierRepository) ExistsByContactPhone(phone string) (bool, error) {
	var count int64
	err := r.db.Model(&data.Supplier{}).Where("contact_phone = ?", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *supplierRepository) CreateSupplierMaterial(material *data.SupplierMaterial) error {
	return r.db.Create(material).Error
}

func (r *supplierRepository) CreateSupplierPrice(price *data.SupplierPrice) error {
	return r.db.Create(price).Error
}

func (r *supplierRepository) GetMaterialsBySupplier(supplierID uint) ([]data.SupplierMaterial, error) {
	var materials []data.SupplierMaterial
	err := r.db.Preload("StockMaterial").
		Preload("StockMaterial.StockMaterialCategory").
		Preload("StockMaterial.Package").
		Preload("StockMaterial.Package.Unit").
		Preload("SupplierPrices").
		Where("supplier_id = ?", supplierID).
		Find(&materials).Error
	return materials, err
}
