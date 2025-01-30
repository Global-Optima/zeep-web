package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type CreateRegionDTO struct {
	Name string `json:"name" binding:"required"`
}

type UpdateRegionDTO struct {
	Name *string `json:"name,omitempty"`
}

type RegionDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type RegionFilter struct {
	utils.BaseFilter
	Name   *string `form:"name,omitempty"`
	Search *string `form:"search,omitempty"`
}
