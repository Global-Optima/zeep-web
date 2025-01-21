package stockMaterialPackage

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type StockMaterialPackageRepository interface {
	Create(packageEntity *data.StockMaterialPackage) error
	CreateMultiplePackages(packages []data.StockMaterialPackage) error
	GetByID(id uint) (*data.StockMaterialPackage, error)
	Update(id uint, packageEntity data.StockMaterialPackage) error
	Delete(id uint) error
	GetAll() ([]data.StockMaterialPackage, error)
}

type stockMaterialPackageRepository struct {
	db *gorm.DB
}

func NewStockMaterialPackageRepository(db *gorm.DB) StockMaterialPackageRepository {
	return &stockMaterialPackageRepository{db: db}
}

func (r *stockMaterialPackageRepository) Create(packageEntity *data.StockMaterialPackage) error {
	return r.db.Create(packageEntity).Error
}

func (r *stockMaterialPackageRepository) CreateMultiplePackages(packages []data.StockMaterialPackage) error {
	return r.db.Create(packages).Error
}

func (r *stockMaterialPackageRepository) GetByID(id uint) (*data.StockMaterialPackage, error) {
	var packageEntity data.StockMaterialPackage
	err := r.db.Preload("StockMaterial").Preload("Unit").First(&packageEntity, id).Error
	if err != nil {
		return nil, err
	}
	return &packageEntity, nil
}

func (r *stockMaterialPackageRepository) Update(id uint, packageEntity data.StockMaterialPackage) error {
	return r.db.Model(&data.StockMaterialPackage{}).Where("id = ?", id).Updates(packageEntity).Error
}

func (r *stockMaterialPackageRepository) Delete(id uint) error {
	return r.db.Delete(&data.StockMaterialPackage{}, id).Error
}

func (r *stockMaterialPackageRepository) GetAll() ([]data.StockMaterialPackage, error) {
	var packages []data.StockMaterialPackage
	err := r.db.Preload("StockMaterial").Preload("Unit").Find(&packages).Error
	return packages, err
}
