package regions

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type RegionRepository interface {
	Create(region *data.Region) error
	Update(id uint, updateData *data.Region) error
	Delete(id uint) error
	GetByID(id uint) (*data.Region, error)
	GetAll(filter *types.RegionFilter) ([]data.Region, error)
}

type regionRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) RegionRepository {
	return &regionRepository{db: db}
}

func (r *regionRepository) Create(region *data.Region) error {
	return r.db.Create(region).Error
}

func (r *regionRepository) Update(id uint, updateData *data.Region) error {
	return r.db.Model(&data.Region{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *regionRepository) Delete(id uint) error {
	return r.db.Delete(&data.Region{}, id).Error
}

func (r *regionRepository) GetByID(id uint) (*data.Region, error) {
	var region data.Region
	if err := r.db.Preload("Warehouses").First(&region, id).Error; err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *regionRepository) GetAll(filter *types.RegionFilter) ([]data.Region, error) {
	var regions []data.Region
	query := r.db.Model(&data.Region{}).Preload("Warehouses")

	if filter.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
	}

	if filter.Search != nil {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ?", searchTerm)
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Region{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}
