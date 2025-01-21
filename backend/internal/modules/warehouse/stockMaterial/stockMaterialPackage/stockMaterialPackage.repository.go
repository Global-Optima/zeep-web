package stockMaterialPackage

import (
	"fmt"

	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	"gorm.io/gorm"
)

type StockMaterialPackageRepository interface {
	Create(packageEntity *data.StockMaterialPackage) error
	CreateMultiplePackages(packages []data.StockMaterialPackage) error
	GetByID(id uint) (*data.StockMaterialPackage, error)
	Update(id uint, packageEntity data.StockMaterialPackage) error
	UpsertPackages(stockMaterialID uint, packages []types.UpdateStockMaterialPackagesDTO) error
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

func (r *stockMaterialPackageRepository) UpsertPackages(stockMaterialID uint, packages []types.UpdateStockMaterialPackagesDTO) error {
	var existingPackages []data.StockMaterialPackage
	if err := r.db.Where("stock_material_id = ?", stockMaterialID).Find(&existingPackages).Error; err != nil {
		return fmt.Errorf("failed to fetch existing packages: %w", err)
	}

	existingMap := make(map[uint]data.StockMaterialPackage)
	for _, pkg := range existingPackages {
		existingMap[pkg.ID] = pkg
	}

	processedIDs := make(map[uint]bool)

	for _, pkgDTO := range packages {
		if err := types.ValidatePackageDTO(pkgDTO); err != nil {
			return fmt.Errorf("invalid package: %w", err)
		}

		if pkgDTO.StockMaterialPackageID != nil {
			existing, exists := existingMap[*pkgDTO.StockMaterialPackageID]
			if !exists {
				return fmt.Errorf("package with ID %d not found", *pkgDTO.StockMaterialPackageID)
			}

			if pkgDTO.Size != nil {
				existing.Size = *pkgDTO.Size
			}
			if pkgDTO.UnitID != nil {
				existing.UnitID = *pkgDTO.UnitID
			}

			if err := r.db.Model(&existing).Updates(existing).Error; err != nil {
				return fmt.Errorf("failed to update package with ID %d: %w", existing.ID, err)
			}

			processedIDs[*pkgDTO.StockMaterialPackageID] = true
		} else {
			newPackage := data.StockMaterialPackage{
				StockMaterialID: stockMaterialID,
				Size:            *pkgDTO.Size,
				UnitID:          *pkgDTO.UnitID,
			}
			if err := r.db.Create(&newPackage).Error; err != nil {
				return fmt.Errorf("failed to create package: %w", err)
			}
		}
	}

	for _, existing := range existingPackages {
		if !processedIDs[existing.ID] {
			if err := r.db.Delete(&existing).Error; err != nil {
				return fmt.Errorf("failed to delete package: %w", err)
			}
		}
	}

	return nil
}
