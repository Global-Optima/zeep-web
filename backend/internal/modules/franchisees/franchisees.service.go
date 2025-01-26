package franchisees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
)

type FranchiseeService interface {
	CreateFranchisee(dto *types.CreateFranchiseeDTO) (uint, error)
	UpdateFranchisee(id uint, dto *types.UpdateFranchiseeDTO) error
	DeleteFranchisee(id uint) error
	GetFranchiseeByID(id uint) (*types.FranchiseeDTO, error)
	GetFranchisees(filter *types.FranchiseeFilter) ([]types.FranchiseeDTO, error)
	IsFranchiseeStore(franchiseeID, storeID uint) (bool, error)
}

type franchiseeService struct {
	repo FranchiseeRepository
}

func NewFranchiseeService(repo FranchiseeRepository) FranchiseeService {
	return &franchiseeService{repo: repo}
}

func (s *franchiseeService) CreateFranchisee(dto *types.CreateFranchiseeDTO) (uint, error) {
	franchisee := types.CreateToFranchisee(dto)
	id, err := s.repo.CreateFranchisee(franchisee)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *franchiseeService) UpdateFranchisee(id uint, dto *types.UpdateFranchiseeDTO) error {
	updateData := &data.Franchisee{}
	if dto.Name != nil {
		updateData.Name = *dto.Name
	}
	if dto.Description != nil {
		updateData.Description = *dto.Description
	}
	return s.repo.UpdateFranchisee(id, updateData)
}

func (s *franchiseeService) DeleteFranchisee(id uint) error {
	return s.repo.DeleteFranchisee(id)
}

func (s *franchiseeService) GetFranchiseeByID(id uint) (*types.FranchiseeDTO, error) {
	franchisee, err := s.repo.GetFranchiseeByID(id)
	if err != nil {
		return nil, err
	}
	return types.ConvertFranchiseeToDTO(franchisee), nil
}

func (s *franchiseeService) GetFranchisees(filter *types.FranchiseeFilter) ([]types.FranchiseeDTO, error) {
	franchisees, err := s.repo.GetFranchisees(filter)
	if err != nil {
		return nil, err
	}
	dtos := make([]types.FranchiseeDTO, len(franchisees))
	for i, franchisee := range franchisees {
		dtos[i] = *types.ConvertFranchiseeToDTO(&franchisee)
	}
	return dtos, nil
}

func (s *franchiseeService) IsFranchiseeStore(franchiseeID, storeID uint) (bool, error) {
	return s.repo.IsFranchiseeStore(franchiseeID, storeID)
}
