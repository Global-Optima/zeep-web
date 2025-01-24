package franchisees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
)

type FranchiseeService interface {
	Create(dto *types.CreateFranchiseeDTO) (uint, error)
	Update(id uint, dto *types.UpdateFranchiseeDTO) error
	Delete(id uint) error
	GetByID(id uint) (*types.FranchiseeDTO, error)
	GetAll(filter *types.FranchiseeFilter) ([]types.FranchiseeDTO, error)
}

type franchiseeService struct {
	repo FranchiseeRepository
}

func NewFranchiseeService(repo FranchiseeRepository) FranchiseeService {
	return &franchiseeService{repo: repo}
}

func (s *franchiseeService) Create(dto *types.CreateFranchiseeDTO) (uint, error) {
	franchisee := types.CreateToFranchisee(dto)
	id, err := s.repo.Create(franchisee)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *franchiseeService) Update(id uint, dto *types.UpdateFranchiseeDTO) error {
	updateData := &data.Franchisee{}
	if dto.Name != nil {
		updateData.Name = *dto.Name
	}
	if dto.Description != nil {
		updateData.Description = *dto.Description
	}
	return s.repo.Update(id, updateData)
}

func (s *franchiseeService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *franchiseeService) GetByID(id uint) (*types.FranchiseeDTO, error) {
	franchisee, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return types.ConvertFranchiseeToDTO(franchisee), nil
}

func (s *franchiseeService) GetAll(filter *types.FranchiseeFilter) ([]types.FranchiseeDTO, error) {
	franchisees, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, err
	}
	dtos := make([]types.FranchiseeDTO, len(franchisees))
	for i, franchisee := range franchisees {
		dtos[i] = *types.ConvertFranchiseeToDTO(&franchisee)
	}
	return dtos, nil
}
