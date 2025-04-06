package regions

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type RegionRepository interface {
	CreateRegion(region *data.Region) (uint, error)
	UpdateRegion(id uint, updateData *data.Region) error
	DeleteRegion(id uint) error
	GetRegionByID(id uint) (*data.Region, error)
	GetRegions(filter *types.RegionFilter) ([]data.Region, error)
	GetAllRegions(filter *types.RegionFilter) ([]data.Region, error)
	IsRegionWarehouse(regionID uint, warehouseID uint) (bool, error)
}

type regionRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) RegionRepository {
	return &regionRepository{db: db}
}

func (r *regionRepository) CreateRegion(region *data.Region) (uint, error) {
	err := r.db.Create(region).Error
	if err != nil {
		return 0, err
	}
	return region.ID, nil
}

func (r *regionRepository) UpdateRegion(id uint, updateData *data.Region) error {
	return r.db.Model(&data.Region{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *regionRepository) DeleteRegion(id uint) error {
	return r.db.Delete(&data.Region{}, id).Error
}

func (r *regionRepository) GetRegionByID(id uint) (*data.Region, error) {
	var region data.Region
	if err := r.db.Preload("Warehouses").First(&region, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrRegionNotFound
		}
		return nil, err
	}
	return &region, nil
}

func (r *regionRepository) GetRegions(filter *types.RegionFilter) ([]data.Region, error) {
	var regions []data.Region
	query := r.db.Model(&data.Region{})

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

func (r *regionRepository) GetAllRegions(filter *types.RegionFilter) ([]data.Region, error) {
	var regions []data.Region
	query := r.db.Model(&data.Region{})

	if filter == nil {
		return nil, fmt.Errorf("fitler is nil")
	}

	if filter.Search != nil {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ?", searchTerm)
	}

	if err := query.Scopes(filter.Sort.SortGorm()).Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

func (r *regionRepository) IsRegionWarehouse(regionID uint, warehouseID uint) (bool, error) {
	var warehouse data.Warehouse
	err := r.db.Model(&data.Warehouse{}).
		Where("region_id = ? AND id = ?", regionID, warehouseID).
		First(&warehouse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
