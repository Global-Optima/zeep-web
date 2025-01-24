package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func CreateToRegion(dto *CreateRegionDTO) *data.Region {
	return &data.Region{
		Name: dto.Name,
	}
}

func UpdateToRegion(dto *UpdateRegionDTO) *data.Region {
	updateData := &data.Region{}
	if dto.Name != nil {
		updateData.Name = *dto.Name
	}
	return updateData
}

func MapRegionEntityToDTO(entity *data.Region) *RegionDTO {
	return &RegionDTO{
		ID:   entity.ID,
		Name: entity.Name,
	}
}
