package inventory

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type PackageRepository interface {
	GetAllPackages() ([]data.StockMaterialPackage, error)
	GetPackageByStockMaterial(stockMaterialID uint) (*data.StockMaterialPackage, error)
	CreatePackage(pkg *data.StockMaterialPackage) error
	UpdatePackage(pkg *data.StockMaterialPackage) error
	DeletePackage(packageID uint) error
}

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) GetAllPackages() ([]data.StockMaterialPackage, error) {
	var packages []data.StockMaterialPackage
	if err := r.db.Find(&packages).Error; err != nil {
		return nil, err
	}
	return packages, nil
}

func (r *packageRepository) GetPackageByStockMaterial(stockMaterialID uint) (*data.StockMaterialPackage, error) {
	var pkg data.StockMaterialPackage
	if err := r.db.Where("stock_material_id = ?", stockMaterialID).First(&pkg).Error; err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepository) CreatePackage(pkg *data.StockMaterialPackage) error {
	return r.db.Create(pkg).Error
}

func (r *packageRepository) UpdatePackage(pkg *data.StockMaterialPackage) error {
	return r.db.Save(pkg).Error
}

func (r *packageRepository) DeletePackage(packageID uint) error {
	return r.db.Delete(&data.StockMaterialPackage{}, packageID).Error
}
