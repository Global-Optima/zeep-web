package regions

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
)

type RegionService interface {
	CreateRegion(dto *types.CreateRegionDTO) (uint, error)
	UpdateRegion(id uint, dto *types.UpdateRegionDTO) error
	DeleteRegion(id uint) error
	GetRegionByID(id uint) (*types.RegionDTO, error)
	GetRegions(filter *types.RegionFilter) ([]types.RegionDTO, error)
	IsRegionWarehouse(regionID uint, warehouseID uint) (bool, error)
}

type regionService struct {
	repo RegionRepository
}

func NewRegionService(repo RegionRepository) RegionService {
	return &regionService{repo: repo}
}

func (s *regionService) CreateRegion(dto *types.CreateRegionDTO) (uint, error) {
	region := &data.Region{
		Name: dto.Name,
	}
	id, err := s.repo.CreateRegion(region)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *regionService) UpdateRegion(id uint, dto *types.UpdateRegionDTO) error {
	updateData := &data.Region{}
	if dto.Name != nil {
		updateData.Name = *dto.Name
	}
	return s.repo.UpdateRegion(id, updateData)
}

func (s *regionService) DeleteRegion(id uint) error {
	return s.repo.DeleteRegion(id)
}

func (s *regionService) GetRegionByID(id uint) (*types.RegionDTO, error) {
	region, err := s.repo.GetRegionByID(id)
	if err != nil {
		return nil, err
	}
	return &types.RegionDTO{
		ID:   region.ID,
		Name: region.Name,
	}, nil
}

func (s *regionService) GetRegions(filter *types.RegionFilter) ([]types.RegionDTO, error) {
	regions, err := s.repo.GetRegions(filter)
	if err != nil {
		return nil, err
	}
	dtos := make([]types.RegionDTO, len(regions))
	for i, region := range regions {
		dtos[i] = types.RegionDTO{
			ID:   region.ID,
			Name: region.Name,
		}
	}
	return dtos, nil
}

func (s *regionService) IsRegionWarehouse(regionID uint, warehouseID uint) (bool, error) {
	return s.repo.IsRegionWarehouse(regionID, warehouseID)
}
