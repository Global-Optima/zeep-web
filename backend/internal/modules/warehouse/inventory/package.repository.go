package inventory

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type PackageRepository interface {
	GetAllPackages() ([]data.Package, error)
	GetPackageBySKU(skuID uint) (*data.Package, error)
	CreatePackage(pkg *data.Package) error
	UpdatePackage(pkg *data.Package) error
	DeletePackage(packageID uint) error
}

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) GetAllPackages() ([]data.Package, error) {
	var packages []data.Package
	if err := r.db.Find(&packages).Error; err != nil {
		return nil, err
	}
	return packages, nil
}

func (r *packageRepository) GetPackageBySKU(skuID uint) (*data.Package, error) {
	var pkg data.Package
	if err := r.db.Where("sku_id = ?", skuID).First(&pkg).Error; err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepository) CreatePackage(pkg *data.Package) error {
	return r.db.Create(pkg).Error
}

func (r *packageRepository) UpdatePackage(pkg *data.Package) error {
	return r.db.Save(pkg).Error
}

func (r *packageRepository) DeletePackage(packageID uint) error {
	return r.db.Delete(&data.Package{}, packageID).Error
}
