package regions

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
	"github.com/gin-gonic/gin"
)

type RegionService interface {
	CreateRegion(dto *types.CreateRegionDTO) (uint, error)
	UpdateRegion(id uint, dto *types.UpdateRegionDTO) error
	DeleteRegion(id uint) error
	GetRegionByID(id uint) (*types.RegionDTO, error)
	GetRegions(filter *types.RegionFilter) ([]types.RegionDTO, error)
	GetAllRegions(filter *types.RegionFilter) ([]types.RegionDTO, error)
	IsRegionWarehouse(regionID, warehouseID uint) *handlerErrors.HandlerError
	CheckRegionWarehouse(c *gin.Context) (uint, *handlerErrors.HandlerError)
	CheckRegionWarehouseWithRole(c *gin.Context) (uint, data.EmployeeRole, *handlerErrors.HandlerError)
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

func (s *regionService) GetAllRegions(filter *types.RegionFilter) ([]types.RegionDTO, error) {
	regions, err := s.repo.GetAllRegions(filter)
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

func (s *regionService) IsRegionWarehouse(regionID, warehouseID uint) *handlerErrors.HandlerError {
	ok, err := s.repo.IsRegionWarehouse(regionID, warehouseID)
	if err != nil {
		return types.ErrFailedToCheckRegionWarehouse
	}
	if !ok {
		return types.ErrRegionWarehouseMismatch
	}
	return nil
}

func (s *regionService) CheckRegionWarehouse(c *gin.Context) (uint, *handlerErrors.HandlerError) {
	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, contexts.ErrUnauthorizedAccess
	}

	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		return 0, errH
	}

	if claims.Role == data.RoleRegionWarehouseManager {
		regionID, errH := contexts.GetRegionId(c)
		if errH != nil {
			return 0, errH
		}

		if regionID == nil {
			return warehouseID, nil
		}

		errH = s.IsRegionWarehouse(*regionID, warehouseID)
		if errH != nil {
			return 0, errH
		}
		return warehouseID, nil
	}

	return warehouseID, nil
}

func (s *regionService) CheckRegionWarehouseWithRole(c *gin.Context) (uint, data.EmployeeRole, *handlerErrors.HandlerError) {
	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, "", contexts.ErrUnauthorizedAccess
	}

	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		return 0, "", errH
	}

	if claims.Role == data.RoleRegionWarehouseManager {
		regionID, errH := contexts.GetRegionId(c)
		if errH != nil {
			return 0, "", errH
		}

		if regionID == nil {
			return warehouseID, claims.Role, nil
		}

		errH = s.IsRegionWarehouse(*regionID, warehouseID)
		if errH != nil {
			return 0, "", errH
		}
		return warehouseID, claims.Role, nil
	}

	return warehouseID, claims.Role, nil
}
