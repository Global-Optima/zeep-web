package units

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"gorm.io/gorm"
)

type UnitService interface {
	Create(dto types.CreateUnitDTO) (uint, error)
	GetAll() ([]types.UnitResponse, error)
	GetByID(id uint) (*types.UnitResponse, error)
	Update(id uint, dto types.UpdateUnitDTO) error
	Delete(id uint) error
}

type unitService struct {
	repo UnitRepository
}

func NewUnitService(repo UnitRepository) UnitService {
	return &unitService{repo: repo}
}

func (s *unitService) Create(dto types.CreateUnitDTO) (uint, error) {
	unit := data.Unit{
		Name:             dto.Name,
		ConversionFactor: dto.ConversionFactor,
	}

	if err := s.repo.Create(&unit); err != nil {
		return 0, fmt.Errorf("failed to create unit: %w", err)
	}

	return unit.ID, nil
}

func (s *unitService) GetAll() ([]types.UnitResponse, error) {
	units, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch units: %w", err)
	}
	return types.ToUnitResponses(units), nil
}

func (s *unitService) GetByID(id uint) (*types.UnitResponse, error) {
	unit, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("unit not found")
		}
		return nil, fmt.Errorf("failed to fetch unit: %w", err)
	}
	response := types.ToUnitResponse(*unit)
	return &response, nil
}

func (s *unitService) Update(id uint, dto types.UpdateUnitDTO) error {
	unit := data.Unit{}
	if dto.Name != nil {
		unit.Name = *dto.Name
	}
	if dto.ConversionFactor != nil {
		unit.ConversionFactor = *dto.ConversionFactor
	}

	if err := s.repo.Update(id, unit); err != nil {
		return fmt.Errorf("failed to update unit: %w", err)
	}
	return nil
}

func (s *unitService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete unit: %w", err)
	}
	return nil
}
