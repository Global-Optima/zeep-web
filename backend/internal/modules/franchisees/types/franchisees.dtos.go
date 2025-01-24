package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type CreateFranchiseeDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type UpdateFranchiseeDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type FranchiseeDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FranchiseeFilter struct {
	utils.BaseFilter
	Name   *string `form:"name,omitempty"`
	Search *string `form:"search,omitempty"`
}
