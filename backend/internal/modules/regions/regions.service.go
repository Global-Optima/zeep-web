package regions

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
)

type RegionService interface {
	Create(dto *types.CreateRegionDTO) (*types.RegionDTO, error)
	Update(id uint, dto *types.UpdateRegionDTO) error
	Delete(id uint) error
	GetByID(id uint) (*types.RegionDTO, error)
	GetAll(filter *types.RegionFilter) ([]types.RegionDTO, error)
}

type regionService struct {
	repo RegionRepository
}

func NewRegionService(repo RegionRepository) RegionService {
	return &regionService{repo: repo}
}

func (s *regionService) Create(dto *types.CreateRegionDTO) (*types.RegionDTO, error) {
	region := &data.Region{
		Name: dto.Name,
	}
	if err := s.repo.Create(region); err != nil {
		return nil, err
	}
	return &types.RegionDTO{
		ID:   region.ID,
		Name: region.Name,
	}, nil
}

func (s *regionService) Update(id uint, dto *types.UpdateRegionDTO) error {
	updateData := &data.Region{}
	if dto.Name != nil {
		updateData.Name = *dto.Name
	}
	return s.repo.Update(id, updateData)
}

func (s *regionService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *regionService) GetByID(id uint) (*types.RegionDTO, error) {
	region, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &types.RegionDTO{
		ID:   region.ID,
		Name: region.Name,
	}, nil
}

func (s *regionService) GetAll(filter *types.RegionFilter) ([]types.RegionDTO, error) {
	regions, err := s.repo.GetAll(filter)
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
