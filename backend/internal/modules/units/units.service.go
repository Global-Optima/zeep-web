package units

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
)

type UnitService interface {
	Create(dto types.CreateUnitDTO) (uint, error)
	GetAll(locale data.LanguageCode, filter *types.UnitFilter) ([]types.UnitsDTO, error)
	GetByID(id uint) (*types.UnitsDTO, error)
	GetTranslatedByID(locale data.LanguageCode, id uint) (*types.UnitsDTO, error)
	Update(id uint, dto types.UpdateUnitDTO) error
	Delete(id uint) error
	UpsertUnitTranslations(unitID uint, dto *types.UnitTranslationsDTO) error
}

type unitService struct {
	repo               UnitRepository
	transactionManager TransactionManager
	logger             *zap.SugaredLogger
}

func NewUnitService(repo UnitRepository, transactionManager TransactionManager, logger *zap.SugaredLogger) UnitService {
	return &unitService{
		repo:               repo,
		transactionManager: transactionManager,
		logger:             logger,
	}
}

func (s *unitService) Create(dto types.CreateUnitDTO) (uint, error) {
	unit := data.Unit{
		Name:             dto.Name,
		ConversionFactor: dto.ConversionFactor,
	}

	if err := s.repo.Create(&unit); err != nil {
		wrappedErr := fmt.Errorf("failed to create unit: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return unit.ID, nil
}

func (s *unitService) GetAll(locale data.LanguageCode, filter *types.UnitFilter) ([]types.UnitsDTO, error) {
	units, err := s.repo.GetAll(locale, filter)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get units: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	return types.ToUnitResponses(units), nil
}

func (s *unitService) GetByID(id uint) (*types.UnitsDTO, error) {
	unit, err := s.repo.GetByID(id)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get unit: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	response := types.ToUnitResponse(*unit)
	return &response, nil
}

func (s *unitService) GetTranslatedByID(locale data.LanguageCode, id uint) (*types.UnitsDTO, error) {
	unit, err := s.repo.GetTranslatedByID(locale, id)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get translated unit: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
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
		wrappedErr := fmt.Errorf("failed to update unit: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *unitService) Delete(id uint) error {
	err := s.repo.Delete(id)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete unit: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *unitService) UpsertUnitTranslations(unitID uint, dto *types.UnitTranslationsDTO) error {
	if dto == nil {
		return fmt.Errorf("translations DTO is nil")
	}
	if err := s.transactionManager.UpsertUnitTranslations(unitID, dto); err != nil {
		wrappedErr := fmt.Errorf("failed to upsert additive category translations: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}
