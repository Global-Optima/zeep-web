package additives

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type AdditiveService interface {
	GetAdditiveCategories(filter *types.AdditiveCategoriesFilterQuery) ([]types.AdditiveCategoryDTO, error)
	CreateAdditiveCategory(dto *types.CreateAdditiveCategoryDTO) (uint, error)
	UpdateAdditiveCategory(id uint, dto *types.UpdateAdditiveCategoryDTO) error
	DeleteAdditiveCategory(categoryID uint) error
	GetAdditiveCategoryByID(categoryID uint) (*types.AdditiveCategoryResponseDTO, error)

	GetAdditives(filter *types.AdditiveFilterQuery) ([]types.AdditiveDTO, error)
	GetAdditiveByID(additiveID uint) (*types.AdditiveDetailsDTO, error)
	GetAdditivesByIDs(additiveIDs []uint) ([]types.AdditiveDTO, error)
	CreateAdditive(dto *types.CreateAdditiveDTO) (uint, error)
	UpdateAdditive(additiveID uint, dto *types.UpdateAdditiveDTO) error
	DeleteAdditive(additiveID uint) error
}

type additiveService struct {
	repo   AdditiveRepository
	logger *zap.SugaredLogger
}

func NewAdditiveService(repo AdditiveRepository, logger *zap.SugaredLogger) AdditiveService {
	return &additiveService{
		repo:   repo,
		logger: logger,
	}
}

func (s *additiveService) GetAdditiveCategories(filter *types.AdditiveCategoriesFilterQuery) ([]types.AdditiveCategoryDTO, error) {
	// Fetch raw data from the repository
	categories, err := s.repo.GetAdditiveCategories(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve additives", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	// Handle case where no categories are found
	if len(categories) == 0 {
		return []types.AdditiveCategoryDTO{}, nil
	}

	// Convert raw data into DTOs
	var categoryDTOs []types.AdditiveCategoryDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, *types.ConvertToAdditiveCategoryDTO(&category))
	}

	return categoryDTOs, nil
}

func (s *additiveService) CreateAdditiveCategory(dto *types.CreateAdditiveCategoryDTO) (uint, error) {
	category := types.ConvertToAdditiveCategoryModel(dto)

	id, err := s.repo.CreateAdditiveCategory(category)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create additive category", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	return id, nil
}

func (s *additiveService) UpdateAdditiveCategory(id uint, dto *types.UpdateAdditiveCategoryDTO) error {
	existingCategory, err := s.repo.GetAdditiveCategoryByID(id)
	if err != nil {
		wrappedErr := utils.WrapError("failed to fetch existing additive categor", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	if existingCategory == nil {
		return fmt.Errorf("additive category with ID %d not found", id)
	}

	updatedCategory := types.ConvertToUpdatedAdditiveCategoryModel(dto, existingCategory)
	if err := s.repo.UpdateAdditiveCategory(updatedCategory); err != nil {
		wrappedErr := utils.WrapError("failed to update additive category", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *additiveService) DeleteAdditiveCategory(categoryID uint) error {
	if err := s.repo.DeleteAdditiveCategory(categoryID); err != nil {
		wrappedErr := utils.WrapError("failed to delete additive category", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *additiveService) GetAdditiveCategoryByID(categoryID uint) (*types.AdditiveCategoryResponseDTO, error) {
	category, err := s.repo.GetAdditiveCategoryByID(categoryID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to fetch additive category", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if category == nil {
		return nil, fmt.Errorf("additive category with ID %d not found", categoryID)
	}

	return types.ConvertToAdditiveCategoryResponseDTO(category), nil
}

func (s *additiveService) GetAdditives(filter *types.AdditiveFilterQuery) ([]types.AdditiveDTO, error) {
	additives, err := s.repo.GetAdditives(filter)
	if err != nil {
		wrappedError := utils.WrapError("failed to retrieve additives", err)
		s.logger.Error(wrappedError)
		return nil, wrappedError
	}

	var additiveDTOs []types.AdditiveDTO
	for _, additive := range additives {
		additiveDTOs = append(additiveDTOs, *types.ConvertToAdditiveDTO(&additive))
	}

	return additiveDTOs, nil
}

func (s *additiveService) GetAdditivesByIDs(additiveIDs []uint) ([]types.AdditiveDTO, error) {
	additives, err := s.repo.GetAdditivesByIDs(additiveIDs)
	if err != nil {
		wrappedError := utils.WrapError("failed to retrieve additives by id list", err)
		s.logger.Error(wrappedError)
		return nil, wrappedError
	}

	var additiveDTOs []types.AdditiveDTO
	for _, additive := range additives {
		additiveDTOs = append(additiveDTOs, *types.ConvertToAdditiveDTO(&additive))
	}

	return additiveDTOs, nil
}

func (s *additiveService) CreateAdditive(dto *types.CreateAdditiveDTO) (uint, error) {
	additive := types.ConvertToAdditiveModel(dto)

	id, err := s.repo.CreateAdditive(additive)
	if err != nil {
		wrappedErr := utils.WrapError("failed to add additive", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *additiveService) UpdateAdditive(additiveID uint, dto *types.UpdateAdditiveDTO) error {
	updateModels := types.ConvertToUpdatedAdditiveModels(dto)

	if err := s.repo.UpdateAdditiveWithAssociations(additiveID, updateModels); err != nil {
		wrappedErr := utils.WrapError("failed to update additive with associations", err)
		s.logger.Error(wrappedErr)
		return err
	}
	return nil
}

func (s *additiveService) DeleteAdditive(additiveID uint) error {
	if err := s.repo.DeleteAdditive(additiveID); err != nil {
		wrappedErr := utils.WrapError("failed to delete additive", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *additiveService) GetAdditiveByID(additiveID uint) (*types.AdditiveDetailsDTO, error) {
	additive, err := s.repo.GetAdditiveByID(additiveID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to fetch additive by ID", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if additive == nil {
		return nil, fmt.Errorf("additive with ID %d not found", additiveID)
	}

	return types.ConvertToAdditiveDetailsDTO(additive), nil
}
