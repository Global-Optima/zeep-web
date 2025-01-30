package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ConvertFranchiseeToDTO(franchisee *data.Franchisee) *FranchiseeDTO {
	return &FranchiseeDTO{
		ID:          franchisee.ID,
		Name:        franchisee.Name,
		Description: franchisee.Description,
	}
}

func CreateToFranchisee(dto *CreateFranchiseeDTO) *data.Franchisee {
	return &data.Franchisee{
		Name:        dto.Name,
		Description: dto.Description,
	}
}

func UpdateToFranchisee(dto *UpdateFranchiseeDTO) *data.Franchisee {
	updateData := &data.Franchisee{}
	if dto.Name != nil {
		updateData.Name = *dto.Name
	}
	if dto.Description != nil {
		updateData.Description = *dto.Description
	}
	return updateData
}
